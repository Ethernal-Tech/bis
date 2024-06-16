package handler

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/messages"
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/json"
	"errors"
	"strconv"
)

type P2PHandler struct {
	*core.Core
}

func CreateP2PHandler(core *core.Core) *P2PHandler {
	return &P2PHandler{core}
}

// TODO: if an error occurs, the channel should be closed so that the listener does not wait forever

// TODO: all or nothing - methods should be constructed as atomic blocks (transactions), if anything fails, all changes are rolled back

func (h *P2PHandler) CreateTransaction(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method CreateTransaction failed to execute properly")

	var messageData common.TransactionDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		errlog.Println(err)
		return returnErr
	}

	transactionType, err := strconv.Atoi(messageData.TransactionType)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	senderID := h.DB.GetOrCreateClient(messageData.SenderLei, messageData.SenderName, "", messageData.OriginatorBankGlobalIdentifier)
	receiverID := h.DB.GetOrCreateClient(messageData.ReceiverLei, messageData.ReceiverName, "", messageData.BeneficiaryBankGlobalIdentifier)

	transaction := models.NewTransaction{
		Id:                messageData.TransactionID,
		OriginatorBankId:  messageData.OriginatorBankGlobalIdentifier,
		BeneficiaryBankId: messageData.BeneficiaryBankGlobalIdentifier,
		SenderId:          senderID,
		ReceiverId:        receiverID,
		Currency:          messageData.Currency,
		Amount:            int(messageData.Amount),
		TransactionTypeId: transactionType,
		LoanId:            int(messageData.LoanID),
	}

	transactionID := h.DB.InsertTransaction(transaction)
	h.DB.UpdateTransactionState(transactionID, 1)

	if config.ResovleIsCentralBank() {
		h.RulesEngine.Do(messageData.TransactionID, "interactive", nil)
	}

	return nil
}

func (h *P2PHandler) GetPolicies(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method GetPolicies failed to execute properly")

	var messageData common.PolicyRequestDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		errlog.Println(err)
		return returnErr
	}

	transactionType, err := strconv.Atoi(messageData.TransactionType)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	// Comercial bank has to update polices from the central bank first
	// TODO: Add reference to this
	if !config.ResovleIsCentralBank() {
		centralBankRequest := common.PolicyRequestDTO{
			Country:                   messageData.Country,
			TransactionType:           messageData.TransactionType,
			RequesterGlobalIdentifier: config.ResolveMyGlobalIdentifier(), // TODO: Add reference to this
		}

		// TODO: Add reference to this
		channel, err := h.P2PClient.Send(config.ResolveCBGlobalIdentifier(), "get-policies", centralBankRequest, 0)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		responseData := (<-channel).(common.PolicyResponseDTO)

		for _, policy := range responseData.Policies {
			policyTypeID := h.DB.GetOrCreatePolicyType(policy.Code, policy.Name)
			h.DB.GetOrCreatePolicy(int(policyTypeID), transactionType, h.DB.GetCountryByCode(config.ResolveCountryCode()).Id, h.DB.GetCountryByCode(messageData.Country).Id, policy.Params)
		}
	}

	policies := h.DB.GetPolicesByCountryCode(messageData.Country, transactionType)

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
		errlog.Println(err)
		return returnErr
	}

	return nil
}

func (h *P2PHandler) SendPolicies(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method SendPolicies failed to execute properly")

	channel, _ := messages.LoadChannel(messageID)

	defer messages.RemoveChannel(messageID)

	var messageData common.PolicyResponseDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		errlog.Println(err)
		return returnErr
	}

	channel <- messageData

	return nil
}

func (h *P2PHandler) CheckConfirmed(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method CheckConfirmed failed to execute properly")

	var messageData common.CheckConfirmedDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		errlog.Println(err)
		return returnErr
	}

	data := map[string]any{
		"vm_address": messageData.VMAddress,
	}
	h.RulesEngine.Do(messageData.CheckID, "interactive", data)

	return nil
}

func (h *P2PHandler) CFMResultBeneficiary(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method CFMResultBeneficiary failed to execute properly")

	var messageData common.CFMCheckDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		errlog.Println(err)
		return returnErr
	}

	applicablePolicies := h.DB.GetPoliciesForTransaction(messageData.TransctionID)
	for _, policy := range applicablePolicies {
		if policy.PolicyType.Code == "CFM" {
			h.DB.UpdateTransactionPolicyStatus(messageData.TransctionID, policy.Policy.Id, messageData.Result)
			if messageData.Result == 2 {
				h.DB.UpdateTransactionState(messageData.TransctionID, 8)
			} else {
				h.DB.UpdateTransactionState(messageData.TransctionID, 7)
			}
		}
	}

	check := h.DB.GetComplianceCheckByID(messageData.TransctionID)

	_, err := h.P2PClient.Send(check.OriginatorBankId, "cfm-result-originator", any(messageData), 0)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	return nil
}

func (h *P2PHandler) CFMResultOriginator(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method CFMResultOriginator failed to execute properly")

	var messageData common.CFMCheckDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		errlog.Println(err)
		return returnErr
	}

	applicablePolicies := h.DB.GetPoliciesForTransaction(messageData.TransctionID)
	for _, policy := range applicablePolicies {
		if policy.PolicyType.Code == "CFM" {
			h.DB.UpdateTransactionPolicyStatus(messageData.TransctionID, policy.Policy.Id, messageData.Result)
			if messageData.Result == 2 {
				h.DB.UpdateTransactionState(messageData.TransctionID, 8)
			} else {
				h.DB.UpdateTransactionState(messageData.TransctionID, 7)
			}
		}
	}

	return nil
}
