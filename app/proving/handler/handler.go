package handler

import (
	"bisgo/app/models"
	"bisgo/app/proving/core"
	"bisgo/app/web/manager"
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

	values := strings.Split(messageData.Value, ";")
	if strings.Split(values[0], ",")[0] == "0" {
		h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.ComplianceCheckID, policyID, false)
	} else {
		h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.ComplianceCheckID, policyID, true)
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
