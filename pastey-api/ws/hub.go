package ws

type Hub struct {
	Clients    map[int64]map[int64]*Client // map[userID]map[deviceID]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int64]map[int64]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Clients[client.UserID]; !ok {
				h.Clients[client.UserID] = make(map[int64]*Client)
			}

			h.Clients[client.UserID][client.DeviceID] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID][client.DeviceID]; ok {
				// Close client message channel
				delete(h.Clients[client.UserID], client.DeviceID)
				close(client.Message)

				// Remove user from clients if no more devices
				if len(h.Clients[client.UserID]) == 0 {
					delete(h.Clients, client.UserID)
				}
			}
		case message := <-h.Broadcast:
			if userClients, ok := h.Clients[message.UserID]; ok {
				if client, ok := userClients[message.ToDeviceID]; ok {
					select {
					case client.Message <- message:
					default:
						close(client.Message)
						delete(userClients, message.ToDeviceID)
						if len(userClients) == 0 {
							delete(h.Clients, message.UserID)
						}
					}
				}
			}
		}
	}
}
