package client

import (
	"bisgo/app/DB"
	"bisgo/app/manager"
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ProvingClient struct {
	gpjcApiAddress       string
	gpjcApiPort          string
	gpjcPort             string
	db                   *DB.DBHandler
	sanctionsListManager *manager.SanctionListManager
}

var client ProvingClient

func init() {
	client = ProvingClient{
		gpjcApiAddress:       config.ResolveGpjcApiAddress(),
		gpjcApiPort:          config.ResolveGpjcApiPort(),
		gpjcPort:             config.ResolveGpjcPort(),
		db:                   DB.GetDBHandler(),
		sanctionsListManager: manager.CreateSanctionListManager(),
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
		err = c.sendNonInteractiveProofRequest(complianceCheckId, policyId)
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

func (c *ProvingClient) sendNonInteractiveProofRequest(complianceCheckId string, policyId int) error {
	returnErr := errors.New("unsucessful send of noninteractive proof request")

	complianceCheck := c.db.GetComplianceCheckByID(complianceCheckId)

	sender, err := c.db.GetBankClientById(int(complianceCheck.SenderId))
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	recever, err := c.db.GetBankClientById(int(complianceCheck.ReceiverId))
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	beneficiaryBank, err := c.db.GetBankByGlobalIdentifier(complianceCheck.BeneficiaryBankId)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	var participantsList [][]int = make([][]int, 0, 3)
	participantsList = append(participantsList, common.HashName(sender.Name))
	participantsList = append(participantsList, common.HashName(recever.Name))
	participantsList = append(participantsList, common.HashName(beneficiaryBank.Name))

	publicSanctionsList, err := c.sanctionsListManager.LoadSanctionListForNoninteractiveCheck()
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	moq := true
	if moq {
		publicSanctionsList = [][]int{{1, 2, 3}, {4, 5, 6}}
	}

	data := models.NonInteractiveComplianceCheckProofRequest{
		ComplianceCheckId:   complianceCheckId,
		PolicyId:            fmt.Sprintf("%d", policyId),
		ParticipantsList:    participantsList,
		PublicSanctionsList: publicSanctionsList,
	}

	fmt.Println(data)

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", config.ResolveNonInteractiveAPIAddress(), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		errlog.Println(err)
		return err
	}

	fmt.Println(string(body))

	return nil
}
