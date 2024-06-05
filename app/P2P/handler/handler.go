package handler

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/messages"
	"bisgo/app/models"
	"bisgo/common"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type P2PHandler struct {
	*core.Core
}

func CreateP2PHandler(core *core.Core) *P2PHandler {
	return &P2PHandler{core}
}

func (h *P2PHandler) CreateTransaction(messageID int, payload []byte) {
	_, err := messages.LoadChannel(messageID)

	if err != nil {
		// handle error
	}

	defer messages.RemoveChannel(messageID)

	// handler logic

	var messageData common.TransactionDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		log.Println(err.Error())
		return
	}

	transactionType, err := strconv.Atoi(messageData.TransactionType)
	if err != nil {
		log.Println(err.Error())
		return
	}

	transaction := models.NewTransaction{
		Id:                "",
		OriginatorBankId:  messageData.OriginatorBankGlobalIdentifier,
		BeneficiaryBankId: messageData.BeneficiaryBankGlobalIdentifier,
		SenderId:          messageData.SenderName,
		ReceiverId:        messageData.ReceiverName,
		Currency:          messageData.Currency,
		Amount:            int(messageData.Amount),
		TransactionTypeId: transactionType,
		LoanId:            int(messageData.LoanID),
	}

	transactionID := h.DB.InsertTransaction(transaction)
	h.DB.UpdateTransactionState(transactionID, 1)

	// TODO: Add policies that are applicable to the created tx
}

func (h *P2PHandler) GetPolicies(messageID int, payload []byte) {
	_, err := messages.LoadChannel(messageID)

	if err != nil {
		// handle error
	}

	defer messages.RemoveChannel(messageID)

	var messageData common.PolicyRequestDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		log.Println(err.Error())
		return
	}

	transactionType, err := strconv.Atoi(messageData.TransactionType)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Comercial bank has to update polices from the central bank first
	// TODO: Add reference to this
	/*
		isCentralBank := true
		fmt.Println("central bank")
		if !isCentralBank {
			centralBankRequest := common.PolicyRequestDTO{
				Country:                   messageData.Country,
				TransactionType:           messageData.TransactionType,
				RequesterGlobalIdentifier: "myGlobalIdentifier", // TODO: Add reference to this
			}

			// TODO: Add reference to this
			channel, err := h.P2PClient.Send("myCentralBankIdentifier", "get-policies", centralBankRequest, 0)
			if err != nil {
				log.Println(err.Error())
				return
			}

			responseData := (<-channel).(common.PolicyResponseDTO)
			fmt.Println(responseData)

			// TODO: Update in DB if needed
		}
	*/

	policies := h.DB.GetPolicesByCountryCode(messageData.Country, transactionType)
	fmt.Println(policies)

	response := common.PolicyResponseDTO{
		Policies: []common.PolicyDTO{},
	}

	for _, policy := range policies {
		response.Policies = append(response.Policies, common.PolicyDTO{
			Code:   policy.PolicyType.Code,
			Name:   policy.PolicyType.Name,
			Params: policy.Policy.Parameters,
		})
	}

	_, err = h.P2PClient.Send(messageData.RequesterGlobalIdentifier, "send-policies", response, messageID)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (h *P2PHandler) SendPolicies(messageID int, payload []byte) {
	channel, err := messages.LoadChannel(messageID)

	if err != nil {
		// handle error
	}

	defer messages.RemoveChannel(messageID)

	var messageData common.PolicyResponseDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		log.Println(err.Error())
		return
	}

	channel <- messageData // send data to the listener
}

func (h *P2PHandler) CheckConfirmed(messageID int, payload []byte) {
	_, err := messages.LoadChannel(messageID)

	if err != nil {
		// handle error
	}

	defer messages.RemoveChannel(messageID)

	var messageData common.CheckConfirmedDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		log.Println(err.Error())
		return
	}

	applicablePolicies := h.DB.GetPoliciesForTransaction(messageData.CheckID)
	check := h.DB.GetComplianceCheckByID(messageData.CheckID)
	sender := h.DB.GetClientNameByID(check.SenderId)
	for _, policy := range applicablePolicies {
		if policy.PolicyType.Code == "SCL" {
			err = h.ProvingClient.SendProofRequest("interactive", messageData.CheckID, policy.Policy.Id, sender, messageData.VMAddress)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}
