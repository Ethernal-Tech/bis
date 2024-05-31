package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (controller *PolicyController) ShowPolicies(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	viewData := map[string]any{}

	viewData["username"] = controller.SessionManager.GetString(r.Context(), "username")
	viewData["bankName"] = controller.SessionManager.GetString(r.Context(), "bankName")
	viewData["country"] = controller.SessionManager.GetString(r.Context(), "country")
	viewData["centralBankEmployee"] = controller.SessionManager.GetBool(r.Context(), "centralBankEmployee")

	policies := controller.DB.PoliciesFromCountry(controller.SessionManager.Get(r.Context(), "bankId").(uint64))

	viewData["policies"] = policies

	ts, err := template.ParseFiles("./static/views/policies.html")
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

func (controller *PolicyController) EditPolicy(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" || !controller.SessionManager.GetBool(r.Context(), "centralBankEmployee") {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error parsing form", 500)
			return
		}

		policyId, _ := strconv.Atoi(r.Form.Get("policyId"))

		viewData := map[string]any{}

		viewData["username"] = controller.SessionManager.GetString(r.Context(), "username")
		viewData["bankName"] = controller.SessionManager.GetString(r.Context(), "bankName")
		viewData["country"] = controller.SessionManager.GetString(r.Context(), "country")
		viewData["policy"] = controller.DB.GetPolicy(viewData["country"].(string), uint64(policyId))
		viewData["centralBankEmployee"] = controller.SessionManager.GetBool(r.Context(), "centralBankEmployee")

		ts, err := template.ParseFiles("./static/views/editpolicy.html")
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
	} else {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error parsing form", 500)
			return
		}

		policyId, _ := strconv.Atoi(r.FormValue("policyId"))
		originalPolicy := controller.DB.GetPolicy(controller.SessionManager.GetString(r.Context(), "country"), uint64(policyId))

		if originalPolicy.Code == "CFM" {
			amount, err := strconv.Atoi(strings.Replace(r.Form.Get("amount"), ",", "", -1))
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Server Error parsing form", 500)
				return
			}

			controller.DB.UpdatePolicyAmount(uint64(policyId), uint64(amount))
		} else {
			fileName, err := controller.SanctionListManager.GetNewestSanctionsList()
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Server Error retrieving file", 500)
				return
			}

			controller.DB.UpdatePolicyChecklist(uint64(policyId), fileName)
		}

		http.Redirect(w, r, "/policies", http.StatusSeeOther)
	}
}
