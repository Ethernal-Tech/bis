package controller

import (
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// AddComplianceCheck handles a web GET/POST "/addcompliancecheck" request. For a GET, it responds with a view for a new
// compliance check creation. On the other hand, POST indicates confirmation and that a compliance check should be created.
// As part of this process, a new compliance check is also sent over the p2p network to the beneficiary bank.
func (c *ComplianceCheckController) AddComplianceCheck(w http.ResponseWriter, r *http.Request) {
	if c.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodGet {

		// TODO: add logic together with a new addcompliancecheck view
		// TODO: modify json names for marshal/unmarshal

	} else if r.Method == http.MethodPost {

		data := struct {
			OriginatorGlobalIdentifier      string `json:"senderLei"`
			OriginatorName                  string `json:"senderName"`
			BeneficiaryGlobalIdentifier     string `json:"beneficiaryLei"`
			BeneficiaryName                 string `json:"beneficiaryName"`
			PaymentTypeId                   string `json:"paymentType"`
			TransactionTypeId               string `json:"transactionType"`
			Currency                        string `json:"currency"`
			Amount                          string `json:"amount"`
			BeneficiaryBankGlobalIdentifier string `json:"beneficiaryBank"`
		}{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		originatorBankGlobalIdentifier, err := c.DB.GetBankIdByName(c.SessionManager.GetString(r.Context(), "bankName"))
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		amount, err := strconv.Atoi(strings.Replace(data.Amount, ",", "", -1))
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		transactionTypeId, err := strconv.Atoi(data.TransactionTypeId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		// TODO: potentially handle loan

		originatorId, err := c.DB.CreateOrGetBankClient(data.OriginatorGlobalIdentifier, data.OriginatorName, "", originatorBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		beneficiaryId, err := c.DB.CreateOrGetBankClient(data.BeneficiaryGlobalIdentifier, data.BeneficiaryName, "", data.BeneficiaryBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		complianceCheck := models.ComplianceCheck{
			OriginatorBankId:  originatorBankGlobalIdentifier,
			BeneficiaryBankId: data.BeneficiaryBankGlobalIdentifier,
			SenderId:          originatorId,
			ReceiverId:        beneficiaryId,
			Currency:          data.Currency,
			Amount:            amount,
			TransactionTypeId: transactionTypeId,
			LoanId:            0,
		}

		complianceCheckId, err := c.DB.AddComplianceCheck(complianceCheck)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		// TODO: a compliance check status manager call instead of a direct state change
		err = c.DB.UpdateComplianceCheckStatus(complianceCheckId, 1)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		transactionType, err := c.DB.GetTransactionTypeById(transactionTypeId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		complianceCheckDTO := common.ComplianceCheckDTO{
			ComplianceCheckId:               complianceCheckId,
			OriginatorGlobalIdentifier:      data.OriginatorGlobalIdentifier,
			OriginatorName:                  data.OriginatorName,
			BeneficiaryGlobalIdentifier:     data.BeneficiaryBankGlobalIdentifier,
			BeneficiaryName:                 data.BeneficiaryName,
			OriginatorBankGlobalIdentifier:  originatorBankGlobalIdentifier,
			BeneficiaryBankGlobalIdentifier: data.BeneficiaryBankGlobalIdentifier,
			PaymentType:                     "",
			TransactionType:                 transactionType.Code,
			Amount:                          amount,
			Currency:                        data.Currency,
			SwiftBICCode:                    "",
			LoanId:                          0,
		}

		_, err = c.P2PClient.Send(data.BeneficiaryBankGlobalIdentifier, "new-compliance-check", complianceCheckDTO, 0)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// ConfirmComplianceCheck handles a web GET/POST "/confirmcompliancecheck" request. For a GET, it responds with
// a view with all the information related to the compliance check and an option to confirm. On the other hand,
// POST means confirmation and that other participants in the system (originator and central bank) should be
// notified. It is done by sending a confirmation message over a p2p network. When this message is send to the
// central bank, in order for alignment with commercial banks, the message additionally contains information
// related to the compliance check as well as the policies that are applied.
func (c *ComplianceCheckController) ConfirmComplianceCheck(w http.ResponseWriter, r *http.Request) {
	if c.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodGet {

		// TODO: add logic together with a new addcompliancecheck view
		// TODO: modify json names for marshal/unmarshal

	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		complianceCheckId := r.Form.Get("id")
		complianceCheck, err := c.DB.GetComplianceCheckById(complianceCheckId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		// send a compliance check confirmation to the originator bank
		_, err = c.P2PClient.Send(complianceCheck.OriginatorBankId, "compliance-check-confirmation", common.ComplianceCheckConfirmationDTO{
			ComplianceCheckId: complianceCheck.Id,
		}, 0)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		originator, err := c.DB.GetBankClientById(complianceCheck.SenderId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		beneficiary, err := c.DB.GetBankClientById(complianceCheck.ReceiverId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		transactionType, err := c.DB.GetTransactionTypeById(complianceCheck.TransactionTypeId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		policies, err := c.DB.GetPoliciesByComplianceCheckId(complianceCheck.Id)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		// a flag indicating whether a private policy exists
		// private policies are never returned individually (nor are their details disclosed),
		// but are grouped into one private policy
		privatePolicy := false

		// convert policies into their transfer form (DTO)
		var policiesDTO []common.PolicyDTO
		for _, policy := range policies {
			if policy.Policy.IsPrivate {
				privatePolicy = true
				continue
			}

			policiesDTO = append(policiesDTO, common.PolicyDTO{
				Code:   policy.PolicyType.Code,
				Name:   policy.PolicyType.Name,
				Params: policy.Policy.Parameters,
				Owner:  policy.Policy.Owner,
			})
		}

		// grouping private policies into one
		if privatePolicy {
			policiesDTO = append(policiesDTO, common.PolicyDTO{
				Code:   "Other",
				Name:   "Internal Checks",
				Params: "",
				Owner:  config.ResolveMyGlobalIdentifier(),
			})
		}

		// send a compliance check confirmation to the beneficiary central bank
		_, err = c.P2PClient.Send(config.ResolveCBGlobalIdentifier(), "compliance-check-confirmation", common.ComplianceCheckConfirmationDTO{
			ComplianceCheckId: complianceCheck.Id,
			Data: common.ComplianceCheckConfirmationData{
				ComplianceCheck: common.ComplianceCheckDTO{
					ComplianceCheckId:               complianceCheck.Id,
					OriginatorGlobalIdentifier:      originator.GlobalIdentifier,
					OriginatorName:                  originator.Name,
					BeneficiaryGlobalIdentifier:     beneficiary.GlobalIdentifier,
					BeneficiaryName:                 beneficiary.Name,
					OriginatorBankGlobalIdentifier:  complianceCheck.OriginatorBankId,
					BeneficiaryBankGlobalIdentifier: complianceCheck.BeneficiaryBankId,
					PaymentType:                     "",
					TransactionType:                 transactionType.Code,
					Amount:                          complianceCheck.Amount,
					Currency:                        complianceCheck.Currency,
					SwiftBICCode:                    "",
					LoanId:                          complianceCheck.LoanId,
				},
				Policies: policiesDTO,
			},
		}, 0)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		go c.RulesEngine.Do(complianceCheck.Id, "interactive")

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}
