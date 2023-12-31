package main

import (
	"bisgo/DB"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") == "yes" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)

		return
	}

	ts, err := template.ParseFiles("./static/views/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, struct{}{})
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
	app.sessionManager.Put(r.Context(), "country", app.db.GetCountry(uint(app.db.GetBank(user.BankId).CountryId)).Name)

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
	viewData["country"] = app.sessionManager.GetString(r.Context(), "country")

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

		loanID := rand.Intn(2500000)

		viewData["loanId"] = loanID
		viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
		viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")
		viewData["country"] = app.sessionManager.GetString(r.Context(), "country")
		viewData["banks"] = app.db.GetBanks()
		viewData["transactionTypes"] = app.db.GetTransactionTypes()

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
		amount, _ := strconv.Atoi(strings.Replace(r.Form.Get("amount"), ",", "", -1))
		transactionType, _ := strconv.Atoi(r.Form.Get("type"))
		loanId, _ := strconv.Atoi(strings.Replace(r.Form.Get("loanId"), ",", "", -1))

		transaction := DB.Transaction{
			OriginatorBank:  uint64(originatorBank),
			BeneficiaryBank: uint64(beneficiaryBank),
			Sender:          sender,
			Receiver:        receiver,
			Currency:        currency,
			Amount:          amount,
			TypeId:          transactionType,
			LoanId:          loanId,
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
	transactionTypeId, _ := strconv.Atoi(data.TransactionTypeId)

	policies := app.db.GetPolices(uint64(bankId), transactionTypeId)

	jsonData, err := json.Marshal(policies)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (app *application) showPolicies(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	viewData := map[string]any{}

	viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
	viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")
	viewData["country"] = app.sessionManager.GetString(r.Context(), "country")

	policies := app.db.PoliciesFromCountry(app.sessionManager.Get(r.Context(), "bankId").(uint64))

	viewData["policies"] = policies

	ts, err := template.ParseFiles("./static/views/policies.html")
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
		viewData["country"] = app.sessionManager.GetString(r.Context(), "country")

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

		// CFM check //

		bank := app.db.GetBank(app.db.GetBankId(transaction.BeneficiaryBank))

		amount := app.db.CheckCFM(app.db.GetBankClientId(transaction.ReceiverName), bank.CountryId)

		policies := app.db.GetPolices(app.db.GetBankId(transaction.BeneficiaryBank), transaction.TypeId)

		var CFMpolicy DB.PolicyModel
		CFMpolicy.Id = 0
		CFMexists := false
		SCLexists := false
		var SCLpolicyId int
		var country string

		for _, policy := range policies {
			country = policy.Country

			if policy.Code == "CFM" {
				CFMpolicy = policy
				CFMexists = true
			} else if policy.Code == "SCL" {
				SCLpolicyId = app.db.GetPolicyId(policy.Code, policy.CountryId)
				SCLexists = true
			}
		}

		policyValid := false

		if CFMpolicy.Id != 0 {
			var ratio = 3.4
			var newAmount = float64(amount+int64(transaction.Amount)) * ratio
			if newAmount >= float64(CFMpolicy.Amount) {
				app.db.UpdateTransactionPolicyStatus(transaction.Id, int(CFMpolicy.Id), 2)
			} else {
				app.db.UpdateTransactionPolicyStatus(transaction.Id, int(CFMpolicy.Id), 1)

				policyValid = true
			}
		}

		if !CFMexists && !SCLexists {
			app.db.UpdateTransactionState(transaction.Id, 6)
			app.db.UpdateTransactionState(transaction.Id, 7)

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		if !SCLexists {
			if policyValid {
				app.db.UpdateTransactionState(transaction.Id, 6)
				app.db.UpdateTransactionState(transaction.Id, 7)
			} else {
				app.db.UpdateTransactionState(transaction.Id, 8)
			}

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		// SCL //

		app.db.UpdateTransactionState(transaction.Id, 3)

		var urlServer string
		var jsonPayloadServer []byte
		var urlClient string
		var jsonPayloadClient []byte

		if country == "Malaysia" {
			urlServer = "http://" + os.Getenv("API_MY") + ":9090/api/start-server"
			jsonPayloadServer = []byte(fmt.Sprintf(`{"tx_id": "%d", "policy_id": "%d"}`, transactionId, SCLpolicyId))

			urlClient = "http://" + os.Getenv("API_SG") + ":9090/api/start-client"
			jsonPayloadClient = []byte(fmt.Sprintf(`{"tx_id": "%d", "receiver": "%s", "to": "%s:10501"}`, transactionId, transaction.ReceiverName, os.Getenv("GPJC_MY")))

		} else if country == "Singapore" {
			urlServer = "http://" + os.Getenv("API_SG") + ":9090/api/start-server"
			jsonPayloadServer = []byte(fmt.Sprintf(`{"tx_id": "%d", "policy_id": "%d"}`, transactionId, SCLpolicyId))

			urlClient = "http://" + os.Getenv("API_MY") + ":9090/api/start-client"
			jsonPayloadClient = []byte(fmt.Sprintf(`{"tx_id": "%d", "receiver": "%s", "to": "%s:10501"}`, transactionId, transaction.ReceiverName, os.Getenv("GPJC_SG")))

		} else {
			log.Println("Error in SCL")
			http.Error(w, "Internal Server Error", 500)

			return
		}

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

	policiesAndStatuses := []struct {
		Policy DB.PolicyModel
		Status int
	}{}

	for _, onePolicy := range policies {
		currentStatus := app.db.GetTransactionPolicyStatus(uint64(transactionId), int(onePolicy.Id))
		policiesAndStatuses = append(policiesAndStatuses, struct {
			Policy DB.PolicyModel
			Status int
		}{onePolicy, currentStatus})
	}

	viewData := map[string]any{}

	viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
	viewData["transaction"] = transaction
	viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")
	viewData["country"] = app.sessionManager.GetString(r.Context(), "country")

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
		policyId, _ := strconv.Atoi(messageData.PolicyId)
		app.db.UpdateTransactionPolicyStatus(uint64(transactionId), policyId, 1)
		app.db.UpdateTransactionState(uint64(transactionId), 4)
	} else {
		policyId, _ := strconv.Atoi(messageData.PolicyId)
		app.db.UpdateTransactionPolicyStatus(uint64(transactionId), policyId, 2)
		app.db.UpdateTransactionState(uint64(transactionId), 5)
	}

	policyStatuses := app.db.GetTransactionPolicyStatuses(uint64(transactionId))

	check := true

	for _, status := range policyStatuses {
		if status.Status != 1 {
			check = false
		}
	}

	if check {
		app.db.UpdateTransactionState(uint64(transactionId), 6)
		app.db.UpdateTransactionState(uint64(transactionId), 7)
	} else {
		app.db.UpdateTransactionState(uint64(transactionId), 8)
	}

	json.NewEncoder(w).Encode("Ok")
}
