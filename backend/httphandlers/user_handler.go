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
    err := h.userService.SignUp(req.Username, req.Email, req.Password)
    if err != nil {
        utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }
    user, err := h.userService.Authenticate(req.Email, req.Password)
    if err != nil || user == nil {
        utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not authenticate user after signup"})
        return
    }
    token, err := utils.GenerateJWT(user.ID.String())
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
        return
    }
    // Return both token and user data (id, username, email)
    utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
        "token": token,
        "user": map[string]interface{}{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
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
    user, err := h.userService.Authenticate(req.Email, req.Password)
    if err != nil || user == nil {
        utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
        return
    }
    token, err := utils.GenerateJWT(user.ID.String())
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
        return
    }
    // Return both token and user data (id, username, email)
    utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
        "token": token,
        "user": map[string]interface{}{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
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

// GetAllUsersExceptHandler returns all users except the authenticated user
func (h *UserHandler) GetAllUsersExceptHandler(w http.ResponseWriter, r *http.Request) {
    // JWT check
    authHeader := r.Header.Get("Authorization")
    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    userID, err := utils.ExtractUserIDFromJWT(tokenStr)
    if err != nil {
        utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or missing token"})
        return
    }
    users, err := h.userService.GetAllUsersExcept(userID)
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch users"})
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

// MeHandler returns the current authenticated user's data
func (h *UserHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    userID, err := utils.ExtractUserIDFromJWT(tokenStr)
    if err != nil {
        utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or missing token"})
        return
    }
    user, err := h.userService.GetUserByID(userID)
    if err != nil || user == nil {
        utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
        return
    }
    // Optionally, only return ID and Username for privacy
    result := map[string]interface{}{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
    }
    utils.WriteJSON(w, http.StatusOK, result)
}
