package client

import (
	"bisgo/config"
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type ProvingClient struct {
	gpjcApiAddress string
	gpjcApiPort    string
	gpjcPort       string
}

var client ProvingClient

func init() {
	client = ProvingClient{
		gpjcApiAddress: config.ResolveGpjcApiAddress(),
		gpjcApiPort:    config.ResolveGpjcApiPort(),
		gpjcPort:       config.ResolveGpjcPort(),
	}
}

func GetProvingClient() *ProvingClient {
	return &client
}

func (c *ProvingClient) GetGpjcApiAddress() string {
	return c.gpjcApiAddress
}

func (c *ProvingClient) GetVMAddress() string {
	return strings.Join([]string{c.gpjcApiAddress, c.gpjcPort}, ":")
}

func (c *ProvingClient) SendProofRequest(proofType string, checkID string, policyID int, receiverName string, targetedServer string /* additional parameters */) error {

	// logic for sending proof requests
	// different logic based on whether the request is for interactive or non-interactive proof
	if proofType == "interactive" {
		gpjcApiURL := strings.Join([]string{c.gpjcApiAddress, c.gpjcApiPort}, ":")
		return sendInteractiveProofRequest(gpjcApiURL, checkID, policyID, receiverName, targetedServer)
	}

	return nil
}

func sendInteractiveProofRequest(address string, checkID string, policyID int, receiverName string, targetedServer string) error {
	targetUrl := strings.Join([]string{"http:/", address}, "/")
	var payload []byte
	if targetedServer == "" {
		targetUrl = strings.Join([]string{targetUrl, "start-server"}, "/")
		payload = []byte(fmt.Sprintf(`{"tx_id": "%s", "policy_id": "%d"}`, checkID, policyID))
	} else {
		targetUrl = strings.Join([]string{targetUrl, "start-client"}, "/")
		payload = []byte(fmt.Sprintf(`{"tx_id": "%s", "policy_id": "%d", "receiver": "%s", "to": "%s:10501"}`, checkID, policyID, receiverName, targetedServer))
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
