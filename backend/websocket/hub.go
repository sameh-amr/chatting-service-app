package websocket

import (
	"encoding/json"
	"chatting-service-app/models"
)

// Extend OnlineStatusSetter to include GetUserByID for user data fetch
// This avoids import cycles and allows the hub to fetch user info
type OnlineStatusSetter interface {
	SetOnlineStatus(userID string, isOnline bool) error
	GetUserByID(userID string) (*models.User, error)
}

type DirectMessage struct {
	ToID string
	Data []byte
}

type Hub struct {
	clients     map[*Client]bool
	clientsByID map[string]*Client
	broadcast   chan []byte
	direct      chan DirectMessage
	register    chan *Client
	unregister  chan *Client
	userService OnlineStatusSetter // Use interface instead of concrete type
}

func NewHub(userService OnlineStatusSetter) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		clientsByID: make(map[string]*Client),
		broadcast:   make(chan []byte),
		direct:      make(chan DirectMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		userService: userService,
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) SendDirect(toID string, data []byte) {
	h.direct <- DirectMessage{ToID: toID, Data: data}
}

func (h *Hub) getOnlineUserIDs() []string {
	userIDs := make([]string, 0, len(h.clientsByID))
	for id := range h.clientsByID {
		userIDs = append(userIDs, id)
	}
	return userIDs
}

func (h *Hub) broadcastUserOnline(userID string) {
	// Fetch user data for the new online user
	var userData map[string]interface{}
	if h.userService != nil {
		if user, err := h.userService.GetUserByID(userID); err == nil && user != nil {
			userData = map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			}
		}
	}
	msg := map[string]interface{}{
		"type":   "user_online",
		"userId": userID,
	}
	if userData != nil {
		msg["user"] = userData
	}
	data, _ := json.Marshal(msg)
	for client := range h.clients {
		client.Send <- data
	}
}

func (h *Hub) broadcastUserOffline(userID string) {
	msg := map[string]interface{}{
		"type":   "user_offline",
		"userId": userID,
	}
	data, _ := json.Marshal(msg)
	for client := range h.clients {
		client.Send <- data
	}
}

func (h *Hub) sendOnlineUsersList(client *Client) {
	msg := map[string]interface{}{
		"type":    "online_users",
		"userIds": h.getOnlineUserIDs(),
	}
	data, _ := json.Marshal(msg)
	client.Send <- data
}

// BroadcastExcept sends a message to all connected clients except the sender (by user ID), using goroutines for concurrency.
func (h *Hub) BroadcastExcept(senderID string, data []byte) {
	for id, client := range h.clientsByID {
		if id == senderID {
			continue
		}
		go func(c *Client) {
			select {
			case c.Send <- data:
			default:
				close(c.Send)
				h.unregister <- c
			}
		}(client)
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.clientsByID[client.ID] = client
			if h.userService != nil {
				_ = h.userService.SetOnlineStatus(client.ID, true)
			}
			h.sendOnlineUsersList(client)
			h.broadcastUserOnline(client.ID)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientsByID, client.ID)
				close(client.Send)
				if h.userService != nil {
					_ = h.userService.SetOnlineStatus(client.ID, false)
				}
				h.broadcastUserOffline(client.ID)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
					delete(h.clientsByID, client.ID)
				}
			}
		case dm := <-h.direct:
			if client, ok := h.clientsByID[dm.ToID]; ok {
				select {
				case client.Send <- dm.Data:
				default:
					close(client.Send)
					delete(h.clients, client)
					delete(h.clientsByID, client.ID)
				}
			}
		}
	}
}