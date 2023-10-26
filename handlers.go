package main

import (
	"bisgo/DB"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, app.sessionManager.GetString(r.Context(), "inside"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := app.db.Login(r.Form.Get("username"), r.Form.Get("password"))

	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	app.sessionManager.Put(r.Context(), "inside", "yes")
	app.sessionManager.Put(r.Context(), "username", user.Name)
	app.sessionManager.Put(r.Context(), "bankId", user.BankId)
	app.sessionManager.Put(r.Context(), "bankName", user.BankName)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Put(r.Context(), "inside", "no")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	transactions := app.db.GetTransactionsForAddress(app.sessionManager.Get(r.Context(), "bankId").(uint64))

	viewData := map[string]any{}

	viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
	viewData["transactions"] = transactions
	viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")

	ts, err := template.ParseFiles("./static/views/home.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, viewData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) addTransaction(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodGet {

		viewData := map[string]any{}

		viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
		viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")
		viewData["banks"] = app.db.GetBanks()

		ts, err := template.ParseFiles("./static/views/addtransaction.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error 1", 500)
			return
		}

		ts.Execute(w, viewData)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error 2", 500)
		}

	} else if r.Method == http.MethodPost {

		r.ParseForm()

		originatorBank := app.db.GetBankId(app.sessionManager.GetString(r.Context(), "bankName"))
		beneficiaryBank, _ := strconv.Atoi(r.Form.Get("bank"))
		sender := app.db.GetBankClientId(r.Form.Get("sender"))
		receiver := app.db.GetBankClientId(r.Form.Get("receiver"))
		currency := r.Form.Get("currency")
		amount, _ := strconv.Atoi(r.Form.Get("amount"))
		transactionType := app.db.GetTransactionTypeId(r.Form.Get("type"))

		transaction := DB.Transaction{
			OriginatorBank:  uint64(originatorBank),
			BeneficiaryBank: uint64(beneficiaryBank),
			Sender:          sender,
			Receiver:        receiver,
			Currency:        currency,
			Amount:          amount,
			TypeId:          transactionType,
		}

		transactionID := app.db.InsertTransaction(transaction)
		app.db.UpdateTransactionState(transactionID, 1)

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (app *application) getPolicies(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
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

	policies := app.db.GetPolices(uint64(bankId), app.db.GetTransactionTypeId(data.TransactionTypeId))

	jsonData, err := json.Marshal(policies)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (app *application) confirmTransaction(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodGet {

		r.ParseForm()

		transactionId, _ := strconv.Atoi(r.Form.Get("transaction"))

		transaction := app.db.GetTransactionHistory(uint64(transactionId))

		bankId := app.db.GetBankId(transaction.BeneficiaryBank)

		policies := app.db.GetPolices(bankId, transaction.TypeId)

		viewData := map[string]any{}

		viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
		viewData["transaction"] = transaction
		viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")

		viewData["CapitalFlowManagement"] = "false"
		viewData["SactionCheckList"] = "false"

		for _, policy := range policies {
			viewData[strings.ReplaceAll(policy.Name, " ", "")] = "true"
			viewData[strings.ReplaceAll(policy.Name, " ", "")+"Content"] = policy
		}

		ts, err := template.ParseFiles("./static/views/confirmtransaction.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error 1", 500)
			return
		}

		ts.Execute(w, viewData)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error 2", 500)
		}

	} else if r.Method == http.MethodPost {

		r.ParseForm()

		transactionId, _ := strconv.Atoi(r.Form.Get("transactionid"))

		transaction := app.db.GetTransactionHistory(uint64(transactionId))

		app.db.UpdateTransactionState(transaction.Id, 2)

		// CFM check//

		bank := app.db.GetBank(app.db.GetBankId(transaction.BeneficiaryBank))

		fmt.Println(app.db.GetBankClientId(transaction.ReceiverName))
		fmt.Println(bank.CountryId)

		amount := app.db.CheckCFM(app.db.GetBankClientId(transaction.ReceiverName), bank.CountryId)

		fmt.Println(amount)

		policies := app.db.GetPolices(app.db.GetBankId(transaction.BeneficiaryBank), transaction.TypeId)

		var CFMpolicy *DB.PolicyModel
		SCLexists := false

		for _, policy := range policies {
			if policy.Name == "Capital Flow Management" {
				CFMpolicy = &policy
			} else if policy.Name == "Saction Check List" {
				SCLexists = true
			}
		}

		policyValid := false

		if CFMpolicy != nil {
			if amount+int64(transaction.Amount) >= int64(CFMpolicy.Amount) {
				app.db.UpdateTransactionPolicyStatus(transaction.Id, int(CFMpolicy.Id), 2)
			} else {
				app.db.UpdateTransactionPolicyStatus(transaction.Id, int(CFMpolicy.Id), 1)

				policyValid = true
			}
		}

		if !SCLexists {
			if policyValid {
				app.db.UpdateTransactionState(transaction.Id, 6)
				app.db.UpdateTransactionState(transaction.Id, 7)
			} else {

				app.db.UpdateTransactionState(transaction.Id, 8)
			}
			return
		}

		return

		// SCL //

		app.db.UpdateTransactionState(transaction.Id, 3)

		urlServer := "http://localhost:9090/api/start-server"
		jsonPayloadServer := []byte(fmt.Sprintf(`{"tx_id": "%d", "policy_id": "1"}`, transactionId))

		urlClient := "http://localhost:9090/api/start-client"
		jsonPayloadClient := []byte(fmt.Sprintf(`{"tx_id": "%d", "receiver": "%s", "to": "0.0.0.0:10501"}`, transactionId, transaction.ReceiverName))

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
	}
}

func (app *application) transactionHistory(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	r.ParseForm()

	transactionId, _ := strconv.Atoi(r.Form.Get("transaction"))

	transaction := app.db.GetTransactionHistory(uint64(transactionId))

	bankId := app.db.GetBankId(transaction.BeneficiaryBank)

	policies := app.db.GetPolices(bankId, transaction.TypeId)

	viewData := map[string]any{}

	viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
	viewData["transaction"] = transaction
	viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")

	viewData["Capital Flow Management"] = "false"
	viewData["Saction Check List"] = "false"

	for _, policy := range policies {
		viewData[strings.ReplaceAll(policy.Name, " ", "")] = "true"
		viewData[strings.ReplaceAll(policy.Name, " ", "")+"Content"] = policy
	}

	ts, err := template.ParseFiles("./static/views/transactionhistory.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, viewData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) submitTransactionProof(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var messageData DB.TransactionProofRequest
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
	app.db.InsertTransactionProof(uint64(transactionId), messageData.Value)

	if messageData.Value == "0" {
		app.db.UpdateTransactionPolicyStatus(uint64(transactionId), 2, 1)
	} else {
		app.db.UpdateTransactionPolicyStatus(uint64(transactionId), 2, 2)
	}

	policyStatuses := app.db.GetTransactionPolicyStatuses(uint64(transactionId))

	check := true

	for _, status := range policyStatuses {
		if status != 1 {
			check = false
		}
	}

	if check {
		app.db.UpdateTransactionState(uint64(transactionId), 4)
		app.db.UpdateTransactionState(uint64(transactionId), 6)
		app.db.UpdateTransactionState(uint64(transactionId), 7)
	} else {
		app.db.UpdateTransactionState(uint64(transactionId), 5)
		app.db.UpdateTransactionState(uint64(transactionId), 8)
	}

	json.NewEncoder(w).Encode("Ok")
}
