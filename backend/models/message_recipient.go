package models

import (
    "time"
    "github.com/google/uuid"
)

type MessageRecipient struct {
    ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    MessageID   uuid.UUID
    RecipientID uuid.UUID
    DeliveredAt *time.Time
    ReadAt      *time.Time
}
