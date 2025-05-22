package httphandlers

import (
	ws "chatting-service-app/websocket"
	"chatting-service-app/utils"
	"fmt"
	"net/http"
	websocket "github.com/gorilla/websocket"
	"chatting-service-app/service"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins for development; restrict in production
			return true
		},
	}
	messageServiceGlobal *service.MessageService
)

func ServeWs(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT from query param instead of Authorization header
		tokenStr := r.URL.Query().Get("token")
		userID, err := utils.ExtractUserIDFromJWT(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Could not upgrade to websocket", http.StatusInternalServerError)
			return
		}
		fmt.Println("WebSocket: Connection upgraded for userIDD:", userID)

		client := &ws.Client{
			Hub:  hub,
			Conn: conn,
			Send: make(chan []byte, 256),
			ID:   userID,
		}
		hub.Register(client)
		fmt.Println("WebSocket: Client registered for userIDDD:", userID)

		// Wire up delivery/read status callbacks
		ws.OnMessageDelivered = func(messageID, recipientID string) {
			if messageServiceGlobal != nil {
				_ = messageServiceGlobal.SetDeliveredAt(messageID, recipientID)
			}
		}
		ws.OnMessageRead = func(messageID, recipientID string) {
			if messageServiceGlobal != nil {
				_ = messageServiceGlobal.SetReadAt(messageID, recipientID)
			}
		}

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("WebSocket: Panic recovered for userID %s: %v\n", userID, r)
				if conn != nil {
					conn.Close()
				}
			}
		}()

		// Start pumps
		go client.WritePump()
		go client.ReadPump()
	}
}
