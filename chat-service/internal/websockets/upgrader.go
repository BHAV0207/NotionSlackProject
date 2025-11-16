package websockets

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// note: pass both broker and dbClient
func ServeWS(hub *Hub, broker *RedisBroker, dbClient *mongo.Client, w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general"
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		ID:     r.RemoteAddr,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
		Room:   roomName,
		Broker: broker,
		DB:     dbClient, // <--- set DB client
	}

	room := hub.GetRoom(roomName)
	room.Register <- client

	// subscribe once per room (subscribe is idempotent because hub.GetRoom creates room only once)
	if broker != nil {
		broker.Subscribe(roomName, hub)
	}

	go client.WritePump()
	go client.ReadPump() // no args now
}
