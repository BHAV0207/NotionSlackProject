package websockets

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

// Upgrader used to upgrade HTTP to WebSocket. Allow all origins for now (dev only).
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub, dbClient *mongo.Client, w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump(dbClient)
}
