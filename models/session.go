package models

import (
    "time"
    "github.com/google/uuid"
)

type Session struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID    uuid.UUID
    Token     string
    ExpiresAt time.Time
    CreatedAt time.Time
}
