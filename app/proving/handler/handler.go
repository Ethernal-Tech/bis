package handler

import (
	"bisgo/app/models"
	"bisgo/app/proving/core"
	"bisgo/errlog"
	"encoding/json"
	"strconv"
	"strings"
)

type ProvingHandler struct {
	*core.Core
}

func CreateProvingHandler(core *core.Core) *ProvingHandler {
	return &ProvingHandler{core}
}

func (h *ProvingHandler) HandleInteractiveProof(body []byte) {
	var messageData models.InteractiveComplianceCheckProofRequest
	if err := json.Unmarshal(body, &messageData); err != nil {
		errlog.Println(err)
		return
	}

	h.DB.InsertTransactionProof(messageData.ComplianceCheckID, messageData.Value)

	values := strings.Split(messageData.Value, ";")
	if strings.Split(values[0], ",")[0] == "0" {
		policyId, err := strconv.Atoi(messageData.PolicyID)
		if err != nil {
			errlog.Println(err)
			return
		}

		h.DB.UpdateTransactionPolicyStatus(messageData.ComplianceCheckID, policyId, 1)
		h.DB.UpdateTransactionState(messageData.ComplianceCheckID, 4)

		// TODO: Move to sep function
		statuses := h.DB.GetTransactionPolicyStatuses(messageData.ComplianceCheckID)
		noOfPassed := 0
		for _, status := range statuses {
			if status.Status == 1 {
				noOfPassed += 1
			} else if status.Status == 2 {
				h.DB.UpdateTransactionState(messageData.ComplianceCheckID, 8)
			}
		}

		if noOfPassed == len(statuses) {
			h.DB.UpdateTransactionState(messageData.ComplianceCheckID, 7)
		}
	} else {
		policyId, err := strconv.Atoi(messageData.PolicyID)
		if err != nil {
			errlog.Println(err)
			return
		}

		h.DB.UpdateTransactionPolicyStatus(messageData.ComplianceCheckID, policyId, 2)
		h.DB.UpdateTransactionState(messageData.ComplianceCheckID, 5)
		h.DB.UpdateTransactionState(messageData.ComplianceCheckID, 8)
	}
}

func (h *ProvingHandler) HandleNonInteractiveProof() {

}
