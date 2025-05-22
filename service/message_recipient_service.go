package service

import (
	"chatting-service-app/models"
	"chatting-service-app/repository"
	"time"
)

type MessageRecipientService struct {
	repo *repository.MessageRecipientRepository
}

func NewMessageRecipientService(repo *repository.MessageRecipientRepository) *MessageRecipientService {
	return &MessageRecipientService{repo: repo}
}

func (s *MessageRecipientService) Create(recipient *models.MessageRecipient) error {
	return s.repo.Create(recipient)
}

func (s *MessageRecipientService) SetDeliveredAt(messageID, recipientID string) error {
	return s.repo.SetDeliveredAt(messageID, recipientID, time.Now())
}

func (s *MessageRecipientService) SetReadAt(messageID, recipientID string) error {
	return s.repo.SetReadAt(messageID, recipientID, time.Now())
}
