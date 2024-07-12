package handler

import (
	"bisgo/app/models"
	"bisgo/app/proving/core"
	"bisgo/app/web/manager"
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

	if check.BeneficiaryBankId == config.ResolveMyGlobalIdentifier() {
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
	fmt.Println("called /proof/noninteractive")

	// Remove the leading and trailing quotes
	trimmedBody := strings.Trim(string(body), "\"")

	// Unescape the JSON string
	unescapedBody, err := strconv.Unquote("\"" + trimmedBody + "\"")
	if err != nil {
		errlog.Println(err)
		return
	}

	var messageData models.NonInteractiveComplianceCheckProofRequest
	if err := json.Unmarshal([]byte(unescapedBody), &messageData); err != nil {
		errlog.Println(err)
		return
	}

	fmt.Println(messageData)

	if messageData.Status == "Failed" {
		// TODO: Should we rerun?
	} else {
		if messageData.SanctionedCheckOutput.NotSanctioned {
			// Passed sanction check

		} else {
			// Failed sanction check
		}
	}
}
