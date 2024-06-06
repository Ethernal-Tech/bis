package controller

import (
	"bisgo/app/models"
	"bisgo/common"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (controller *TransactionController) SearchTransaction(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	searchValue := r.URL.Query().Get("searchValue")

	viewData := map[string]any{}
	var transactions []models.TransactionModel
	if controller.SessionManager.GetBool(r.Context(), "centralBankEmployee") {
		var countryId int
		transactions, countryId = controller.DB.GetTransactionsForCentralbank(controller.SessionManager.Get(r.Context(), "bankId").(string), searchValue)
		viewData["countryId"] = countryId
	} else {
		transactions = controller.DB.GetTransactionsForAddress(controller.SessionManager.Get(r.Context(), "bankId").(string), searchValue)
	}

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

		fmt.Println(data)

		originatorBankId := controller.DB.GetBankId(controller.SessionManager.GetString(r.Context(), "bankName"))

		amount, _ := strconv.Atoi(strings.Replace(data.Amount, ",", "", -1))
		transactionType, _ := strconv.Atoi(data.TransactionTypeID)
		//loanId, _ := strconv.Atoi(data.PaymentTypeID)

		// TODO: If client is not in the DB add it

		transaction := models.NewTransaction{
			OriginatorBankId:  originatorBankId,
			BeneficiaryBankId: data.BeneficiaryBank,
			SenderId:          data.SenderLei,
			ReceiverId:        data.BeneficiaryLei,
			Currency:          data.Currency,
			Amount:            amount,
			TransactionTypeId: transactionType,
			LoanId:            0,
		}

		transactionID := controller.DB.InsertTransaction(transaction)
		controller.DB.UpdateTransactionState(transactionID, 1)
		fmt.Println(transactionID)

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

		transactionId, _ := strconv.Atoi(r.Form.Get("transaction"))

		transaction := controller.DB.GetTransactionHistory(uint64(transactionId))

		bankId := controller.DB.GetBankId(transaction.BeneficiaryBank)

		policies := controller.DB.GetPolices(bankId, transaction.TypeId)

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
		applicablePolicies := controller.DB.GetPoliciesForTransaction(complianceCheckId)
		check := controller.DB.GetComplianceCheckByID(complianceCheckId)

		controller.DB.UpdateTransactionState(check.Id, 2)

		for _, policy := range applicablePolicies {
			if policy.PolicyType.Code == "CFM" {
				// CFM check //
				// TODO: Notify central bank about the CFM check if it exists
			} else if policy.PolicyType.Code == "SCL" {
				// SCL //
				// TODO: Start gpjc proving server
				err = controller.ProvingClient.SendProofRequest("interactive", complianceCheckId, policy.Policy.Id, "", "")
				if err != nil {
					http.Error(w, fmt.Sprint("Internal Server Error %w", err), 500)
				}
			}
		}

		controller.DB.UpdateTransactionState(check.Id, 3)

		checkConfirmedData := common.CheckConfirmedDTO{
			CheckID:   complianceCheckId,
			VMAddress: controller.ProvingClient.GetVMAddress(),
		}

		_, err = controller.P2PClient.Send(check.OriginatorBankId, "check-confirmed", checkConfirmedData, 0)
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

	transactionId, _ := strconv.Atoi(r.Form.Get("transaction"))

	transaction := controller.DB.GetTransactionHistory(uint64(transactionId))

	bankId := controller.DB.GetBankId(transaction.BeneficiaryBank)

	policies := controller.DB.GetPolices(bankId, transaction.TypeId)

	policiesAndStatuses := []struct {
		Policy models.PolicyModel
		Status int
	}{}

	for _, onePolicy := range policies {
		currentStatus := controller.DB.GetTransactionPolicyStatus(uint64(transactionId), int(onePolicy.Id))
		policiesAndStatuses = append(policiesAndStatuses, struct {
			Policy models.PolicyModel
			Status int
		}{onePolicy, currentStatus})
	}

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
