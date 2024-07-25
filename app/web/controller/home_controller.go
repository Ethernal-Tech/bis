package controller

import (
	"bisgo/config"
	"bisgo/errlog"
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

func (c *HomeController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if c.SessionManager.GetString(r.Context(), "inside") == "yes" {
			http.Redirect(w, r, "/home", http.StatusSeeOther)

			return
		}

		ts, err := template.ParseFiles("./app/web/static/views/login.html")
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = ts.Execute(w, struct{}{})
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		user, err := c.DB.GetBankEmployee(r.Form.Get("username"), r.Form.Get("password"))
		if err != nil {
			if err == errlog.ErrBankEmployee404 {
				http.Error(w, "Bank employee can't be found", 404)
			}

			errlog.Println(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		if config.ResolveIsCentralBank() {
			c.SessionManager.Put(r.Context(), "centralBankEmployee", "yes")
		} else {
			c.SessionManager.Put(r.Context(), "centralBankEmployee", "no")
		}

		c.SessionManager.Put(r.Context(), "inside", "yes")
		c.SessionManager.Put(r.Context(), "username", user.Username)
		c.SessionManager.Put(r.Context(), "bankId", user.BankId)

		jurisdiction, err := c.DB.GetBankJurisdiction(user.BankId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}
		c.SessionManager.Put(r.Context(), "jurisdiction", jurisdiction.Name)

		bank, err := c.DB.GetBankByGlobalIdentifier(user.BankId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		c.SessionManager.Put(r.Context(), "bankName", bank.Name)

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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

	viewData["username"] = controller.SessionManager.GetString(r.Context(), "username")
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
