package httphandlers

import (
    "chatting-service-app/dto"
    "chatting-service-app/service"
    "chatting-service-app/utils"
    "net/http"
    "strings"
    "github.com/google/uuid"
)

type MessageHandler struct {
    messageService *service.MessageService
}

func NewMessageHandler(ms *service.MessageService) *MessageHandler {
    return &MessageHandler{messageService: ms}
}

func (h *MessageHandler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    var req dto.SendMessageRequest
    if !utils.DecodeJSON(r, &req, w) {
        return
    }
    authHeader := r.Header.Get("Authorization")
    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    senderID, err := utils.ExtractUserIDFromJWT(tokenStr)
    if err != nil {
        utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or missing token"})
        return
    }
    senderUUID, err := uuid.Parse(senderID)
    if err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid sender id"})
        return
    }
    req.SenderID = senderUUID
    // Optionally validate RecipientID is a valid uuid.UUID (if needed)
    err = h.messageService.SendMessage(req)
    if err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }
    utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "message sent"})
}

func (h *MessageHandler) GetMessagesBetweenUsersHandler(w http.ResponseWriter, r *http.Request) {
    user1ID := r.URL.Query().Get("user1")
    user2ID := r.URL.Query().Get("user2")
    if user1ID == "" || user2ID == "" {
        utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "user1 and user2 are required"})
        return
    }
    messages, err := h.messageService.GetMessagesBetweenUsers(user1ID, user2ID)
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch messages"})
        return
    }
    utils.WriteJSON(w, http.StatusOK, messages)
}
