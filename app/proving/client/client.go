package client

import (
	"bisgo/app/DB"
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
	db             *DB.DBHandler
}

var client ProvingClient

func init() {
	client = ProvingClient{
		gpjcApiAddress: config.ResolveGpjcApiAddress(),
		gpjcApiPort:    config.ResolveGpjcApiPort(),
		gpjcPort:       config.ResolveGpjcPort(),
		db:             DB.GetDBHandler(),
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

func (c *ProvingClient) SendProofRequest(proofType string, checkID string, policyID int, targetedServer string /* additional parameters */) error {

	// logic for sending proof requests
	// different logic based on whether the request is for interactive or non-interactive proof
	if proofType == "interactive" {
		return c.sendInteractiveProofRequest(checkID, policyID, targetedServer)
	}

	return nil
}

func (c *ProvingClient) sendInteractiveProofRequest(checkID string, policyID int, targetedServer string) error {
	gpjcApiURL := strings.Join([]string{c.gpjcApiAddress, c.gpjcApiPort}, ":")
	targetUrl := strings.Join([]string{"http:/", gpjcApiURL}, "/")

	check := c.db.GetComplianceCheckByID(checkID)
	senderName := c.db.GetClientNameByID(check.SenderId)
	receiverName := c.db.GetClientNameByID(check.ReceiverId)

	var payload []byte
	if targetedServer == "" {
		targetUrl = strings.Join([]string{targetUrl, "api", "start-server"}, "/")
		payload = []byte(fmt.Sprintf(`{"compliance_check_id": "%s", "policy_id": "%d"}`, checkID, policyID))
	} else {
		targetUrl = strings.Join([]string{targetUrl, "api", "start-client"}, "/")
		payload = []byte(fmt.Sprintf(`{"compliance_check_id": "%s", "policy_id": "%d", "participants": ["%s", "%s"], "to": "%s"}`, checkID, policyID, senderName, receiverName, targetedServer))
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("gpjc not started successfully")
	}

	return nil
}
