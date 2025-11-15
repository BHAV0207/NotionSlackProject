package websockets

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/BHAV0207/chat-service/pkg/models"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) ReadPump(dbClient *mongo.Client) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	collection := dbClient.Database("chatdb").Collection("messages")

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var message models.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		message.Timestamp = time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = collection.InsertOne(ctx, message)
		cancel()

		if err != nil {
			log.Println("DB insert error:", err)
			continue
		}

		// Broadcast to other clients
		c.Hub.Broadcast <- msg

	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
