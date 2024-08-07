package handler

import (
	"bisgo/app/models"
	"bisgo/app/proving/core"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
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

	policyID, err := strconv.Atoi(messageData.PolicyID)
	if err != nil {
		errlog.Println(err)
		return
	}

	var result int
	values := strings.Split(messageData.Value, ";")
	if strings.Split(values[0], ",")[1] == "0" {
		result = 1
	} else {
		result = 2
	}

	err = h.DB.UpdatePolicyStatus(messageData.ComplianceCheckID, policyID, result)
	if err != nil {
		errlog.Println(err)
		return
	}
	state, went, err := h.ComplianceCheckStateManager.Transition(messageData.ComplianceCheckID)
	fmt.Println(state, went, err)

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

	// // Remove the leading and trailing quotes
	// trimmedBody := strings.Trim(string(body), "\"")

	// // Unescape the JSON string
	// unescapedBody, err := strconv.Unquote("\"" + trimmedBody + "\"")
	// if err != nil {
	// 	errlog.Println(err)
	// 	return
	// }

	var messageData models.NonInteractiveComplianceCheckProofResponse
	if err := json.Unmarshal(body, &messageData); err != nil {
		errlog.Println(err)
		return
	}

	marshaledSanctionedCheckOutput, err := json.Marshal(messageData.SanctionedCheckOutput)
	if err != nil {
		errlog.Println(err)
		return
	}

	h.DB.InsertTransactionProof(messageData.SanctionedCheckOutput.ComplianceCheckID, string(marshaledSanctionedCheckOutput))

	policyID, err := strconv.Atoi(messageData.SanctionedCheckOutput.PolicyID)
	if err != nil {
		errlog.Println(err)
		return
	}

	check, err := h.DB.GetComplianceCheckById(messageData.SanctionedCheckOutput.ComplianceCheckID)
	if err != nil {
		errlog.Println(err)
		return
	}

	// TODO: Remove, needed for smart contract testing
	{ // testing
		fmt.Println(messageData.SanctionedCheckOutput.Proof[0:4])
		fmt.Println(common.IntArrayToHexString(messageData.SanctionedCheckOutput.Proof[0:4]))
		b := make([]byte, 0)
		for i := 0; i < 4; i++ {
			b = append(b, byte(messageData.SanctionedCheckOutput.Proof[i]))
		}
		fmt.Println(hex.EncodeToString(b))
	}

	result := 0
	if messageData.Status == "Failed" {
		// err = h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.SanctionedCheckOutput.ComplianceCheckID, policyID, true, "Proof genereation failed")
		// if err != nil {
		// 	errlog.Println(err)
		// 	return
		// }
		result = 2
	} else {
		if messageData.SanctionedCheckOutput.NotSanctioned {
			// Passed sanction check
			// err = h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.SanctionedCheckOutput.ComplianceCheckID, policyID, false, "")
			// if err != nil {
			// 	errlog.Println(err)
			// 	return
			// }
			result = 1
		} else {
			// Failed sanction check
			sender, err := h.DB.GetBankClientById(check.SenderId)
			if err != nil {
				errlog.Println(err)
				return
			}

			receiver, err := h.DB.GetBankClientById(check.ReceiverId)
			if err != nil {
				errlog.Println(err)
				return
			}

			beneficiaryBank, err := h.DB.GetBankByGlobalIdentifier(check.BeneficiaryBankId)
			if err != nil {
				errlog.Println(err)
				return
			}

			names := []string{sender.Name, receiver.Name, beneficiaryBank.Name}

			// Iterate over elements to determine which entity is sanctioned
			description := "Sanctioned hit on"
			for _, entity := range messageData.SanctionedCheckInput.ParticipantsList {
				for _, sanctioned := range messageData.SanctionedCheckInput.PubSanctionsList {
					if reflect.DeepEqual(entity, sanctioned) {
						for i, name := range names {
							if reflect.DeepEqual(common.HashName(name), entity) {
								if i == 0 {
									description = strings.Join([]string{description}, " originator")
								} else if i == 1 {
									if len(description) == 0 {
										description = strings.Join([]string{description}, " beneficiary")
									} else {
										description = strings.Join([]string{description}, " and beneficiary")
									}
								} else {
									if len(description) == 0 {
										description = strings.Join([]string{description}, " beneficiary bank")
									} else {
										description = strings.Join([]string{description}, " and beneficiary bank")
									}
								}
							}
						}
					}
				}
			}
			// err = h.ComplianceCheckStateManager.UpdateComplianceCheckPolicyStatus(h.DB, messageData.SanctionedCheckOutput.ComplianceCheckID, policyID, true, description)
			// if err != nil {
			// 	errlog.Println(err)
			// 	return
			// }
			result = 2
		}
	}

	policy, err := h.DB.GetPolicyById(policyID)
	if err != nil {
		errlog.Println(err)
		return
	}

	// originator needs to notify beneficiary with the result and send the proof through
	if check.OriginatorBankId == config.ResolveMyGlobalIdentifier() {
		// notify beneficiary bank
		policyCheckResult := common.PolicyCheckResultDTO{
			ComplianceCheckId: check.Id,
			Code:              policy.PolicyType.Code,
			Name:              policy.PolicyType.Name,
			Owner:             policy.Policy.Owner,
			Result:            result,
			Proof:             string(marshaledSanctionedCheckOutput),
		}

		_, err = h.P2PClient.Send(check.BeneficiaryBankId, "policy-check-result", policyCheckResult, 0)
		if err != nil {
			errlog.Println(err)
			return
		}
	}
}
