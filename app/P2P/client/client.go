package client

import (
	"bisgo/app/P2P/messages"
	manager "bisgo/app/P2P/peers_manager"
	"bisgo/common"
	"bisgo/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (c *P2PClient) Send(receivingBankID string, method string, data any, messageID int) (<-chan any, error) {
	receivingBankPeerID, err := manager.GetBankPeerID(receivingBankID)

	if err != nil {
		return nil, err
	}

	if messageID == 0 {
		messageID = rand.Int()
	}

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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, handleError(body)
	}

	channel := make(chan any, 1)

	messages.StoreChannel(messageID, channel)

	return channel, nil
}

func handleError(responseBody []byte) error {
	var resp common.ErrorResponse
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return err
	}
	return fmt.Errorf("call to p2p failed with error: " + resp.Message)
}
