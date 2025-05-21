package models

import (
    "time"
    "github.com/google/uuid"
)

type Message struct {
    ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    SenderID    uuid.UUID
    Content     string
    MediaURL    string
    IsBroadcast bool
    CreatedAt   time.Time
}
