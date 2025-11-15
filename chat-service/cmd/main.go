package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BHAV0207/chat-service/internal/websockets"
)

var hub = websockets.NewHub()

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websockets.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &websockets.Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func main() {
	go hub.Run()

	http.HandleFunc("/ws", serveWS)
	fmt.Println("🚀 Chat server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
