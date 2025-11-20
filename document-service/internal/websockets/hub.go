package websockets

import (
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			room, ok := h.Rooms[client.DocID]
			if !ok {
				room = &Room{
					ID:      client.DocID,
					Clients: make(map[*Client]bool),
				}
				h.Rooms[client.DocID] = room
			}
			room.Clients[client] = true
			log.Printf("Client Registered for document %s. Total clients: %d", client.DocID, len(room.Clients))

		case client := <-h.Unregister:
			room, ok := h.Rooms[client.DocID]
			if ok {
				delete(room.Clients, client)
				if len(room.Clients) == 0 {
					delete(h.Rooms, client.DocID)
					log.Printf("Room for document %s deleted (no clients)", client.DocID)
				} else {
					log.Printf("Client Unregistered from document %s. Remaining clients: %d", client.DocID, len(room.Clients))
				}
			}

		case msg := <-h.Broadcast:
			room, ok := h.Rooms[msg.DocumentID]
			if !ok {
				continue
			}
			for client := range room.Clients {
				if client == msg.Sender {
					continue
				}

				err := client.Conn.WriteMessage(websocket.BinaryMessage, msg.Data)
				if err != nil {
					log.Printf("Error Broadcasting to client in document %s: %v", msg.DocumentID, err)
					h.Unregister <- client
				}
			}
		}
	}
}
