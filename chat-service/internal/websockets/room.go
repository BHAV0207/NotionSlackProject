package websockets

type Room struct {
	Name       string
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
			// fmt.Printf("✅ %s joined room %s\n", client.ID, r.Name)

		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
				// fmt.Printf("❌ %s left room %s\n", client.ID, r.Name)
			}

		case message := <-r.Broadcast:
			for c := range r.Clients {
				select {
				case c.Send <- message:
				default:
					close(c.Send)
					delete(r.Clients, c)
				}
			}
		}
	}
}
