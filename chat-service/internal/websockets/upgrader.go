package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader used to upgrade HTTP to WebSocket. Allow all origins for now (dev only).
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
