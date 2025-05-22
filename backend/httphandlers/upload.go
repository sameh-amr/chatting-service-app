package httphandlers

import (
	"chatting-service-app/utils"
	"fmt"
	"net/http"
	"strings"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Require JWT authentication
	authHeader := r.Header.Get("Authorization")
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}
	_, err := utils.ExtractUserIDFromJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileURL, err := utils.SaveUploadedFile(file, handler, "uploads")
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"url":"%s"}`, fileURL)))
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Require JWT authentication
	authHeader := r.Header.Get("Authorization")
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}
	_, err := utils.ExtractUserIDFromJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Get file URL from query param and extract filename
	fileURL := r.URL.Query().Get("file")
	if fileURL == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}
	// Remove any leading slashes and 'uploads/' prefix if present
	filename := fileURL
	if idx := strings.LastIndex(fileURL, "/"); idx != -1 {
		filename = fileURL[idx+1:]
	}
	filePath := "uploads/" + filename

	http.ServeFile(w, r, filePath)
}
