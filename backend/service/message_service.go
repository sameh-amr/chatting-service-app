package service

import (
    "chatting-service-app/dto"
    "chatting-service-app/models"
    "chatting-service-app/repository"
    "chatting-service-app/websocket"
    "encoding/json"
    "errors"
    "time"
)

type MessageService struct {
    repo              *repository.MessageRepository
    hub               *websocket.Hub
    recipientService *MessageRecipientService
}

func NewMessageService(repo *repository.MessageRepository, hub *websocket.Hub, recipientService *MessageRecipientService) *MessageService {
    return &MessageService{repo: repo, hub: hub, recipientService: recipientService}
}

func (s *MessageService) SendMessage(req dto.SendMessageRequest) error {
    if req.SenderID == (models.User{}).ID || req.Content == "" {
        return errors.New("missing required fields")
    }
    msg := &models.Message{
        SenderID:    req.SenderID,
        RecipientID: req.RecipientID,
        Content:     req.Content,
        MediaURL:    req.MediaURL,
        IsBroadcast: req.IsBroadcast,
        CreatedAt:   time.Now(),
    }
    err := s.repo.SendMessage(msg)
    if err != nil {
        return err
    }
    msgBytes, _ := json.Marshal(req)
    if req.IsBroadcast {
        // Broadcast: create MessageRecipient for all users except sender, send concurrently
        users, err := s.recipientService.userRepo.GetAllUsersExcept(req.SenderID.String())
        if err != nil {
            return err
        }
        for _, user := range users {
            recipient := &models.MessageRecipient{
                MessageID:   msg.ID,
                RecipientID: user.ID,
                DeliveredAt: nil,
                ReadAt:      nil,
            }
            _ = s.recipientService.Create(recipient)
        }
        if s.hub != nil {
            s.hub.BroadcastExcept(req.SenderID.String(), msgBytes)
        }
    } else {
        // 1:1: create MessageRecipient for recipient only
        recipient := &models.MessageRecipient{
            MessageID:   msg.ID,
            RecipientID: req.RecipientID,
            DeliveredAt: nil,
            ReadAt:      nil,
        }
        _ = s.recipientService.Create(recipient)
        if s.hub != nil {
            s.hub.SendDirect(req.RecipientID.String(), msgBytes)
        }
    }
    return nil
}

func (s *MessageService) SetDeliveredAt(messageID, recipientID string) error {
    return s.recipientService.SetDeliveredAt(messageID, recipientID)
}

func (s *MessageService) SetReadAt(messageID, recipientID string) error {
    return s.recipientService.SetReadAt(messageID, recipientID)
}

func (s *MessageService) GetMessagesBetweenUsers(user1ID, user2ID string) ([]models.Message, error) {
    return s.repo.GetMessagesBetweenUsers(user1ID, user2ID)
}

func (s *MessageService) GetAllMessagesForUser(userID string) ([]models.Message, error) {
    return s.repo.GetAllMessagesForUser(userID)
}
