package controller

import (
	"bisgo/common"
	"bisgo/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

func (controller *APIController) GetPolicies(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	data := struct {
		BankId            string
		TransactionTypeId string
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	bankId, _ := strconv.Atoi(data.BankId)
	transactionTypeId, _ := strconv.Atoi(data.TransactionTypeId)

	policies := controller.DB.GetPolices(uint64(bankId), transactionTypeId)

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

func (controller *APIController) SubmitTransactionProof(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var messageData models.TransactionProofRequest
	if err := json.Unmarshal(body, &messageData); err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	transactionId, err := strconv.Atoi(messageData.TransactionId)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}
	controller.DB.InsertTransactionProof(uint64(transactionId), messageData.Value)

	if messageData.Value == "0" {
		policyId, _ := strconv.Atoi(messageData.PolicyId)
		controller.DB.UpdateTransactionPolicyStatus(uint64(transactionId), policyId, 1)
		controller.DB.UpdateTransactionState(uint64(transactionId), 4)
	} else {
		policyId, _ := strconv.Atoi(messageData.PolicyId)
		controller.DB.UpdateTransactionPolicyStatus(uint64(transactionId), policyId, 2)
		controller.DB.UpdateTransactionState(uint64(transactionId), 5)
	}

	policyStatuses := controller.DB.GetTransactionPolicyStatuses(uint64(transactionId))

	check := true

	for _, status := range policyStatuses {
		if status.Status != 1 {
			check = false
		}
	}

	if check {
		controller.DB.UpdateTransactionState(uint64(transactionId), 6)
		controller.DB.UpdateTransactionState(uint64(transactionId), 7)
	} else {
		controller.DB.UpdateTransactionState(uint64(transactionId), 8)
	}

	err = json.NewEncoder(w).Encode("Ok")
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

	policies := controller.DB.GetPolicy(data.BankCountry, uint64(policyId))

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

	transaction := models.Transaction{
		OriginatorBank:  controller.DB.GetBankIdByIdentifier(messageData.OriginatorBankGlobalIdentifier),
		BeneficiaryBank: controller.DB.GetBankIdByIdentifier(messageData.BeneficiaryBankGlobalIdentifier),
		Sender:          controller.DB.GetBankClientId(messageData.SenderName),
		Receiver:        controller.DB.GetBankClientId(messageData.ReceiverName),
		Currency:        messageData.Currency,
		Amount:          int(messageData.Amount),
		TypeId:          transactionType,
		LoanId:          int(messageData.LoanID),
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
