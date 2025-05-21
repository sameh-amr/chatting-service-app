package utils

import (
    "encoding/json"
    "net/http"
)

// Helper: decode JSON request
func DecodeJSON(r *http.Request, v interface{}, w http.ResponseWriter) bool {
    err := json.NewDecoder(r.Body).Decode(v)
    if err != nil {
        http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
        return false
    }
    return true
}

// Helper: write JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
