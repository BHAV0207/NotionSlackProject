package websockets

type Room struct {
	ID      string
	Clients map[*Client]bool
}
