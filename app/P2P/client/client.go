package client

import (
	"bisgo/app/P2P/messages"
	manager "bisgo/app/P2P/peers_manager"
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strings"
)

type P2PClient struct {
	p2pNodeAddress string
}

var client P2PClient

func init() {
	// configuration settings (e.g. for p2pNodeAddress)
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

	messagePayload, err := messages.CreateMessagePayload(messageID, method, data)

	if err != nil {
		return nil, err
	}

	message := struct {
		PeerID  string
		method  string
		Payload []byte
	}{receivingBankPeerID, method, messagePayload}

	messageJSON, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", strings.Join([]string{"http:/", c.p2pNodeAddress, "v1", "p2p", "passthru"}, "/"), bytes.NewBuffer(messageJSON))

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
