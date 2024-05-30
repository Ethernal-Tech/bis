package messages

import (
	"bisgo/app/models"
	"encoding/json"
	"errors"
)

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

func CreateMessagePayload(messageID int, method string, data any) ([]byte, error) {
	switch method {
	case "method1":
		tx, ok := data.(models.Transaction)

		if !ok {
			return nil, errors.New("error: can't create message payload")
		}

		return json.Marshal(tx)
	case "method2":
		return nil, errors.New("error: can't create message payload")
	default:
		return nil, errors.New("error: can't create message payload")
	}
}
