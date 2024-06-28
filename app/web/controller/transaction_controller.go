package controller

import (
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/config"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (controller *TransactionController) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	var searchModel models.SearchModel
	err := json.NewDecoder(r.Body).Decode(&searchModel)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	viewData := map[string]any{}
	var transactions []models.TransactionModel
	if controller.SessionManager.GetBool(r.Context(), "centralBankEmployee") {
		var countryId string
		transactions, countryId = controller.DB.GetCentralBankTransactions(controller.SessionManager.Get(r.Context(), "bankId").(string), searchModel)
		viewData["countryId"] = countryId
	} else {
		transactions = controller.DB.GetCommercialBankTransactions(controller.SessionManager.Get(r.Context(), "bankId").(string), searchModel)
	}

	viewData["bankName"] = controller.SessionManager.GetString(r.Context(), "bankName")
	viewData["transactions"] = transactions
	viewData["country"] = controller.SessionManager.GetString(r.Context(), "country")
	viewData["centralBankEmployee"] = controller.SessionManager.GetBool(r.Context(), "centralBankEmployee")

	ts, err := template.ParseFiles("./static/views/_transactionPartial.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	err = ts.Execute(w, viewData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (controller *TransactionController) AddTransaction(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodGet {

		viewData := map[string]any{}

		loanID := rand.Intn(2500000)

		viewData["loanId"] = loanID
		viewData["username"] = controller.SessionManager.GetString(r.Context(), "username")
		viewData["bankName"] = controller.SessionManager.GetString(r.Context(), "bankName")
		viewData["country"] = controller.SessionManager.GetString(r.Context(), "country")
		viewData["banks"] = controller.DB.GetBanks()
		viewData["transactionTypes"] = controller.DB.GetTransactionTypes()

		ts, err := template.ParseFiles("./app/web/static/views/addcompliance.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error 1", 500)
			return
		}

		err = ts.Execute(w, viewData)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error 2", 500)
		}

	} else if r.Method == http.MethodPost {

		data := struct {
			SenderLei         string `json:"senderLei"`
			SenderName        string `json:"senderName"`
			BeneficiaryLei    string `json:"beneficiaryLei"`
			BeneficiaryName   string `json:"beneficiaryName"`
			PaymentTypeID     string `json:"paymentType"`
			TransactionTypeID string `json:"transactionType"`
			Currency          string `json:"currency"`
			Amount            string `json:"amount"`
			BeneficiaryBank   string `json:"beneficiaryBank"`
		}{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&data); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		originatorBankId := controller.DB.GetBankId(controller.SessionManager.GetString(r.Context(), "bankName"))

		amount, _ := strconv.Atoi(strings.Replace(data.Amount, ",", "", -1))
		transactionType, _ := strconv.Atoi(data.TransactionTypeID)
		//loanId, _ := strconv.Atoi(data.PaymentTypeID)

		// If client is not in the DB add it
		senderID := controller.DB.GetOrCreateClient(data.SenderLei, data.SenderName, "", originatorBankId)
		receiverID := controller.DB.GetOrCreateClient(data.BeneficiaryLei, data.BeneficiaryName, "", data.BeneficiaryBank)

		transaction := models.NewTransaction{
			OriginatorBankId:  originatorBankId,
			BeneficiaryBankId: data.BeneficiaryBank,
			SenderId:          senderID,
			ReceiverId:        receiverID,
			Currency:          data.Currency,
			Amount:            amount,
			TransactionTypeId: transactionType,
			LoanId:            0,
		}

		transactionID := controller.DB.InsertTransaction(transaction)
		controller.DB.UpdateTransactionState(transactionID, 1)

		// Call P2P create-transaction of beneficiary bank
		transactionDto := common.TransactionDTO{
			TransactionID:                   transactionID,
			SenderLei:                       data.SenderLei,
			SenderName:                      data.SenderName,
			ReceiverLei:                     data.BeneficiaryLei,
			ReceiverName:                    data.BeneficiaryName,
			OriginatorBankGlobalIdentifier:  transaction.OriginatorBankId,
			BeneficiaryBankGlobalIdentifier: transaction.BeneficiaryBankId,
			PaymentType:                     "",
			TransactionType:                 fmt.Sprint(transactionType),
			Amount:                          uint64(amount),
			Currency:                        data.Currency,
			SwiftBICCode:                    "",
			LoanID:                          uint64(0),
		}

		ch, err := controller.P2PClient.Send(transactionDto.BeneficiaryBankGlobalIdentifier, "create-transaction", transactionDto, 0)

		_ = ch

		if err != nil {
			log.Println(err.Error())
			http.Error(w, fmt.Sprint("Internal Server Error sending create tx %w", err), 500)
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (controller *TransactionController) ConfirmTransaction(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodGet {

		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error parsing form", 500)
		}

		transaction := controller.DB.GetTransactionHistory(r.Form.Get("transaction"))

		policies := controller.DB.GetPoliciesForTransaction(transaction.Id)

		viewData := map[string]any{}

		viewData["username"] = controller.SessionManager.GetString(r.Context(), "username")
		viewData["transaction"] = transaction
		viewData["bankName"] = controller.SessionManager.GetString(r.Context(), "bankName")
		viewData["country"] = controller.SessionManager.GetString(r.Context(), "country")

		viewData["policies"] = policies
		viewData["policiesApplied"] = "false"

		if len(policies) != 0 {
			viewData["policiesApplied"] = "true"
		}

		ts, err := template.ParseFiles("./static/views/confirmtransaction.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error 1", 500)
			return
		}

		err = ts.Execute(w, viewData)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error 2", 500)
		}

	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error parsing form", 500)
		}

		complianceCheckId := r.Form.Get("transactionid")
		check := controller.DB.GetComplianceCheckByID(complianceCheckId)

		controller.DB.UpdateTransactionState(check.Id, 2)

		controller.RulesEngine.Do(complianceCheckId, "interactive", map[string]any{"vm_address": ""})

		// TODO: This should probably be handled inside of the rules engine
		//		 CB should be notified about the compliance proof request so it knows about the checks
		//		 and their states
		controller.DB.UpdateTransactionState(check.Id, 3)

		checkConfirmedData := common.CheckConfirmedDTO{
			CheckID:   complianceCheckId,
			VMAddress: controller.RulesEngine.GetVMAddress(),
		}

		_, err = controller.P2PClient.Send(check.OriginatorBankId, "check-confirmed", checkConfirmedData, 0)
		if err != nil {
			http.Error(w, fmt.Sprint("Internal Server Error %w", err), 500)
		}

		sender := controller.DB.GetClientByID(check.SenderId)
		receiver := controller.DB.GetClientByID(check.ReceiverId)

		transactionDto := common.TransactionDTO{
			TransactionID:                   check.Id,
			SenderLei:                       sender.GlobalIdentifier,
			SenderName:                      sender.Name,
			ReceiverLei:                     receiver.GlobalIdentifier,
			ReceiverName:                    receiver.Name,
			OriginatorBankGlobalIdentifier:  check.OriginatorBankId,
			BeneficiaryBankGlobalIdentifier: check.BeneficiaryBankId,
			PaymentType:                     "",
			TransactionType:                 fmt.Sprint(check.TransactionTypeId),
			Amount:                          uint64(check.Amount),
			Currency:                        check.Currency,
			SwiftBICCode:                    "",
			LoanID:                          uint64(0),
		}

		_, err = controller.P2PClient.Send(config.ResolveCBGlobalIdentifier(), "create-transaction", transactionDto, 0)
		if err != nil {
			http.Error(w, fmt.Sprint("Internal Server Error %w", err), 500)
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (controller *TransactionController) TransactionHistory(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error parsing form", 500)
	}

	transaction := controller.DB.GetTransactionHistory(r.Form.Get("transaction"))

	bankId := controller.DB.GetBankId(transaction.BeneficiaryBank)

	policies := controller.DB.GetPolices(bankId, transaction.TypeId)

	policiesAndStatuses := []struct {
		Policy models.PolicyModel
		Status int
	}{}

	// for _, onePolicy := range policies {
	// 	currentStatus := controller.DB.GetTransactionPolicyStatus(uint64(transactionId), int(onePolicy.Id))
	// 	policiesAndStatuses = append(policiesAndStatuses, struct {
	// 		Policy models.PolicyModel
	// 		Status int
	// 	}{onePolicy, currentStatus})
	// }

	viewData := map[string]any{}

	viewData["username"] = controller.SessionManager.GetString(r.Context(), "username")
	viewData["transaction"] = transaction
	viewData["bankName"] = controller.SessionManager.GetString(r.Context(), "bankName")
	viewData["country"] = controller.SessionManager.GetString(r.Context(), "country")
	viewData["centralBankEmployee"] = controller.SessionManager.GetBool(r.Context(), "centralBankEmployee")

	viewData["policies"] = policiesAndStatuses
	viewData["policiesApplied"] = "false"

	if len(policies) != 0 {
		viewData["policiesApplied"] = "true"
	}

	ts, err := template.ParseFiles("./static/views/transactionhistory.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	err = ts.Execute(w, viewData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}
