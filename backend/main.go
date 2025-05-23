package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	"chatting-service-app/db"
	"chatting-service-app/httphandlers"
	"chatting-service-app/models"
	"chatting-service-app/repository"
	"chatting-service-app/service"
	"chatting-service-app/websocket"
)

var messageServiceGlobal *service.MessageService

func main() {
	// Connect to the database
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Drop all tables using GORM for testing and development purposes
	err = db.DB.Migrator().DropTable(&models.MessageRecipient{}, &models.Message{}, &models.User{}, &models.Session{})
	if err != nil {
	    log.Fatal("Failed to drop tables:", err)
	}

	// GORM AutoMigrate for your models
	err = db.DB.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.MessageRecipient{},
		&models.Session{},
	)
	if err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}

	// Set up repository, service, and handler
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := httphandlers.NewUserHandler(userService)

	// Message recipient repository and service
	messageRecipientRepo := repository.NewMessageRecipientRepository()
	messageRecipientService := service.NewMessageRecipientService(messageRecipientRepo, userRepo)

	// Message repository, service, and handler
	messageRepo := repository.NewMessageRepository()
	messageService := service.NewMessageService(messageRepo, nil, messageRecipientService)
	messageServiceGlobal = messageService
	messageHandler := httphandlers.NewMessageHandler(messageService, messageRecipientService)

	// Start the WebSocket hub
	hub := websocket.NewHub(userService) // userService implements OnlineStatusSetter
	go hub.Run()

	router := httphandlers.SetupRouter(userHandler, hub, messageHandler, messageRecipientService)

	// Add CORS middleware
	h := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	)(router)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
