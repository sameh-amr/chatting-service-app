package db

import (
    "database/sql"
    "fmt"
    "os"
    "time"

    _ "github.com/lib/pq"
)

// ConnectDB opens a connection to PostgreSQL and returns *sql.DB
func ConnectDB() (*sql.DB, error) {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    // Build the connection string
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    var db *sql.DB
    var err error

    // Retry logic: try connecting up to 10 times with 2 seconds interval
    for i := 0; i < 10; i++ {
        db, err = sql.Open("postgres", psqlInfo)
        if err == nil {
            err = db.Ping()
            if err == nil {
                fmt.Println("Successfully connected to database!")
                return db, nil
            }
        }
        fmt.Printf("Attempt %d: Unable to connect to DB, retrying in 2 seconds...\n", i+1)
        time.Sleep(2 * time.Second)
    }

    return nil, fmt.Errorf("could not connect to database after 10 attempts: %w", err)
}
