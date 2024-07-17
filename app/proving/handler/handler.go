package handler

import (
	"bisgo/app/manager"
	"bisgo/app/models"
	"bisgo/app/proving/core"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type ProvingHandler struct {
	*core.Core
	*manager.ComplianceCheckStateManager
}

func CreateProvingHandler(core *core.Core) *ProvingHandler {
	return &ProvingHandler{core, manager.CreateComplianceCheckStateManager()}
}

func (h *ProvingHandler) HandleInteractiveProof(body []byte) {
	var messageData models.InteractiveComplianceCheckProofRequest
	if err := json.Unmarshal(body, &messageData); err != nil {
		errlog.Println(err)
		return
	}

	h.DB.InsertTransactionProof(messageData.ComplianceCheckID, messageData.Value)

	policyID, err := strconv.Atoi(messageData.PolicyID)
	if err != nil {
		errlog.Println(err)
		return
	}

	var result int
	values := strings.Split(messageData.Value, ";")
	if strings.Split(values[0], ",")[1] == "0" {
		result = 1
		h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.ComplianceCheckID, policyID, false)
	} else {
		result = 2
		h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.ComplianceCheckID, policyID, true)
	}

	// originator does't need to notify its central bank about the result
	check, err := h.DB.GetComplianceCheckById(messageData.ComplianceCheckID)
	if err != nil {
		errlog.Println(err)
		return
	}

	if check.BeneficiaryBankId == config.ResolveMyGlobalIdentifier() && config.ResolveCBGlobalIdentifier() != "" {
		// notify cetntra bank
		policy, err := h.DB.GetPolicyById(policyID)
		if err != nil {
			errlog.Println(err)
			return
		}

		policyCheckResult := common.PolicyCheckResultDTO{
			ComplianceCheckId: check.Id,
			Code:              policy.PolicyType.Code,
			Name:              policy.PolicyType.Name,
			Owner:             policy.Policy.Owner,
			Result:            result,
		}

		_, err = h.P2PClient.Send(config.ResolveCBGlobalIdentifier(), "policy-check-result", policyCheckResult, 0)
		if err != nil {
			errlog.Println(err)
			return
		}
	}
}

func (h *ProvingHandler) HandleNonInteractiveProof(body []byte) {

	// Remove the leading and trailing quotes
	trimmedBody := strings.Trim(string(body), "\"")

	// Unescape the JSON string
	unescapedBody, err := strconv.Unquote("\"" + trimmedBody + "\"")
	if err != nil {
		errlog.Println(err)
		return
	}

	var messageData models.NonInteractiveComplianceCheckProofResponse
	if err := json.Unmarshal([]byte(unescapedBody), &messageData); err != nil {
		errlog.Println(err)
		return
	}

	h.DB.InsertTransactionProof(messageData.SanctionedCheckOutput.ComplianceCheckID, unescapedBody)

	fmt.Println(messageData)
	policyID, err := strconv.Atoi(messageData.SanctionedCheckOutput.PolicyID)
	if err != nil {
		errlog.Println(err)
		return
	}

	result := 0
	if messageData.Status == "Failed" {
		h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.SanctionedCheckOutput.ComplianceCheckID, policyID, true)
		result = 2
	} else {
		if messageData.SanctionedCheckOutput.NotSanctioned {
			// Passed sanction check
			h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.SanctionedCheckOutput.ComplianceCheckID, policyID, false)
			result = 1
		} else {
			// Failed sanction check
			h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.SanctionedCheckOutput.ComplianceCheckID, policyID, true)
			result = 2
		}
	}

	// originator needs to notify beneficiary with the result and send the proof through
	check, err := h.DB.GetComplianceCheckById(messageData.SanctionedCheckOutput.ComplianceCheckID)
	if err != nil {
		errlog.Println(err)
		return
	}

	if check.OriginatorBankId == config.ResolveMyGlobalIdentifier() {
		// notify beneficiary bank
		policy, err := h.DB.GetPolicyById(policyID)
		if err != nil {
			errlog.Println(err)
			return
		}

		policyCheckResult := common.PolicyCheckResultDTO{
			ComplianceCheckId: check.Id,
			Code:              policy.PolicyType.Code,
			Name:              policy.PolicyType.Name,
			Owner:             policy.Policy.Owner,
			Result:            result,
		}

		_, err = h.P2PClient.Send(check.BeneficiaryBankId, "policy-check-result", policyCheckResult, 0)
		if err != nil {
			errlog.Println(err)
			return
		}

	}
}
