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
	"time"
)

func (controller *APIController) GetPolicies(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	time.Sleep(4 * time.Second)

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
	originatorBankCountryCode := controller.DB.GetCountryOfBank(bankIdOfLoggedUser).Code

	// TODO: Handle RequesterGlobalIdentifier
	policyRequestDto := common.PolicyRequestDTO{
		Country:                   originatorBankCountryCode,
		TransactionType:           data.TransactionTypeId,
		RequesterGlobalIdentifier: bankIdOfLoggedUser,
	}

	ch, err := controller.P2PClient.Send(data.BeneficiaryBankId, "get-policies", policyRequestDto, 0)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	responseData := (<-ch).(common.PolicyResponseDTO)

	fmt.Println("Received from request")
	fmt.Println(responseData)

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

	controller.DB.InsertTransactionProof(messageData.TransactionId, messageData.Value)

	if messageData.Value == "0" {
		policyId, _ := strconv.Atoi(messageData.PolicyId)
		controller.DB.UpdateTransactionPolicyStatus(messageData.TransactionId, policyId, 1)
		controller.DB.UpdateTransactionState(messageData.TransactionId, 4)
	} else {
		policyId, _ := strconv.Atoi(messageData.PolicyId)
		controller.DB.UpdateTransactionPolicyStatus(messageData.TransactionId, policyId, 2)
		controller.DB.UpdateTransactionState(messageData.TransactionId, 5)
	}

	policyStatuses := controller.DB.GetTransactionPolicyStatuses(messageData.TransactionId)

	check := true

	for _, status := range policyStatuses {
		if status.Status != 1 {
			check = false
		}
	}

	if check {
		controller.DB.UpdateTransactionState(messageData.TransactionId, 6)
		controller.DB.UpdateTransactionState(messageData.TransactionId, 7)
	} else {
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

	// TODO: Add new client if not exists

	transaction := models.NewTransaction{
		Id:                messageData.TransactionID,
		OriginatorBankId:  messageData.OriginatorBankGlobalIdentifier,
		BeneficiaryBankId: messageData.BeneficiaryBankGlobalIdentifier,
		SenderId:          messageData.SenderLei,
		ReceiverId:        messageData.ReceiverLei,
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
