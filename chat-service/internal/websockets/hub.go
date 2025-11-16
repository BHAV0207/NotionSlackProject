package websockets

import (
	"fmt"
	"sync"
)

type Hub struct {
	Rooms map[string]*Room
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{Rooms: make(map[string]*Room)}
}

func (h *Hub) GetRoom(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Rooms[name]
	if !exists {
		room = NewRoom(name)
		h.Rooms[name] = room
		go room.Run()
		fmt.Println("🏠 Created new room:", name)
	}
	return room
}