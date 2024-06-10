package controller

import (
	"log"
	"math"
	"net/http"
	"text/template"
)

func (controller *CBController) ShowAnalytics(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" || !controller.SessionManager.GetBool(r.Context(), "centralBankEmployee") {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	viewData := map[string]any{}

	viewData["username"] = controller.SessionManager.GetString(r.Context(), "username")
	viewData["bankName"] = controller.SessionManager.GetString(r.Context(), "bankName")
	viewData["country"] = controller.SessionManager.GetString(r.Context(), "country")
	viewData["centralBankEmployee"] = controller.SessionManager.GetBool(r.Context(), "centralBankEmployee")

	centralBankId := controller.SessionManager.Get(r.Context(), "bankId").(string)
	transactions, countryId := controller.DB.GetCentralBankTransactions(centralBankId, "")

	sentAmount := 0
	receivedAmount := 0
	successfulTxs := 0
	failedBecauseOfSanctionsCheck := 0
	failedBecauseOfCFMCheck := 0

	for _, tx := range transactions {
		if tx.Status == "COMPLETED" {
			successfulTxs += 1

			if tx.OriginatorBankCountryId == countryId {
				sentAmount += tx.Amount
			} else {
				receivedAmount += tx.Amount
			}
		} else if tx.Status == "CANCELED" {
			policyStatuses := controller.DB.GetTransactionPolicyStatuses(tx.Id)

			for _, policyStatus := range policyStatuses {
				policy := controller.DB.GetPolicyById(policyStatus.PolicyId)
				if policy.Code == "CFM" && policyStatus.Status == 2 {
					failedBecauseOfCFMCheck += 1
				}
				if policy.Code == "SCL" && policyStatus.Status == 2 {
					failedBecauseOfSanctionsCheck += 1
				}
			}
		}
	}

	viewData["sent"] = sentAmount
	viewData["received"] = receivedAmount
	viewData["successful"] = successfulTxs
	viewData["initialized"] = len(transactions)

	if successfulTxs == 0 {
		viewData["percentage"] = 0
	} else {
		percentage := float64(successfulTxs) / float64(len(transactions)) * 100
		viewData["percentage"] = math.Floor(percentage*100) / 100
	}

	viewData["sclFails"] = failedBecauseOfSanctionsCheck
	viewData["cfmFails"] = failedBecauseOfCFMCheck

	ts, err := template.ParseFiles("./static/views/analytics.html")
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
