package messages

// client message
type P2PClientMessage struct {
	PeerID    string `json:"peer_id"`
	MessageID int    `json:"message_id"`
	Method    string `json:"uri"`
	Payload   []byte `json:"payload"`
}

// server message
type P2PServerMessasge struct {
	MessageID int    `json:"message_id"`
	Method    string `json:"name"`
	Payload   []byte `json:"payload"`
}
