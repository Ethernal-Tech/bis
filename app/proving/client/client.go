package client

import (
	"bisgo/app/DB"
	"bisgo/config"
	"bisgo/errlog"
	"bytes"
	"errors"
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

func (c *ProvingClient) SendProofRequest(proofType string, complianceCheckId string, policyId int, targetedServer string) error {
	var err error

	if proofType == "interactive" {
		err = c.sendInteractiveProofRequest(complianceCheckId, policyId, targetedServer)
	} else {
		// non-interactive ...
	}

	if err != nil {
		errlog.Println(err)
		return errors.New("unsucessful send of proof request")
	}

	return nil
}

func (c *ProvingClient) sendInteractiveProofRequest(complianceCheckId string, policyId int, targetedServer string) error {
	returnErr := errors.New("unsucessful send of interactive proof request")

	gpjcApiURL := strings.Join([]string{c.gpjcApiAddress, c.gpjcApiPort}, ":")
	targetUrl := strings.Join([]string{"http:/", gpjcApiURL}, "/")

	complianceCheck, err := c.db.GetComplianceCheckById(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	originator, err := c.db.GetBankClientById(complianceCheck.SenderId)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	beneficiary, err := c.db.GetBankClientById(complianceCheck.ReceiverId)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	var payload []byte
	if targetedServer == "" {
		targetUrl = strings.Join([]string{targetUrl, "api", "start-server"}, "/")
		payload = []byte(fmt.Sprintf(`{"compliance_check_id": "%s", "policy_id": "%d"}`, complianceCheckId, policyId))
	} else {
		targetUrl = strings.Join([]string{targetUrl, "api", "start-client"}, "/")
		payload = []byte(fmt.Sprintf(`{"compliance_check_id": "%s", "policy_id": "%d", "participants": ["%s", "%s"], "to": "%s"}`, complianceCheckId, policyId, originator.Name, beneficiary.Name, targetedServer))
	}

	request, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(payload))
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "close")

	var client http.Client

	response, err := client.Do(request)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	if response.StatusCode != 200 {
		errlog.Println(errors.New("gpjc not started successfully"))
		return returnErr
	}

	return nil
}
