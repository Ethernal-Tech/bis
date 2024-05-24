package controller

import (
	"bisgo/core/DB/models"
	"log"
	"net/http"
	"text/template"
)

func (controller *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") == "yes" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)

		return
	}

	ts, err := template.ParseFiles("./static/views/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	err = ts.Execute(w, struct{}{})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (controller *HomeController) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error parsing form", 500)
	}

	user := controller.DB.Login(r.Form.Get("username"), r.Form.Get("password"))

	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	centralBankEmploye := controller.DB.IsCentralBankEmployee(user.Username)

	controller.SessionManager.Put(r.Context(), "inside", "yes")
	controller.SessionManager.Put(r.Context(), "username", user.Name)
	controller.SessionManager.Put(r.Context(), "bankId", user.BankId)
	controller.SessionManager.Put(r.Context(), "bankName", user.BankName)
	controller.SessionManager.Put(r.Context(), "country", controller.DB.GetCountry(uint(controller.DB.GetBank(user.BankId).CountryId)).Name)
	controller.SessionManager.Put(r.Context(), "centralBankEmployee", centralBankEmploye)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (controller *HomeController) Logout(w http.ResponseWriter, r *http.Request) {
	controller.SessionManager.Put(r.Context(), "inside", "no")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (controller *HomeController) Home(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}
	viewData := map[string]any{}

	var transactions []models.TransactionModel
	if controller.SessionManager.GetBool(r.Context(), "centralBankEmployee") {
		var countryId int
		transactions, countryId = controller.DB.GetTransactionsForCentralbank(controller.SessionManager.Get(r.Context(), "bankId").(uint64), "")
		viewData["countryId"] = countryId
	} else {
		transactions = controller.DB.GetTransactionsForAddress(controller.SessionManager.Get(r.Context(), "bankId").(uint64), "")
	}

	viewData["username"] = controller.SessionManager.GetString(r.Context(), "username")
	viewData["transactions"] = transactions
	viewData["bankName"] = controller.SessionManager.GetString(r.Context(), "bankName")
	viewData["country"] = controller.SessionManager.GetString(r.Context(), "country")
	viewData["centralBankEmployee"] = controller.SessionManager.GetBool(r.Context(), "centralBankEmployee")

	ts, err := template.ParseFiles("./static/views/home.html")
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
