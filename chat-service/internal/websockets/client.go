package websockets

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

// IMP => message form browser of the sending client goes to => readPump() => from read pump to Hub via the Broadcast channel , then =>   the message received by the hub is send to al the clients in the room via a for loop whihc have the send channels of all the clinets that is c.send <- message //

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

		var m models.Message
		if err := json.Unmarshal(msgBytes, &m); err != nil {
			log.Println("invalid message format:", err)
			continue
		}

		switch m.Type {
		case "suggested_replies":
			go c.handleSuggestReplies(m.Content)
			continue // do NOT save to Mongo, do NOT broadcast
		case "send_message", "":
			if m.RoomID == "" {
				m.RoomID = c.Room // fallback to the client's room
			}
			if m.Timestamp.IsZero() {
				m.Timestamp = time.Now().UTC()
			}

			// 2) Persist to Mongo (originating instance only)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, insertErr := collection.InsertOne(ctx, m)
			cancel()
			if insertErr != nil {
				log.Println("mongo insert error:", insertErr)
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
}

func (c *Client) handleSuggestReplies(text string) {
	suggestions, err := callAISuggest(text)
	if err != nil {
		log.Println("AI error:", err)
		return
	}

	response := map[string]interface{}{
		"type":        "ai_suggestions",
		"suggestions": suggestions,
	}

	data, _ := json.Marshal(response)
	c.Send <- data
}

func callAISuggest(text string) ([]string, error) {
	payload := map[string]string{
		"context": text,
	}

	b, _ := json.Marshal(payload)

	resp, err := http.Post(
		"http://ai-service:8005/suggest",
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Suggestions string `json:"suggestions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// AI returns a string. Convert to a slice.
	return strings.Split(result.Suggestions, "\n"), nil
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
