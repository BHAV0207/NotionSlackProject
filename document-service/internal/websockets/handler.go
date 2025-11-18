package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServerWs(hub *Hub, w http.ResponseWriter, r *http.Request, DocId string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Hub: hub,
		Conn: conn,
		DocID: DocId,
	}

	hub.register <- client
	go client.ReadPump()
}
