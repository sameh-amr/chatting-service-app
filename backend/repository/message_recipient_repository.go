package repository

import (
	"chatting-service-app/models"
	"chatting-service-app/db"
	"time"
)

type MessageRecipientRepository struct {}

func NewMessageRecipientRepository() *MessageRecipientRepository {
	return &MessageRecipientRepository{}
}

func (r *MessageRecipientRepository) Create(recipient *models.MessageRecipient) error {
	return db.DB.Create(recipient).Error
}

func (r *MessageRecipientRepository) SetDeliveredAt(messageID, recipientID string, deliveredAt time.Time) error {
	return db.DB.Model(&models.MessageRecipient{}).
		Where("message_id = ? AND recipient_id = ?", messageID, recipientID).
		Update("delivered_at", deliveredAt).Error
}

func (r *MessageRecipientRepository) SetReadAt(messageID, recipientID string, readAt time.Time) error {
	return db.DB.Model(&models.MessageRecipient{}).
		Where("message_id = ? AND recipient_id = ?", messageID, recipientID).
		Update("read_at", readAt).Error
}
