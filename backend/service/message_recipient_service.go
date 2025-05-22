package service

import (
	"chatting-service-app/models"
	"chatting-service-app/repository"
	"time"
)

type MessageRecipientService struct {
	repo     *repository.MessageRecipientRepository
	userRepo *repository.UserRepository
}

func NewMessageRecipientService(repo *repository.MessageRecipientRepository, userRepo *repository.UserRepository) *MessageRecipientService {
	return &MessageRecipientService{repo: repo, userRepo: userRepo}
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
