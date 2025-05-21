package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq"
)

// ConnectDB opens a connection to PostgreSQL and returns *sql.DB
func ConnectDB() (*sql.DB, error) {
    host := os.Getenv("db")
    port := os.Getenv("5432")
    user := os.Getenv("postgres")
    password := os.Getenv("secret")
    dbname := os.Getenv("chatdb")

    // Build the connection string
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }

    // Test the connection
    err = db.Ping()
    if err != nil {
        return nil, err
    }

    fmt.Println("Successfully connected to database!")
    return db, nil
}
