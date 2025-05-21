package main

import (
    "fmt"
    "log"
    "net/http"
    "chatting-service-app/db" 
)

func main() {
    // Connect to the database
    database, err := db.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }
    defer database.Close()

    // Just a simple HTTP handler for testing
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from my Go project!")
    })

    fmt.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
