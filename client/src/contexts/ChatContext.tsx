import { createContext, useContext, useState, useEffect, ReactNode, useRef, MutableRefObject } from 'react';
import { useAuth } from './AuthContext';
import { Message, User } from '../types';
import api from '../services/api';
import { setupWebSocket, disconnectWebSocket, sendWebSocketMessage } from '../services/websocket';

interface ChatContextType {
  messages: Message[];
  onlineUsers: User[];
  selectedUser: User | null;
  isLoadingMessages: boolean;
  setSelectedUser: (user: User | null) => void;
  sendMessage: (content: string, recipientId: string, mediaUrl?: string, isBroadcast?: boolean) => Promise<void>;
  uploadFile: (file: File) => Promise<string>;
  downloadFile: (filename: string) => Promise<void>;
  markAsDelivered: (messageId: string, recipientId: string) => Promise<void>;
  markAsRead: (messageId: string, recipientId: string) => Promise<void>;
}

const ChatContext = createContext<ChatContextType | undefined>(undefined);

export const useChat = () => {
  const context = useContext(ChatContext);
  if (!context) {
    throw new Error('useChat must be used within a ChatProvider');
  }
  return context;
};

// Utility to convert a single message from PascalCase to camelCase
function mapMessageFromApi(msg: any): Message {
  return {
    id: msg.ID || msg.id || String(msg.id || msg.ID),
    sender_id: msg.SenderID || msg.sender_id || String(msg.sender_id || msg.SenderID),
    recipient_id: msg.RecipientID || msg.recipient_id || String(msg.recipient_id || msg.RecipientID),
    content: msg.Content || msg.content,
    media_url: msg.MediaURL || msg.media_url,
    is_broadcast: msg.IsBroadcast ?? msg.is_broadcast ?? false,
    created_at: msg.CreatedAt || msg.created_at || (msg.timestamp ? msg.timestamp : new Date().toISOString()),
    delivered: msg.Delivered ?? msg.delivered ?? false,
    read: msg.Read ?? msg.read ?? false,
  };
}

// Utility to map an array of messages
function mapMessagesFromApi(msgs: any[]): Message[] {
  return msgs.map(mapMessageFromApi);
}

