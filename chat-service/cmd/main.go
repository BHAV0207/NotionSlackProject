package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

/*
WebSocket connections are full-duplex, meaning reading and writing happen independently. But the websocket.Conn type in Go is not safe for concurrent writes — if multiple goroutines try to write at once, it’ll panic or corrupt data.
So we use:
a read loop (one goroutine) that listens for messages from the client and sends them to the hub.
a write loop (another goroutine) that listens on a channel (client.Send) and writes to the socket.
That way, no two goroutines ever write directly to the socket simultaneously.
*/

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("✅ New WebSocket connection established!")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reding message: %v", err)
			break
		}
		fmt.Printf("receive: %s\n", msg)

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", HandleConnection)
	fmt.Println("🚀 WebSocket server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
