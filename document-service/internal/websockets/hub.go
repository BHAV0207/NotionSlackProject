package websockets

import "github.com/gorilla/websocket"

type Hub struct {
	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			room, ok := h.rooms[client.DocID]
			if !ok {
				room = &Room{
					ID:      client.DocID,
					Clients: make(map[*Client]bool),
				}
				h.rooms[client.DocID] = room
			}
			room.Clients[client] = true

		case client := <-h.unregister:
			room, ok := h.rooms[client.DocID]
			if ok {
				delete(room.Clients, client)
			}
			if len(room.Clients) == 0 {
				delete(h.rooms, client.DocID)
			}

		case msg := <-h.broadcast:
			room, ok := h.rooms[msg.DocumentID]
			if !ok {
				continue
			}
			for client := range room.Clients {
				if client == msg.Sender {
					continue
				}

				err := client.Conn.WriteMessage(websocket.BinaryMessage, msg.Data)
				if err != nil {
					h.unregister <- client
				}
			}
		}
	}
}