export const ChatProvider = ({ children }: { children: ReactNode }) => {
  const { token, userId, isAuthenticated } = useAuth();
  const [messages, setMessages] = useState<Message[]>([]);
  const [onlineUsers, setOnlineUsers] = useState<User[]>([]);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);

  const selectedUserRef = useRef<User | null>(null) as MutableRefObject<User | null>;
  const userIdRef = useRef<string | null | undefined>(undefined) as MutableRefObject<string | null | undefined>;

  useEffect(() => {
    selectedUserRef.current = selectedUser;
  }, [selectedUser]);

  useEffect(() => {
    userIdRef.current = userId;
  }, [userId]);

  // Fetch online users on login/logout
  useEffect(() => {
    if (isAuthenticated && token) {
      fetchOnlineUsers();
    } else {
      // Only clear if not authenticated (true logout)
      setOnlineUsers([]);
    }
  }, [isAuthenticated, token]);

  // Handler for WebSocket message status updates
  const handleWebSocketMessageStatusUpdate = (messageId: string, status: "delivered" | "read") => {
    setMessages(prev =>
      prev.map(msg =>
        msg.id === messageId ? { ...msg, [status]: true } : msg
      )
    );
  };

  useEffect(() => {
    if (isAuthenticated && token) {
      setupWebSocket(token, handleWebSocketMessage, handleWebSocketConnection, handleWebSocketMessageStatusUpdate);
      fetchOnlineUsers();

      return () => {
        disconnectWebSocket();
      };
    }
  }, [isAuthenticated, token]);

  useEffect(() => {
    if (selectedUser && userId) {
      fetchMessages(userId, selectedUser.id);
    }
  }, [selectedUser, userId]);

  const fetchOnlineUsers = async () => {
    try {
      const response = await api.get('/auth/users');
      setOnlineUsers(response.data);
    } catch (error) {
      console.error('Error fetching online users:', error);
    }
  };

  const fetchMessages = async (user1Id: string, user2Id: string) => {
    try {
      setIsLoadingMessages(true);
      const response = await api.get(`/messages?user1=${user1Id}&user2=${user2Id}`);
      setMessages(mapMessagesFromApi(response.data));
    } catch (error) {
      console.error('Error fetching messages:', error);
    } finally {
      setIsLoadingMessages(false);
    }
  };

  const handleWebSocketMessage = (message: any) => {
    console.log('[WebSocket] Raw incoming message:', message);
    const mapped = mapMessageFromApi(message);
    const currentSelectedUser = selectedUserRef.current;
    const currentUserId = userIdRef.current;
    console.log('[WebSocket] Mapped message:', mapped, 'selectedUser:', currentSelectedUser, 'userId:', currentUserId);
    // Allow messages with content (even if no id) for real-time recipient display
    if (!mapped.content) {
      console.warn('Received message with no content:', message);
      return;
    }
    setMessages(prev => {
      if (
        currentSelectedUser &&
        (
          (mapped.sender_id === currentSelectedUser.id && mapped.recipient_id === currentUserId) ||
          (mapped.sender_id === currentUserId && mapped.recipient_id === currentSelectedUser.id)
        )
      ) {
        if (mapped.sender_id === currentSelectedUser.id) {
          markAsDelivered(mapped.id, currentUserId!);
        }
        return [...prev, mapped];
      }
      if (mapped.is_broadcast) {
        return [...prev, mapped];
      }
      return prev;
    });
  };

  const handleWebSocketConnection = (usersList: User[]) => {
    setOnlineUsers(usersList);
  };

  const sendMessage = async (
    content: string,
    recipientId: string,
    mediaUrl?: string,
    isBroadcast: boolean = false
  ) => {
    try {
      if (!userId) return;

      console.log('Sending message:', {
        sender_id: userId,
        recipient_id: recipientId,
        content,
        media_url: mediaUrl,
        is_broadcast: isBroadcast,
      });;

      // Convert to snake_case for backend
      const newMessageApi = {
        sender_id: userId,
        recipient_id: recipientId,
        content,
        media_url: mediaUrl,
        is_broadcast: isBroadcast,
        created_at: new Date().toISOString(), // Always send created_at
      };

      // Optimistically add the message for the sender
      const optimisticId = `optimistic-${Date.now()}`;
      const optimisticMessage = {
        id: optimisticId,
        sender_id: userId!,
        recipient_id: recipientId,
        content,
        media_url: mediaUrl,
        is_broadcast: isBroadcast,
        created_at: newMessageApi.created_at, // Use the same timestamp
        delivered: false,
        read: false,
      };
      setMessages(prev => [...prev, optimisticMessage]);

      sendWebSocketMessage(newMessageApi);

      const response = await api.post('/messages', newMessageApi);
      const mapped = mapMessageFromApi(response.data);
      if (mapped.id && mapped.content) {
        setMessages(prev => [
          ...prev.filter(m => m.id !== optimisticId),
          mapped
        ]);
      }

     // return response.data;
    } catch (error) {
      console.error('Error sending message:', error);
      throw error;
    }
  };

  const uploadFile = async (file: File): Promise<string> => {
    try {
      const formData = new FormData();
      formData.append('file', file);

      const response = await api.post('/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      return response.data.url;
    } catch (error) {
      console.error('Error uploading file:', error);
      throw error;
    }
  };

  const downloadFile = async (filename: string) => {
    try {
      const response = await api.get(`/download?file=${filename}`, {
        responseType: 'blob',
      });

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
    } catch (error) {
      console.error('Error downloading file:', error);
      throw error;
    }
  };

  const markAsDelivered = async (messageId: string, recipientId: string) => {
    try {
      await api.post(`/messages/delivered?message_id=${messageId}&recipient_id=${recipientId}`);
      setMessages(prev =>
        prev.map(msg =>
          msg.id === messageId ? { ...msg, delivered: true } : msg
        )
      );
    } catch (error) {
      console.error('Error marking message as delivered:', error);
    }
  };

  const markAsRead = async (messageId: string, recipientId: string) => {
    try {
      await api.post(`/messages/read?message_id=${messageId}&recipient_id=${recipientId}`);
      setMessages(prev =>
        prev.map(msg =>
          msg.id === messageId ? { ...msg, read: true } : msg
        )
      );
    } catch (error) {
      console.error('Error marking message as read:', error);
    }
  };

  const value = {
    messages,
    onlineUsers,
    selectedUser,
    isLoadingMessages,
    setSelectedUser,
    sendMessage,
    uploadFile,
    downloadFile,
    markAsDelivered,
    markAsRead,
  };

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
};
