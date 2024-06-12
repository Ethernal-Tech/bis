package controller

import (
	"bisgo/app/models"
	"bisgo/common"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (controller *APIController) GetPolicies(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	//time.Sleep(2 * time.Second)

	data := struct {
		BeneficiaryBankId string
		TransactionTypeId string
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	bankIdOfLoggedUser := controller.SessionManager.GetString(r.Context(), "bankId") //get logged user's bank ID
	originatorBankCountry := controller.DB.GetCountryOfBank(bankIdOfLoggedUser)

	policyRequestDto := common.PolicyRequestDTO{
		Country:                   originatorBankCountry.Code,
		TransactionType:           data.TransactionTypeId,
		RequesterGlobalIdentifier: bankIdOfLoggedUser,
	}

	ch, err := controller.P2PClient.Send(data.BeneficiaryBankId, "get-policies", policyRequestDto, 0)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	responseData := (<-ch).(common.PolicyResponseDTO)

	// Add policies to the DB if not exist
	for _, policy := range responseData.Policies {
		// 1. Insert policy type if not exists
		policyTypeID := controller.DB.GetOrCreatePolicyType(policy.Code, policy.Name)
		// 2. Insert Policy if not exists
		transactionTypeID, err := strconv.Atoi(data.TransactionTypeId)
		if err != nil {
			http.Error(w, fmt.Sprint("Internal Server Error %w", err), 500)
		}

		policyEnforcingCountryId := controller.DB.GetCountryOfBank(data.BeneficiaryBankId).Id
		controller.DB.GetOrCreatePolicy(int(policyTypeID), transactionTypeID, policyEnforcingCountryId, originatorBankCountry.Id, policy.Params)
	}

	jsonData, err := json.Marshal(responseData)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (controller *APIController) SubmitTransactionProof(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var messageData models.TransactionProofRequest
	if err := json.Unmarshal(body, &messageData); err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	values := strings.Split(messageData.Value, ";")

	controller.DB.InsertTransactionProof(messageData.TransactionId, values[1]+values[2])

	if values[0] == "0" {
		policyId, _ := strconv.Atoi(messageData.PolicyId)
		controller.DB.UpdateTransactionPolicyStatus(messageData.TransactionId, policyId, 1)
		controller.DB.UpdateTransactionState(messageData.TransactionId, 4)
	} else {
		policyId, _ := strconv.Atoi(messageData.PolicyId)
		controller.DB.UpdateTransactionPolicyStatus(messageData.TransactionId, policyId, 2)
		controller.DB.UpdateTransactionState(messageData.TransactionId, 5)
		controller.DB.UpdateTransactionState(messageData.TransactionId, 8)
	}

	err := json.NewEncoder(w).Encode("Ok")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
		return
	}
}

func (controller *APIController) GetPolicy(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	data := struct {
		BankCountry string
		PolicyId    string
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	policyId, err := strconv.Atoi(data.PolicyId)
	if err != nil {
		http.Error(w, "Failed to decode policy id", http.StatusInternalServerError)
		return
	}

	policies := controller.DB.GetPolicy(uint64(policyId))

	jsonData, err := json.Marshal(policies)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// ------------------------------------------------------------------------------------------
func (controller *APIController) CreateTx(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var messageData common.TransactionDTO
	if err := json.Unmarshal(body, &messageData); err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	transactionType, err := strconv.Atoi(messageData.TransactionType)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	senderID := controller.DB.GetOrCreateClient(messageData.SenderLei, messageData.SenderName, "", messageData.OriginatorBankGlobalIdentifier)
	receiverID := controller.DB.GetOrCreateClient(messageData.ReceiverLei, messageData.ReceiverName, "", messageData.BeneficiaryBankGlobalIdentifier)

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

	transactionID := controller.DB.InsertTransaction(transaction)
	controller.DB.UpdateTransactionState(transactionID, 1)

	err = json.NewEncoder(w).Encode("Ok")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
		return
	}
}
