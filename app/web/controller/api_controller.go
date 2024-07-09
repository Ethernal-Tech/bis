package controller

import (
	"bisgo/common"
	"bisgo/errlog"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// GetBeneficiaryBankPolicies handles web API call ("/api/getbeneficiarybankpolicies") for obtaining the
// policies for the beneficiary bank. As input it expects a JSON object of the form (beneficiary bank
// global id, transaction type) and as a result it responds with a JSON object containing the policies.
// It internally sends a request over a p2p network to the benefiricary bank to get its policies.
func (c *APIController) GetBeneficiaryBankPolicies(w http.ResponseWriter, r *http.Request) {
	if c.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//time.Sleep(2 * time.Second)

	data := struct {
		BeneficiaryBankGlobalIdentifier string
		TransactionTypeId               string
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		errlog.Println(err)

		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	originatorBankGlobalIdentifier := c.SessionManager.GetString(r.Context(), "bankId")
	originatorJurisdiction, err := c.DB.GetBankJurisdiction(originatorBankGlobalIdentifier)
	if err != nil {
		errlog.Println(err)

		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	policyRequestDTO := common.PolicyRequestDTO{
		Jurisdiction:              originatorJurisdiction.Id,
		TransactionType:           data.TransactionTypeId,
		RequesterGlobalIdentifier: originatorBankGlobalIdentifier,
	}

	ch, err := c.P2PClient.Send(data.BeneficiaryBankGlobalIdentifier, "get-policies", policyRequestDTO, 0)
	if err != nil {
		errlog.Println(err)

		http.Error(w, "Internal Server Error", 500)
		return
	}

	responseData := (<-ch).(common.PolicyResponseDTO)

	beneficiaryJurisdiction, err := c.DB.GetBankJurisdiction(data.BeneficiaryBankGlobalIdentifier)
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

	for _, policy := range responseData.Policies {
		policyTypeId, err := c.DB.CreateOrGetPolicyType(policy.Code, policy.Name)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		_, _, err = c.DB.CreateOrUpdatePolicy(policyTypeId, policy.Owner, transactionTypeId, beneficiaryJurisdiction.Id, originatorJurisdiction.Id, policy.Params, 0)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}
	}

	response, err := json.Marshal(responseData)
	if err != nil {
		errlog.Println(fmt.Errorf("%v %w", responseData, err))

		http.Error(w, "Internal Server Error", 500)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		errlog.Println(err)

		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (controller *APIController) GetPolicy(w http.ResponseWriter, r *http.Request) {
	if controller.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	data := struct {
		BankCountry string
		PolicyId    string
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	policyId, err := strconv.Atoi(data.PolicyId)
	if err != nil {
		http.Error(w, "Failed to decode policy id", http.StatusInternalServerError)
		return
	}

	policies := controller.DB.GetPolicy(uint64(policyId))

	jsonData, err := json.Marshal(policies)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
