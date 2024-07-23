package client

import (
	"bisgo/app/P2P/messages"
	manager "bisgo/app/P2P/peers_manager"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"bytes"
	"encoding/json"
	"errors"
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
		errlog.Println(err)
		return nil, errors.New("unknown receiving bank ID")
	}

	if messageID == 0 {
		messageID = rand.Int()
	}

	bytePayload, err := json.Marshal(data)
	if err != nil {
		errlog.Println(fmt.Errorf("%v %w", data, err))
		return nil, errors.New("serialization failed due to malformed input data")
	}

	// TODO: describe why it is needed
	messagePayload := make([]int, len(bytePayload))
	for i, v := range bytePayload {
		messagePayload[i] = int(v)
	}

	messageData := messages.P2PClientMessage{
		PeerID:    receivingBankPeerID,
		MessageID: uint64(messageID),
		Method:    method,
		Payload:   messagePayload,
	}

	message, err := json.Marshal(messageData)
	if err != nil {
		errlog.Println(fmt.Errorf("%v %w", messageData, err))
		return nil, errors.New("cannot create p2p message")
	}

	client := &http.Client{}

	// TODO: implement a timeout on a p2p node's response (context)

	request, err := http.NewRequest("POST", c.p2pNodeAddress, bytes.NewBuffer(message))
	if err != nil {
		errlog.Println(err)
		return nil, errors.New("cannot create a p2p wrapping (HTTP POST) message due to internal system problems; not caused by incorrect parameters")
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		errlog.Println(err)
		return nil, errors.New("cannot send a p2p wrapping (HTTP POST) message due to internal system problems; not caused by incorrect parameters")
	}

	if response.StatusCode != 200 {
		returnErr := errors.New("p2p node rejected the message")
		body, err := io.ReadAll(response.Body)
		if err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		var response common.ErrorResponse
		if err := json.Unmarshal(body, &response); err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		errlog.Println(errors.New("p2p node response: " + response.Message))

		return nil, returnErr
	}

	channel := make(chan any, 1)

	messages.StoreChannel(messageID, channel)

	return channel, nil
}
