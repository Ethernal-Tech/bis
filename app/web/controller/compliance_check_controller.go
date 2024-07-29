package controller

import (
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

// ComplianceCheckIndex returns a partial view with filters for users to filter compliance checks.
func (c *ComplianceCheckController) ComplianceCheckIndex(w http.ResponseWriter, r *http.Request) {
	if c.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}
	http.ServeFile(w, r, "./app/web/static/views/complianceCheckIndex.html")
}

// ComplianceChecks handles a web POST "/compliancechecks" request. It responds with a view (HTML partial) containing all
// compliance checks associated with the current bank, as well as all tools for their successful management.
func (c *ComplianceCheckController) ComplianceCheck(w http.ResponseWriter, r *http.Request) {
	if c.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var searchModel models.SearchModel
	err := json.NewDecoder(r.Body).Decode(&searchModel)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	viewData := map[string]any{}
	var transactions []models.TransactionModel
	if c.SessionManager.GetBool(r.Context(), "isCentralBank") {
		var countryId string
		transactions, countryId = c.DB.GetCentralBankTransactions(c.SessionManager.Get(r.Context(), "bankId").(string), searchModel)
		viewData["countryId"] = countryId
	} else {
		transactions = c.DB.GetCommercialBankTransactions(c.SessionManager.Get(r.Context(), "bankId").(string), searchModel)
	}

	viewData["bankName"] = c.SessionManager.GetString(r.Context(), "bankName")
	viewData["transactions"] = transactions
	viewData["country"] = c.SessionManager.GetString(r.Context(), "country")
	viewData["isCentralBank"] = c.SessionManager.GetBool(r.Context(), "isCentralBank")

	ts, err := template.ParseFiles("./app/web/static/views/complianceCheck.html")
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
	// http.ServeFile(w, r, "./app/web/static/views/compliancecheck.html")
}

// AddComplianceCheck handles a web GET/POST "/addcompliancecheck" request. For a GET, it responds with a view for a new
// compliance check creation. On the other hand, POST indicates confirmation and that a compliance check should be created.
// As part of this process, a new compliance check is also sent over p2p network to the beneficiary bank.
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
			Parameter                       string `json:"parameter"`
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
		min := 1000000
		max := 9999999
		loanID := rand.Intn(max-min) + min

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
			LoanId:            loanID,
		}

		complianceCheckId, err := c.DB.AddComplianceCheck(complianceCheck)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		if data.Parameter != "" {
			err := c.handleAdditionalParameter(data.Parameter, complianceCheckId, transactionTypeId)
			if err != nil {
				errlog.Println(err)

				http.Error(w, "Internal Server Error", 500)
				return
			}
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

		// Query OB's applicable policies so they can be sent and displayed
		// for the BB's compliance check history
		obApplicablePolicies, err := c.DB.GetAppliedPolicies(config.ResolveJurisdictionCode(), transactionTypeId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}
		obPolicies := make([]common.PolicyDTO, 0)

		if len(obApplicablePolicies) > 0 {
			beneficiaryJurisdiction, err := c.DB.GetBankJurisdiction(data.BeneficiaryBankGlobalIdentifier)
			if err != nil {
				errlog.Println(err)

				http.Error(w, "Internal Server Error", 500)
				return
			}

			for _, obPolicy := range obApplicablePolicies {
				if obPolicy.Policy.BeneficiaryJurisdictionId == beneficiaryJurisdiction.Id {
					obPolicies = append(obPolicies, common.PolicyDTO{
						Code:   obPolicy.PolicyType.Code,
						Name:   obPolicy.PolicyType.Name,
						Params: obPolicy.Policy.Parameters,
						Owner:  obPolicy.Policy.Owner,
					})
				}
			}
		}

		complianceCheckDTO := common.ComplianceCheckDTO{
			ComplianceCheckId:               complianceCheckId,
			OriginatorGlobalIdentifier:      data.OriginatorGlobalIdentifier,
			OriginatorName:                  data.OriginatorName,
			BeneficiaryGlobalIdentifier:     data.BeneficiaryGlobalIdentifier,
			BeneficiaryName:                 data.BeneficiaryName,
			OriginatorBankGlobalIdentifier:  originatorBankGlobalIdentifier,
			BeneficiaryBankGlobalIdentifier: data.BeneficiaryBankGlobalIdentifier,
			PaymentType:                     "",
			TransactionType:                 transactionType.Code,
			Amount:                          amount,
			Currency:                        data.Currency,
			SwiftBICCode:                    "",
			LoanId:                          loanID,
			OBApplicabePolicies:             obPolicies,
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
		if config.ResolveCBGlobalIdentifier() != "" {
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
		}

		go c.RulesEngine.Do(complianceCheck.Id, config.ResolveRuleEngineProofType())

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (c *ComplianceCheckController) handleAdditionalParameter(additionalParameter string, complianceCheckId string, transactionTypeId int) error {
	returnErr := errors.New("failed to add additional parameters to policies")

	transactionTypes := c.DB.GetTransactionTypes()
	for _, transactionType := range transactionTypes {
		if transactionType.Id == transactionTypeId && transactionType.Code == "SECU" {
			// There are 3 different amounts defined in the process of compliance check creation
			// Acquisition amount (A): The total amount for acquiring securities
			// Offset amount (B): The amount being offset in the transaction
			// Payment amount (C): The actual amount being transferred (C = A - B)
			// For the purposes of the checks we will pass C and A so we can calculate B

			A, err := strconv.Atoi(additionalParameter)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}

			check, err := c.DB.GetComplianceCheckById(complianceCheckId)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}

			B := A - check.Amount

			applicablePolices, err := c.DB.GetPoliciesByComplianceCheckId(complianceCheckId)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}

			for _, policy := range applicablePolices {
				if policy.PolicyType.Code == "AMT" {
					// AMT policy should check the Acquisition amount (A)
					err = c.DB.UpdateTransactionPolicyAdditionalParameters(complianceCheckId, policy.Policy.Id, additionalParameter)
					if err != nil {
						errlog.Println(err)
						return returnErr
					}
					continue
				}

				if policy.PolicyType.Code == "NETT" {
					// NETT policy should check the Offset amount (B)
					err = c.DB.UpdateTransactionPolicyAdditionalParameters(complianceCheckId, policy.Policy.Id, fmt.Sprintf("%d", B))
					if err != nil {
						errlog.Println(err)
						return returnErr
					}
					continue
				}
			}
		}
	}
	return nil
}
