package httphandlers

import (
    "net/http"

    "github.com/gorilla/mux"
    "chatting-service-app/websocket"
)

func SetupRouter(userHandler *UserHandler, hub *websocket.Hub, messageHandler *MessageHandler) *mux.Router {
    r := mux.NewRouter()

    authRouter := r.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/signup", userHandler.SignUpHandler).Methods("POST")
    authRouter.HandleFunc("/login", userHandler.LoginHandler).Methods("POST")

    // Message routes
    r.HandleFunc("/messages", messageHandler.SendMessageHandler).Methods("POST")
    r.HandleFunc("/messages", messageHandler.GetMessagesBetweenUsersHandler).Methods("GET")

    // WebSocket route
    r.HandleFunc("/ws", ServeWs(hub)).Methods("GET")

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello from my Go project!"))
    })

    return r
}