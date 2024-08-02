package controller

import (
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// GetPolicies handles the web API call to "/api/getpolicies" for retrieving applicable policies
// for a given transaction type and beneficiary jurisdiction.
//
// Input, JSON object containing:
//   - bb_gid : string // beneficiary commercial bank global id
//   - tx_type : string // transaction type id
//
// Output, JSON object (array/list) containing applicable policies, union of the following:
//   - Originator commercial bank policies
//   - Beneficiary commercial bank policies
//   - Beneficiary central bank policies
//
// GetPolicies internally sends a request over a p2p network to the beneficiary commerical bank
// with the intention of obtaining policies of the beneficiary side.
func (c *APIController) GetPolicies(w http.ResponseWriter, r *http.Request) {
	if c.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//time.Sleep(2 * time.Second)

	data := struct {
		BeneficiaryBankGlobalIdentifier string `json:"bb_gid"`
		TransactionTypeId               string `json:"tx_type"`
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

		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}

	// TODO: transaction type ids are internal to each bank and may differ, we should send transaction type codes instead

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

		_, _, err = c.DB.CreateOrUpdatePolicy(policyTypeId, policy.Owner, transactionTypeId, beneficiaryJurisdiction.Id, originatorJurisdiction.Id, beneficiaryJurisdiction.Id, policy.Params, 0)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}
	}

	// TODO: send a request to the originating CB to pull and save its policies as well

	// originator commercial bank may also have policies, so we take them
	originatorPolicies, err := c.DB.GetBankPolicies(config.ResolveMyGlobalIdentifier(), originatorJurisdiction.Id, beneficiaryJurisdiction.Id, transactionTypeId)
	if err != nil {
		errlog.Println(err)

		http.Error(w, "Internal Server Error", 500)
		return
	}

	// append originator policies to the received (beneficiary side) policies
	for _, policy := range originatorPolicies {
		responseData.Policies = append(responseData.Policies, common.PolicyDTO{
			Code:   policy.PolicyType.Code,
			Name:   policy.PolicyType.Name,
			Params: policy.Policy.Parameters,
			Owner:  policy.Policy.Owner,
		})
	}

	// TODO: append originator CB policies to the previous list (use c.DB.GetAllOriginatorBankPolicies instead of c.DB.GetBankPolicies)

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
