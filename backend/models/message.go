package models

import (
    "time"
    "github.com/google/uuid"
)

type Message struct {
    ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    SenderID    uuid.UUID
    RecipientID uuid.UUID // Add this field for 1:1 messaging
    Content     string
    MediaURL    string
    IsBroadcast bool
    CreatedAt   time.Time
}
