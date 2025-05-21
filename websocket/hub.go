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
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		clientsByID: make(map[string]*Client),
		broadcast:   make(chan []byte),
		direct:      make(chan DirectMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) SendDirect(toID string, data []byte) {
	h.direct <- DirectMessage{ToID: toID, Data: data}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.clientsByID[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientsByID, client.ID)
				close(client.Send)
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