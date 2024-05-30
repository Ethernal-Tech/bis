package messaging

import (
	"bisgo/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MessagingHandler provides functionalities for communicating with Peer-to-Peer node.
type MessagingHandler struct {
	P2PNodeAddress string
}

// CreateMessagingHandler initializes and returns a pointer to a new MessagingHandler instance.
// Parameters:
//   - p2pNodeAddress: a string representing the address of the Peer-to-Peer node.
func CreateMessagingHandler(p2pNodeAddress string) *MessagingHandler {
	return &MessagingHandler{
		P2PNodeAddress: p2pNodeAddress,
	}
}

// GetAvailablePeers calls the Peer-to-Peer node to discover available peers in the system.
func (m *MessagingHandler) GetAvailablePeers() ([]common.Peer, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", strings.Join([]string{m.P2PNodeAddress, "v1", "p2p", "peers"}, "/"), nil)
	if err != nil {
		return []common.Peer{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []common.Peer{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []common.Peer{}, err
	}

	if resp.StatusCode != 200 {
		return []common.Peer{}, handleError(body)
	}

	var ret []common.Peer
	err = json.Unmarshal(body, &ret)

	return ret, err
}

// SendPassthruMessage sends the message to the receiveing bank over Peer-to-Peer network.
// Parameters:
//   - receiveingBankPeerId: a string representing the id of the targeted Peer-to-Peer node.
//   - receiveingBankURL: a string representing the url of the targeted bank.
//   - method: a string representing the targeted method of the targeted bank.
//   - requestData: a byte array representing the data that is being sent.
func (m *MessagingHandler) SendPassthruMessage(receiveingBankPeerId string, receiveingBankURL string, method string, requestData []byte) error {
	reqObj := common.PassThruRequest{
		PeerID:  receiveingBankPeerId,
		URI:     strings.Join([]string{"http:/", receiveingBankURL, "api", method}, "/"),
		Payload: requestData,
	}

	reqBytes, err := json.Marshal(reqObj)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", strings.Join([]string{"http:/", m.P2PNodeAddress, "v1", "p2p", "passthru"}, "/"), bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return handleError(body)
	}

	return nil
}

func handleError(responseBody []byte) error {
	var resp common.ErrorResponse
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return err
	}

	return fmt.Errorf("call to p2p failed with error: " + resp.Message)
}
