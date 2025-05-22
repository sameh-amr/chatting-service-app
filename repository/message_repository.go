package repository

import (
    "chatting-service-app/models"
    "chatting-service-app/db"
    "time"
)

type MessageRepository struct {}

func NewMessageRepository() *MessageRepository {
    return &MessageRepository{}
}

func (r *MessageRepository) SendMessage(msg *models.Message) error {
    return db.DB.Create(msg).Error
}

func (r *MessageRepository) GetMessagesBetweenUsers(user1ID, user2ID string) ([]models.Message, error) {
    var messages []models.Message
    err := db.DB.Where(
        "(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)", 
        user1ID, user2ID, user2ID, user1ID,
    ).Order("created_at asc").Find(&messages).Error
    return messages, err
}

func (r *MessageRepository) GetAllMessagesForUser(userID string) ([]models.Message, error) {
    var messages []models.Message
    err := db.DB.Where(
        "sender_id = ? OR recipient_id = ?",
        userID, userID,
    ).Order("created_at asc").Find(&messages).Error
    return messages, err
}

func (r *MessageRepository) CreateMessageRecipient(recipient *models.MessageRecipient) error {
    return db.DB.Create(recipient).Error
}

func (r *MessageRepository) SetDeliveredAt(messageID, recipientID string, deliveredAt time.Time) error {
    return db.DB.Model(&models.MessageRecipient{}).
        Where("message_id = ? AND recipient_id = ?", messageID, recipientID).
        Update("delivered_at", deliveredAt).Error
}

func (r *MessageRepository) SetReadAt(messageID, recipientID string, readAt time.Time) error {
    return db.DB.Model(&models.MessageRecipient{}).
        Where("message_id = ? AND recipient_id = ?", messageID, recipientID).
        Update("read_at", readAt).Error
}