let socket: WebSocket | null = null;
let reconnectTimeout: NodeJS.Timeout | null = null;
let currentToken: string | null = null;

// Types for message handlers (kept for type safety, but handlers won't be called)
type MessageHandler = (message: any) => void;
type ConnectionHandler = (data: any, type: string) => void;
type MessageStatusUpdateHandler = (messageId: string, status: "delivered" | "read") => void;

// Store handlers (kept for potential reconnection logic, but won't be actively used for messages)
let messageHandler: MessageHandler | null = null;
let connectionHandler: ConnectionHandler | null = null;
let messageStatusUpdateHandler: MessageStatusUpdateHandler | null = null;

/**
 * Set up WebSocket connection
 */
export const setupWebSocket = (
  token: string,
  onMessage: MessageHandler, // This handler will not be called in this modified version
  onConnection: ConnectionHandler, // This handler will not be called in this modified version
  onMessageStatusUpdate: MessageStatusUpdateHandler // New handler for status updates
) => {
  // Store handlers for reconnection (kept for the reconnect logic in onclose)
  messageHandler = onMessage;
  connectionHandler = onConnection;
  messageStatusUpdateHandler = onMessageStatusUpdate;
  currentToken = token;

  // Close existing connection if any
  if (socket) {
    socket.close();
  }

  // Create new WebSocket connection with token as query param
  let wsHost: string;
  if (import.meta.env.VITE_API_URL) {
    // Use backend service name in Docker Compose
    wsHost = import.meta.env.VITE_API_URL.replace(/^https?:\/\//, '').replace(/\/.*$/, '');
  } else if (
    window.location.hostname === "localhost" ||
    window.location.hostname === "127.0.0.1"
  ) {
    wsHost = "localhost";
  } else {
    wsHost = "backend";
  }
  const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws";
  const wsUrl = `${wsProtocol}://${wsHost}:8080/ws?token=${token}`;
  socket = new WebSocket(wsUrl);

  // Set up event handlers
  socket.onopen = () => {
    // WebSocket connected
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }
  };

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);

      // Accept both wrapped and raw message formats
      if (data.type === "message" && data.payload) {
        onMessage(data.payload);
      } else if (data.type === undefined && data.Content && data.SenderID && data.RecipientID) {
        // Backend sent a raw message object (PascalCase)
        onMessage(data);
      } else if (data.type === undefined && data.content && data.sender_id && data.recipient_id) {
        // Backend sent a raw message object (camelCase)
        onMessage(data);
      } else if (typeof data.type === "string") {
        // Only call connectionHandler for known presence types
        switch (data.type) {
          case "delivered":
            if (data.message_id) {
              updateMessageStatus(data.message_id, "delivered");
            }
            break;
          case "read":
            if (data.message_id) {
              updateMessageStatus(data.message_id, "read");
            }
            break;
          case "online_users":
          case "user_online":
          case "user_offline":
            if (connectionHandler) connectionHandler(data, data.type);
            break;
          case "error":
            console.error("WebSocket error:", data.payload);
            break;
          default:
            console.log("Unknown message type:", data);
        }
      }
    } catch (error) {
      console.error("Error parsing WebSocket message:", error);
    }
  };

  socket.onerror = (error) => {
    console.error("WebSocket error:", error);
  };

  socket.onclose = (event) => {
    // Attempt reconnection only if token is present (user is authenticated)
    if (currentToken) {
      reconnectTimeout = setTimeout(() => {
        if (currentToken && messageHandler && connectionHandler && messageStatusUpdateHandler) {
          setupWebSocket(currentToken, messageHandler, connectionHandler, messageStatusUpdateHandler);
        }
      }, 2000); // 2s delay before reconnect
    }
  };

  return socket;
};

/**
 * Send a message through WebSocket
 */
export const sendWebSocketMessage = (message: any) => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(
      JSON.stringify({
        type: "message",
        payload: message,
      })
    );
    return true;
  }
  return false;
};

/**
 * Send delivered status
 */
export const sendDeliveredStatus = (messageId: string) => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(
      JSON.stringify({
        type: "delivered",
        message_id: messageId,
      })
    );
    return true;
  }
  return false;
};

/**
 * Send read status
 */
export const sendReadStatus = (messageId: string) => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(
      JSON.stringify({
        type: "read",
        message_id: messageId,
      })
    );
    return true;
  }
  return false;
};

/**
 * Update message status in local state
 */
const updateMessageStatus = (
  messageId: string,
  status: "delivered" | "read"
) => {
  if (messageStatusUpdateHandler) {
    messageStatusUpdateHandler(messageId, status);
  }
};

/**
 * Disconnect WebSocket and prevent reconnection
 */
export const disconnectWebSocket = () => {
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout);
    reconnectTimeout = null;
  }
  currentToken = null;
  if (socket) {
    socket.close();
    socket = null;
  }
};
