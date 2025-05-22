package websocket

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
	userService *UserService // Add this for online status updates
}

func NewHub(userService *UserService) *Hub {
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
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientsByID, client.ID)
				close(client.Send)
				if h.userService != nil {
					_ = h.userService.SetOnlineStatus(client.ID, false)
				}
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