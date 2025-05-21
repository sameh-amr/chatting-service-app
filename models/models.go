package models

import "time"

type User struct {
    ID       string
    Username string
    Password string
    Email    string
}

type Message struct {
    ID         string
    SenderID   string
    Content    string
    MediaURL   *string 
    CreatedAt  time.Time
}

type MessageRecipient struct {
    ID         string
    MessageID  string
    RecipientID string
    IsRead     bool
}
