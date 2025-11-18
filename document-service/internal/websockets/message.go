package websockets

type Message struct {
	DocumentID string
	Data       []byte
	Sender     *Client
}
