package controller

import (
	"bisgo/app/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"time"
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
			SenderLei       string `json:"senderLei"`
			SenderName      string `json:"senderName"`
			BeneficiaryLei  string `json:"beneficiaryLei"`
			BeneficiaryName string `json:"beneficiaryName"`
		}{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&data); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		fmt.Println(data)
		// err := r.ParseForm()
		// if err != nil {
		// 	log.Println(err.Error())
		// 	http.Error(w, "Internal Server Error parsing form", 500)
		// }

		// originatorBank := controller.DB.GetBankId(controller.SessionManager.GetString(r.Context(), "bankName"))
		// //beneficiaryBank, _ := strconv.Atoi(r.Form.Get("bank"))
		// sender := controller.DB.GetBankClientId(r.Form.Get("sender"))
		// receiver := controller.DB.GetBankClientId(r.Form.Get("receiver"))
		// currency := r.Form.Get("currency")
		// amount, _ := strconv.Atoi(strings.Replace(r.Form.Get("amount"), ",", "", -1))
		// transactionType, _ := strconv.Atoi(r.Form.Get("type"))
		// loanId, _ := strconv.Atoi(strings.Replace(r.Form.Get("loanId"), ",", "", -1))

		// /*
		// 	OriginatorBank:  uint64(originatorBank),
		// 		BeneficiaryBank: uint64(beneficiaryBank),
		// 		Sender:          sender,
		// 		Receiver:        receiver,
		// 		Currency:        currency,
		// 		Amount:          amount,
		// 		TypeId:          transactionType,
		// 		LoanId:          loanId, */

		// // TODO: Generate tx id
		// // TODO: Get beneficiary lei
		// transaction := models.NewTransaction{
		// 	Id:                "",
		// 	OriginatorBankId:  originatorBank,
		// 	BeneficiaryBankId: "",
		// 	SenderId:          sender,
		// 	ReceiverId:        receiver,
		// 	Currency:          currency,
		// 	Amount:            amount,
		// 	TransactionTypeId: transactionType,
		// 	LoanId:            loanId,
		// }

		// // Call P2P create-transaction of beneficiary bank
		// // TODO: Get beneficiary lei
		// transactionDto := common.TransactionDTO{
		// 	TransactionID:                   "0",
		// 	SenderLei:                       "",
		// 	SenderName:                      r.Form.Get("sender"),
		// 	ReceiverLei:                     "",
		// 	ReceiverName:                    r.Form.Get("receiver"),
		// 	OriginatorBankGlobalIdentifier:  originatorBank,
		// 	BeneficiaryBankGlobalIdentifier: "",
		// 	PaymentType:                     "",
		// 	TransactionType:                 fmt.Sprint(transactionType),
		// 	Amount:                          uint64(amount),
		// 	Currency:                        currency,
		// 	SwiftBICCode:                    "",
		// 	LoanID:                          uint64(loanId),
		// }

		// ch, err := controller.P2PClient.Send(transactionDto.BeneficiaryBankGlobalIdentifier, "create-transaction", transactionDto, 0)

		// _ = ch

		// if err != nil {
		// 	log.Println(err.Error())
		// 	http.Error(w, "Internal Server Error sending create tx", 500)
		// }

		// transactionID := controller.DB.InsertTransaction(transaction)
		// controller.DB.UpdateTransactionState(transactionID, 1)

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

		transactionId, _ := strconv.Atoi(r.Form.Get("transactionid"))

		transaction := controller.DB.GetTransactionHistory(uint64(transactionId))

		controller.DB.UpdateTransactionState(transaction.Id, 2)

		// CFM check //

		bank := controller.DB.GetBank(controller.DB.GetBankId(transaction.BeneficiaryBank))

		amount := controller.DB.CheckCFM(controller.DB.GetBankClientId(transaction.ReceiverName), int(bank.CountryId))

		policies := controller.DB.GetPolices(controller.DB.GetBankId(transaction.BeneficiaryBank), transaction.TypeId)

		var CFMpolicy models.PolicyModel
		CFMpolicy.Id = 0
		CFMexists := false
		SCLexists := false
		var SCLpolicyId int

		for _, policy := range policies {
			if policy.Code == "CFM" {
				CFMpolicy = policy
				CFMexists = true
			} else if policy.Code == "SCL" {
				SCLpolicyId = controller.DB.GetPolicyId(policy.Code, policy.CountryId)
				SCLexists = true
			}
		}

		policyValid := false

		if CFMpolicy.Id != 0 {
			var ratio = 3.4
			var newAmount = float64(amount+int64(transaction.Amount)) * ratio
			if newAmount >= float64(CFMpolicy.Amount) {
				controller.DB.UpdateTransactionPolicyStatus(transaction.Id, int(CFMpolicy.Id), 2)
			} else {
				controller.DB.UpdateTransactionPolicyStatus(transaction.Id, int(CFMpolicy.Id), 1)

				policyValid = true
			}
		}

		if !CFMexists && !SCLexists {
			controller.DB.UpdateTransactionState(transaction.Id, 6)
			controller.DB.UpdateTransactionState(transaction.Id, 7)

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		if !SCLexists {
			if policyValid {
				controller.DB.UpdateTransactionState(transaction.Id, 6)
				controller.DB.UpdateTransactionState(transaction.Id, 7)
			} else {
				controller.DB.UpdateTransactionState(transaction.Id, 8)
			}

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		// SCL //

		controller.DB.UpdateTransactionState(transaction.Id, 3)

		var urlServer string
		var jsonPayloadServer []byte
		var urlClient string
		var jsonPayloadClient []byte

		// TODO: This will be updated in upcoming messaging system changes

		urlServer = "http://" + "controller.Config.GpjcApiAddress" + ":9090/api/start-server"
		jsonPayloadServer = []byte(fmt.Sprintf(`{"tx_id": "%d", "policy_id": "%d"}`, transactionId, SCLpolicyId))

		urlClient = "http://" + "gpjc_client" + ":9090/api/start-client"
		jsonPayloadClient = []byte(fmt.Sprintf(`{"tx_id": "%d", "policy_id": "%d", "receiver": "%s", "to": "%s:10501"}`, transactionId, SCLpolicyId, transaction.ReceiverName, "controller.Config.GpjcApiAddress"))

		client := &http.Client{}

		req, err := http.NewRequest("POST", urlServer, bytes.NewBuffer(jsonPayloadServer))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Connection", "close")

		_, err = client.Do(req)
		if err != nil {
			panic(err)
		}

		time.Sleep(100 * time.Millisecond)

		req, err = http.NewRequest("POST", urlClient, bytes.NewBuffer(jsonPayloadClient))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Connection", "close")

		_, err = client.Do(req)
		if err != nil {
			panic(err)
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
