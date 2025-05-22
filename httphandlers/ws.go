package httphandlers

import (
	ws "chatting-service-app/websocket"
	"chatting-service-app/utils"
	"net/http"
	"strings"
	websocket "github.com/gorilla/websocket"
	"chatting-service-app/service"
)

var (
	upgrader          = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // For dev only; restrict in prod
	}
	messageServiceGlobal *service.MessageService
)

func ServeWs(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from JWT in Authorization header
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
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

		client := &ws.Client{
			Hub:  hub,
			Conn: conn,
			Send: make(chan []byte, 256),
			ID:   userID,
		}
		hub.Register(client)

		// Wire up delivery/read status callbacks
		ws.OnMessageDelivered = func(messageID, recipientID string) {
			// You may want to get the messageService from context or global
			// For demo, assume a global variable or singleton
			if messageServiceGlobal != nil {
				_ = messageServiceGlobal.SetDeliveredAt(messageID, recipientID)
			}
		}
		ws.OnMessageRead = func(messageID, recipientID string) {
			if messageServiceGlobal != nil {
				_ = messageServiceGlobal.SetReadAt(messageID, recipientID)
			}
		}

		// Start pumps
		go client.WritePump()
		go client.ReadPump()
	}
}
