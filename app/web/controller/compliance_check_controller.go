package controller

import (
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

// AddComplianceCheck handles a web GET/POST "/addcompliancecheck" request. For a GET, it responds with a view
// for a new compliance check creation. On the other hand, POST indicates confirmation and that a compliance
// check should be created. As part of this process, a new compliance check is also sent over the p2p network
// to the beneficiary bank.
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

		min := 1000000
		max := 9999999
		loanId := rand.Intn(max-min) + min

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
			LoanId:            loanId,
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

		allPolicies, err := c.DB.GetRoutePolicies(config.ResolveMyGlobalIdentifier(), data.BeneficiaryBankGlobalIdentifier, transactionTypeId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		// a flag indicating whether a private policy exists
		// private policies are never send individually (nor their details disclosed),
		// but are grouped into one private policy
		privatePolicy := false

		var policies []common.PolicyDTO

		for _, policy := range allPolicies {
			if policy.IsPrivate {
				privatePolicy = true
				continue
			}

			// only policies owned by the originating bank or its central bank are sent
			if policy.Owner == config.ResolveMyGlobalIdentifier() || policy.Owner == config.ResolveCBGlobalIdentifier() {

				policyType, err := c.DB.GetPolicyTypeById(policy.PolicyTypeId)
				if err != nil {
					errlog.Println(err)

					http.Error(w, "Internal Server Error", 500)
					return
				}

				policies = append(policies, common.PolicyDTO{
					Code:   policyType.Code,
					Name:   policyType.Name,
					Params: policy.Parameters,
					Owner:  policy.Owner,
				})
			}
		}

		// grouping private policies into one
		if privatePolicy {
			policies = append(policies, common.PolicyDTO{
				Code:   "Other",
				Name:   "Internal Checks",
				Params: "",
				Owner:  config.ResolveMyGlobalIdentifier(),
			})
		}

		dto := common.ComplianceCheckAndPoliciesDTO{
			ComplianceCheck: common.ComplianceCheckDTO{
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
				LoanId:                          loanId,
			},
			Policies: policies,
		}

		_, err = c.P2PClient.Send(data.BeneficiaryBankGlobalIdentifier, "new-compliance-check", dto, 0)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// ConfirmComplianceCheck handles a web GET/POST "/confirmcompliancecheck" request. For a GET, it responds with
// a view containing all the information related to the compliance check and confirm option. On the other hand
// POST means confirmation and that other participants in the system (originator and central banks) should be
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
		// private policies are never send individually (nor their details disclosed),
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
				Data: common.ComplianceCheckAndPoliciesDTO{
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
