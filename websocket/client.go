package websocket

import (
	"github.com/gorilla/websocket"
	"strings"
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
		// Here, parse the message for delivery/read events
		// For example, expect JSON: {"type":"delivered","message_id":"..."}
		msgStr := string(message)
		if strings.Contains(msgStr, "\"type\":\"delivered\"") {
			// Extract message_id (simple parsing for demo)
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
		}
		c.Hub.broadcast <- message
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