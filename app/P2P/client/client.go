package client

import (
	"bisgo/app/P2P/messages"
	manager "bisgo/app/P2P/peers_manager"
	"bisgo/config"
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
)

type P2PClient struct {
	p2pNodeAddress string
}

var client P2PClient

func init() {
	client = P2PClient{config.ResolveP2PNodeAPIAddress()}
}

func GetP2PClient() *P2PClient {
	return &client
}

func (c *P2PClient) Send(receivingBankID string, method string, data any) (<-chan any, error) {
	receivingBankPeerID, err := manager.GetBankPeerID(receivingBankID)

	if err != nil {
		return nil, err
	}

	messageID := rand.Int()

	messagePayload, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	messageData := messages.P2PClientMessage{
		PeerID:    receivingBankPeerID,
		MessageID: messageID,
		Method:    method,
		Payload:   messagePayload,
	}

	message, err := json.Marshal(messageData)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", c.p2pNodeAddress, bytes.NewBuffer(message))

	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("error: unsuccessful request to the peer")
	}

	channel := make(chan any, 1)

	messages.StoreChannel(messageID, channel)

	return channel, nil
}
