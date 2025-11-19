package websockets

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub   *Hub
	Conn  *websocket.Conn
	DocID string
}

func (c *Client) ReadPump() {
	defer func() {
		log.Printf("Client disconnected from document: %s", c.DocID)
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for document %s: %v", c.DocID, err)
			}
			break
		}

		// Only process binary messages (Yjs updates)
		if messageType == websocket.BinaryMessage {
			// Broadcast CRDT update to others
			c.Hub.broadcast <- Message{
				DocumentID: c.DocID,
				Data:       data,
				Sender:     c,
			}
		} else {
			log.Printf("Received non-binary message type %d for document %s, ignoring", messageType, c.DocID)
		}
	}
}
