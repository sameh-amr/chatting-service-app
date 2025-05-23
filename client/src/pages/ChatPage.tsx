import { useState, useEffect, useRef } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useChat } from '../contexts/ChatContext';
import { LogOut, Send, Upload, Paperclip, Users } from 'lucide-react';
import Sidebar from '../components/Sidebar';
import Message from '../components/Message';
import Button from '../components/Button';
import Avatar from '../components/Avatar';
import { Message as MessageType } from '../types';

const ChatPage = () => {
  const { user, logout } = useAuth();
  const { 
    messages, 
    onlineUsers, 
    onlineUserIds, // Get onlineUserIds for real-time presence
    selectedUser, 
    setSelectedUser, 
    sendMessage,
    uploadFile,
    isLoadingMessages,
    markAsRead
  } = useChat();
  
  const [messageText, setMessageText] = useState('');
  const [isSending, setIsSending] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [showMobileSidebar, setShowMobileSidebar] = useState(false);
  
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Scroll to bottom when messages change
  useEffect(() => {
    scrollToBottom();
  }, [messages]);
  
  // Mark messages as read when they become visible
  // useEffect(() => {
  //   if (selectedUser && user) {
  //     messages.forEach(message => {
  //       if (message.sender_id === selectedUser.id && message.recipient_id === user.id && !message.read) {
  //         markAsRead(message.id, user.id);
  //       }
  //     });
  //   }
  // }, [messages, selectedUser, user]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!messageText.trim() || !selectedUser) return;
    
    try {
      setIsSending(true);
      await sendMessage(messageText, selectedUser.id);
      setMessageText('');
    } catch (error) {
      console.error('Error sending message:', error);
    } finally {
      setIsSending(false);
    }
  };

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !selectedUser) return;
    
    try {
      setIsUploading(true);
      const fileUrl = await uploadFile(file);
      await sendMessage(
        `Sent a file: ${file.name}`, 
        selectedUser.id,
        fileUrl
      );
    } catch (error) {
      console.error('Error uploading file:', error);
    } finally {
      setIsUploading(false);
      // Reset the file input
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  // New: Broadcast handler
  const handleBroadcastMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!messageText.trim() || !user?.id) return;
    try {
      setIsSending(true);
      // For broadcast, send the sender's own UUID as recipientId
      await sendMessage(messageText, user.id, undefined, true);
      setMessageText('');
    } catch (error) {
      console.error('Error sending broadcast message:', error);
    } finally {
      setIsSending(false);
    }
  };

  const toggleMobileSidebar = () => {
    setShowMobileSidebar(!showMobileSidebar);
  };

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Mobile Sidebar Toggle */}
      <div className="md:hidden fixed bottom-4 left-4 z-10">
        <button
          onClick={toggleMobileSidebar}
          className="bg-blue-600 text-white p-3 rounded-full shadow-lg hover:bg-blue-700 transition-colors"
        >
          <Users size={24} />
        </button>
      </div>
      
      {/* Sidebar */}
      <Sidebar
        onlineUsers={onlineUsers}
        onlineUserIds={onlineUserIds} // Pass onlineUserIds for real-time presence
        selectedUser={selectedUser}
        setSelectedUser={setSelectedUser}
        showMobile={showMobileSidebar}
        toggleMobile={toggleMobileSidebar}
        currentUser={user}
      />
      
      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col h-full">
        {/* Chat Header */}
        <header className="bg-white shadow-sm py-3 px-4 flex justify-between items-center">
          <div className="flex items-center">
            {selectedUser ? (
              <>
                <Avatar 
                  user={selectedUser} 
                  showStatus 
                  isOnline={onlineUsers.some(u => u.id === selectedUser.id)}
                />
                <div className="ml-3">
                  <h2 className="font-semibold text-gray-800">{selectedUser.username}</h2>
                  <p className="text-xs text-gray-500">
                    {onlineUsers.some(u => u.id === selectedUser.id) ? 'Online' : 'Offline'}
                  </p>
                </div>
              </>
            ) : (
              <h2 className="font-semibold text-gray-800">Select a conversation</h2>
            )}
          </div>
          
          <button
            onClick={logout}
            className="text-gray-500 hover:text-gray-700 flex items-center gap-1"
          >
            <LogOut size={18} />
            <span className="hidden sm:inline">Logout</span>
          </button>
        </header>
        
        {/* Messages Area */}
        <div className="flex-1 overflow-y-auto p-4 bg-gray-50">
          {selectedUser ? (
            <>
              {isLoadingMessages ? (
                <div className="flex justify-center items-center h-full">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                </div>
              ) : messages.length > 0 ? (
                <div className="space-y-4">
                  {messages.map((msg: MessageType) => (
                    <Message
                      key={msg.id}
                      message={msg}
                      isCurrentUser={msg.sender_id === user?.id}
                    />
                  ))}
                  <div ref={messagesEndRef} />
                </div>
              ) : (
                <div className="flex items-center justify-center h-full">
                  <p className="text-gray-500">No messages yet. Start the conversation!</p>
                </div>
              )}
            </>
          ) : (
            <div className="flex flex-col items-center justify-center h-full">
              <div className="bg-blue-50 p-6 rounded-full mb-4">
                <MessageCircle className="h-12 w-12 text-blue-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-800 mb-2">Your Messages</h3>
              <p className="text-gray-500 text-center max-w-sm">
                Select a contact to start chatting
              </p>
            </div>
          )}
        </div>
        
        {/* Message Input */}
        <form onSubmit={selectedUser ? handleSendMessage : (e) => e.preventDefault()} className="bg-white border-t p-3">
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className="text-gray-500 hover:text-blue-600 p-2 rounded-full hover:bg-gray-100"
              disabled={isUploading}
            >
              {isUploading ? <Upload className="h-5 w-5 animate-pulse" /> : <Paperclip className="h-5 w-5" />}
            </button>
            <input
              type="text"
              value={messageText}
              onChange={(e) => setMessageText(e.target.value)}
              placeholder="Type a message..."
              className="flex-1 border rounded-full py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-200"
              disabled={isSending}
            />
            {/* Send to selected user */}
            {selectedUser && (
              <Button
                type="submit"
                variant="primary"
                isLoading={isSending}
                className="rounded-full px-3 py-2"
                disabled={!messageText.trim()}
              >
                <Send className="h-5 w-5" />
              </Button>
            )}
            {/* Broadcast button, always visible */}
            <Button
              type="button"
              variant="secondary"
              isLoading={isSending}
              className="rounded-full px-3 py-2"
              disabled={!messageText.trim()}
              onClick={handleBroadcastMessage}
              title="Send broadcast message"
            >
              <Send className="h-5 w-5" />
              <span className="ml-1 hidden sm:inline">Broadcast</span>
            </Button>
            <input
              ref={fileInputRef}
              type="file"
              onChange={handleFileSelect}
              className="hidden"
            />
          </div>
        </form>
      </div>
    </div>
  );
};

const MessageCircle = ({ className }: { className?: string }) => (
  <svg 
    className={className}
    viewBox="0 0 24 24" 
    fill="none" 
    stroke="currentColor" 
    strokeWidth="2" 
    strokeLinecap="round" 
    strokeLinejoin="round"
  >
    <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z" />
  </svg>
);

export default ChatPage;