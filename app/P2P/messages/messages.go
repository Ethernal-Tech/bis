package messages

// client message
type P2PClientMessage struct {
	PeerID    string `json:"peer_id"`
	MessageID uint64 `json:"message_id"`
	Method    string `json:"method"`
	Payload   []byte `json:"payload"`
}

// server message
type P2PServerMessasge struct {
	MessageID int    `json:"message_id"`
	Method    string `json:"method"`
	Payload   []byte `json:"payload"`
}
