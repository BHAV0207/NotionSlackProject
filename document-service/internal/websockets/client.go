package websockets

import "github.com/gorilla/websocket"

type Client struct {
	Hub   *Hub
	Conn  *websocket.Conn
	DocID string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// Broadcast CRDT update to others
		c.Hub.broadcast <- Message{
			DocumentID: c.DocID,
			Data:       data,
			Sender:     c,
		}
	}
}
