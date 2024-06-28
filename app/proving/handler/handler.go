package handler

import (
	"bisgo/app/models"
	"bisgo/app/proving/core"
	"bisgo/app/web/manager"
	"bisgo/errlog"
	"encoding/json"
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

func (h *ProvingHandler) HandleNonInteractiveProof() {

}
