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
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	Room   string
	Broker *RedisBroker
	DB     *mongo.Client // <--- add Mongo client here
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.GetRoom(c.Room).Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	if c.DB == nil {
		log.Println("DB client is nil")
		return
	}
	if c.Broker == nil {
		log.Println("Broker is nil")
		return
	}

	collection := c.DB.Database("chatdb").Collection("messages")

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// 1) Parse incoming JSON into Message
		var m models.Message
		if err := json.Unmarshal(msgBytes, &m); err != nil {
			log.Println("invalid message format:", err)
			continue
		}

		// Ensure required fields
		if m.RoomID == "" {
			m.RoomID = c.Room // fallback to the client's room
		}
		if m.Timestamp.IsZero() {
			m.Timestamp = time.Now().UTC()
		}
		// If you have authentication, set SenderID from auth, not client-sent field:
		// m.SenderID = c.UserID

		// 2) Persist to Mongo (originating instance only)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, insertErr := collection.InsertOne(ctx, m)
		cancel()
		if insertErr != nil {
			log.Println("mongo insert error:", insertErr)
			// Continue: we can still publish even if DB write failed, or choose to skip publish.
		}

		// 3) Publish to Redis so other instances receive it
		// Ensure message payload is JSON string
		payload, _ := json.Marshal(m)
		c.Broker.Publish(m.RoomID, payload)

		// Optionally also broadcast locally immediately (can skip since Redis subscription will deliver
		// back to this instance if you subscribe to room channel). But to reduce latency you may want both:
		room := c.Hub.GetRoom(m.RoomID)
		room.Broadcast <- payload
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
