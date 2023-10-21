package main

import (
	"bisgo/DB"
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

	viewData := map[string]any{}

	viewData["username"] = app.sessionManager.GetString(r.Context(), "username")

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

	fmt.Println(policies)

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

	fmt.Println(policies)

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

	json.NewEncoder(w).Encode("Ok")
}
