package httphandlers

import (
    "chatting-service-app/service"
    "chatting-service-app/utils"
    "net/http"
    "strings"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
    return &UserHandler{userService: us}
}

type signUpRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Helper: check if method is POST
func requirePost(w http.ResponseWriter, r *http.Request) bool {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return false
    }
    return true
}

func (h *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if !requirePost(w, r) {
        return
    }
    var req signUpRequest
    if !utils.DecodeJSON(r, &req, w) {
        return
    }
    token, err := h.userService.SignUpAndToken(req.Username, req.Email, req.Password)
    if err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }
    utils.WriteJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
    if !requirePost(w, r) {
        return
    }
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if !utils.DecodeJSON(r, &req, w) {
        return
    }
    token, err := h.userService.LoginAndToken(req.Email, req.Password)
    if err != nil {
        utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
        return
    }
    utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) GetOnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
    // JWT check
    authHeader := r.Header.Get("Authorization")
    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    _, err := utils.ExtractUserIDFromJWT(tokenStr)
    if err != nil {
        utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or missing token"})
        return
    }
    users, err := h.userService.GetOnlineUsers()
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch online users"})
        return
    }
    // Optionally, only return ID and Username for privacy
    var result []map[string]interface{}
    for _, u := range users {
        result = append(result, map[string]interface{}{
            "id":       u.ID,
            "username": u.Username,
        })
    }
    utils.WriteJSON(w, http.StatusOK, result)
}
