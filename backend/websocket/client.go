package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
	"strings"
	"chatting-service-app/models"
	"time"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
	ID   string // user ID or unique identifier
}

// Add a callback type for delivery/read status
// This should be set by the main app or handler to call the service
var OnMessageDelivered func(messageID, recipientID string)
var OnMessageRead func(messageID, recipientID string)

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		msgStr := string(message)
		if strings.Contains(msgStr, "\"type\":\"delivered\"") {
			idIdx := strings.Index(msgStr, "message_id")
			if idIdx != -1 {
				start := strings.Index(msgStr[idIdx:], ":") + idIdx + 2
				end := strings.Index(msgStr[start:], "\"") + start
				messageID := msgStr[start:end]
				if OnMessageDelivered != nil {
					OnMessageDelivered(messageID, c.ID)
				}
			}
		} else if strings.Contains(msgStr, "\"type\":\"read\"") {
			idIdx := strings.Index(msgStr, "message_id")
			if idIdx != -1 {
				start := strings.Index(msgStr[idIdx:], ":") + idIdx + 2
				end := strings.Index(msgStr[start:], "\"") + start
				messageID := msgStr[start:end]
				if OnMessageRead != nil {
					OnMessageRead(messageID, c.ID)
				}
			}
		} else {
			// Try to parse as wrapped message (with payload)
			var envelope map[string]interface{}
			if err := json.Unmarshal(message, &envelope); err == nil {
				if envelope["type"] == "message" && envelope["payload"] != nil {
					payloadBytes, _ := json.Marshal(envelope["payload"])
					var raw map[string]interface{}
					if err := json.Unmarshal(payloadBytes, &raw); err == nil {
						var chatMsg models.Message
						if sender, ok := raw["sender_id"].(string); ok {
							if uuidVal, err := uuid.Parse(sender); err == nil {
								chatMsg.SenderID = uuidVal
							}
						}
						if recipient, ok := raw["recipient_id"].(string); ok {
							if uuidVal, err := uuid.Parse(recipient); err == nil {
								chatMsg.RecipientID = uuidVal
							}
						}
						if content, ok := raw["content"].(string); ok {
							chatMsg.Content = content
						}
						if isBroadcast, ok := raw["is_broadcast"].(bool); ok {
							chatMsg.IsBroadcast = isBroadcast
						}
						if createdAt, ok := raw["created_at"].(string); ok {
							if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
								chatMsg.CreatedAt = t
							}
						}
						if mediaURL, ok := raw["media_url"].(string); ok {
							chatMsg.MediaURL = mediaURL
						}
						msgBytes, _ := json.Marshal(chatMsg)
						if chatMsg.Content != "" && chatMsg.RecipientID != uuid.Nil {
							if chatMsg.IsBroadcast {
								c.Hub.BroadcastExcept(c.ID, msgBytes)
							} else {
								c.Hub.SendDirect(chatMsg.RecipientID.String(), msgBytes)
							}
							continue
						}
					}
				}
			}
			// Fallback: Try to parse as raw message (no envelope)
			var raw map[string]interface{}
			if err := json.Unmarshal(message, &raw); err == nil {
				var chatMsg models.Message
				if sender, ok := raw["sender_id"].(string); ok {
					if uuidVal, err := uuid.Parse(sender); err == nil {
						chatMsg.SenderID = uuidVal
					}
				}
				if recipient, ok := raw["recipient_id"].(string); ok {
					if uuidVal, err := uuid.Parse(recipient); err == nil {
						chatMsg.RecipientID = uuidVal
					}
				}
				if content, ok := raw["content"].(string); ok {
					chatMsg.Content = content
				}
				if isBroadcast, ok := raw["is_broadcast"].(bool); ok {
					chatMsg.IsBroadcast = isBroadcast
				}
				if mediaURL, ok := raw["media_url"].(string); ok {
					chatMsg.MediaURL = mediaURL
				}
				msgBytes, _ := json.Marshal(chatMsg)
				if chatMsg.Content != "" && chatMsg.RecipientID != uuid.Nil {
					if chatMsg.IsBroadcast {
						c.Hub.BroadcastExcept(c.ID, msgBytes)
					} else {
						c.Hub.SendDirect(chatMsg.RecipientID.String(), msgBytes)
					}
					continue
				}
			}
			// If not a chat message, ignore or handle as needed
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}