package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, you should validate the origin
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// func ServerWs(hub *Hub, w http.ResponseWriter, r *http.Request, DocId string) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("WebSocket upgrade error for document %s: %v", DocId, err)
// 		return
// 	}

// 	log.Printf("WebSocket connection established for document: %s", DocId)

// 	client := &Client{
// 		Hub:   hub,
// 		Conn:  conn,
// 		DocID: DocId,
// 	}

// 	hub.register <- client
// 	go client.ReadPump()
// }
