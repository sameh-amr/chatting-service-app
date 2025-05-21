package dto

import (
	"github.com/google/uuid"
	"time"
)

type SendMessageRequest struct {
	SenderID    uuid.UUID `json:"sender_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Content     string    `json:"content"`
	MediaURL    string    `json:"media_url"`
	IsBroadcast bool      `json:"is_broadcast"`
	CreatedAt   time.Time `json:"created_at"`
}
