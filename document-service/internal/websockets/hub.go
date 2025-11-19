package websockets

import (
	"log"

	"github.com/gorilla/websocket"
)

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
			log.Printf("Client registered for document %s. Total clients: %d", client.DocID, len(room.Clients))

		case client := <-h.unregister:
			room, ok := h.rooms[client.DocID]
			if ok {
				delete(room.Clients, client)
				if len(room.Clients) == 0 {
					delete(h.rooms, client.DocID)
					log.Printf("Room for document %s deleted (no clients)", client.DocID)
				} else {
					log.Printf("Client unregistered from document %s. Remaining clients: %d", client.DocID, len(room.Clients))
				}
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
					log.Printf("Error broadcasting to client in document %s: %v", msg.DocumentID, err)
					h.unregister <- client
				}
			}
		}
	}
}
