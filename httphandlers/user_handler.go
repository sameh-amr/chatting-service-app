package httphandlers

import (
    "chatting-service-app/service"
    "chatting-service-app/utils"
    "net/http"
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
