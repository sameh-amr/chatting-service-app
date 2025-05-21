package main

import (
    "fmt"
    "log"
    "net/http"

    "chatting-service-app/db"
    "chatting-service-app/httphandlers"
    "chatting-service-app/models"
    "chatting-service-app/repository"
    "chatting-service-app/service"
)

func main() {
    // Connect to the database
    _, err := db.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    // Drop all tables using GORM for testing and development purposes
    // err = db.DB.Migrator().DropTable(&models.MessageRecipient{}, &models.Message{}, &models.User{}, &models.Session{})
    // if err != nil {
    //     log.Fatal("Failed to drop tables:", err)
    // }

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

    router := httphandlers.SetupRouter(userHandler)

    fmt.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
