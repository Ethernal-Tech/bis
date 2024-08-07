package handler

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/messages"
	"bisgo/app/P2P/subscribe"
	"bisgo/app/manager"
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/json"
	"errors"
	"strconv"
)

type P2PHandler struct {
	*core.Core
	*manager.ComplianceCheckStateManager
}

func CreateP2PHandler(core *core.Core) *P2PHandler {
	return &P2PHandler{core, manager.CreateComplianceCheckStateManager()}
}

// TODO: if an error occurs, the channel should be closed so that the listener does not wait forever

// TODO: all or nothing - methods should be constructed as atomic blocks (transactions), if anything fails, all changes are rolled back

// AddComplianceCheck p2p handler method, as the name suggests, adds a new compliance check obtained from
// the p2p network. It is invoked when a "new-compliance-check" message arrives from a p2p network.
func (h *P2PHandler) AddComplianceCheck(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method AddComplianceCheck failed to execute properly")

	var complianceCheck common.ComplianceCheckDTO
	err := json.Unmarshal(payload, &complianceCheck)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	transactionType, err := h.DB.GetTransactionTypeByCode(complianceCheck.TransactionType)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	// create a new user (originator) if needed
	originatorId, err := h.DB.CreateOrGetBankClient(complianceCheck.OriginatorGlobalIdentifier, complianceCheck.OriginatorName, "", complianceCheck.OriginatorBankGlobalIdentifier)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	// create a new user (beneficiary) if needed
	beneficiaryId, err := h.DB.CreateOrGetBankClient(complianceCheck.BeneficiaryGlobalIdentifier, complianceCheck.BeneficiaryName, "", complianceCheck.BeneficiaryBankGlobalIdentifier)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	// Insert potentialy sent policies by OB in the DB
	if len(complianceCheck.OBApplicabePolicies) > 0 {
		originatorJurisdiction, err := h.DB.GetBankJurisdiction(complianceCheck.OriginatorBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		beneficiaryJurisdiction, err := h.DB.GetBankJurisdiction(complianceCheck.BeneficiaryBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		for _, policy := range complianceCheck.OBApplicabePolicies {
			if policy.Owner == config.ResolveMyGlobalIdentifier() {
				continue
			}

			policyTypeId, err := h.DB.CreateOrGetPolicyType(policy.Code, policy.Name)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}

			_, _, err = h.DB.CreateOrUpdatePolicy(policyTypeId, policy.Owner, transactionType.Id, beneficiaryJurisdiction.Id, originatorJurisdiction.Id, beneficiaryJurisdiction.Id, policy.Params, 0)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}
		}
	}

	_, err = h.DB.AddComplianceCheck(models.ComplianceCheck{
		Id:                complianceCheck.ComplianceCheckId,
		OriginatorBankId:  complianceCheck.OriginatorBankGlobalIdentifier,
		BeneficiaryBankId: complianceCheck.BeneficiaryBankGlobalIdentifier,
		SenderId:          originatorId,
		ReceiverId:        beneficiaryId,
		Currency:          complianceCheck.Currency,
		Amount:            complianceCheck.Amount,
		TransactionTypeId: transactionType.Id,
		LoanId:            complianceCheck.LoanId,
	})
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	// TODO: a compliance check status manager call instead of a direct state change
	err = h.DB.UpdateComplianceCheckStatus(complianceCheck.ComplianceCheckId, 1)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	return nil
}

// GetPolicies p2p handler method handles a request for policies. It is invoked when a "get-policies" message
// arrives from the p2p network. Requests for policies can be sent by commercial banks to each other, or by
// a commercial bank to its central bank. The first case includes the second, so that the commercial bank
// always return the union of its policies and the policies of its central bank. Applicable (and returned)
// policies are determined based on the transaction type and the originating jurisdiction.
func (h *P2PHandler) GetPolicies(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method GetPolicies failed to execute properly")

	var request common.PolicyRequestDTO
	if err := json.Unmarshal(payload, &request); err != nil {
		errlog.Println(err)
		return returnErr
	}

	transactionTypeId, err := strconv.Atoi(request.TransactionType)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	// in the case of a commercial bank, it is first necessary to send a request to the central bank to obtain its policies as well
	if !config.ResolveIsCentralBank() && config.ResolveCBGlobalIdentifier() != "" {
		requestToCB := common.PolicyRequestDTO{
			Jurisdiction:              request.Jurisdiction,
			TransactionType:           request.TransactionType,
			RequesterGlobalIdentifier: config.ResolveMyGlobalIdentifier(),
		}

		ch, err := h.P2PClient.Send(config.ResolveCBGlobalIdentifier(), "get-policies", requestToCB, 0)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		responseData := (<-ch).(common.PolicyResponseDTO)

		// loop through the CB policies and insert (policy type)/policy if necessary
		for _, policy := range responseData.Policies {
			policyTypeId, err := h.DB.CreateOrGetPolicyType(policy.Code, policy.Name)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}

			// isPrivate flag is always set to 0 regardless of the policy type
			// private CB policies are always public for commercial banks in the sense that they know about their existence, but they are not familiar with the details
			_, _, err = h.DB.CreateOrUpdatePolicy(policyTypeId, policy.Owner, transactionTypeId, config.ResolveJurisdictionCode(), request.Jurisdiction, config.ResolveJurisdictionCode(), policy.Params, 0)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}
		}
	}

	var policies []models.PolicyAndItsType

	// (if) in the case of a commercial bank, all (commerical bank + CB) policies are taken
	// (else) otherwise, only policies owned by the CB are taken
	if !config.ResolveIsCentralBank() {
		policies, err = h.DB.GetAppliedPolicies(request.Jurisdiction, transactionTypeId)
	} else {
		policies, err = h.DB.GetAppliedPoliciesByOwner(config.ResolveMyGlobalIdentifier(), request.Jurisdiction, transactionTypeId)
	}

	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	var response common.PolicyResponseDTO

	// a flag indicating whether a private policy exists
	// private policies are never returned individually (nor are their details disclosed),
	// but are grouped into one private policy
	privatePolicy := false

	for _, policy := range policies {
		if policy.Policy.IsPrivate {
			privatePolicy = true
			continue
		}

		response.Policies = append(response.Policies, common.PolicyDTO{
			Code:   policy.PolicyType.Code,
			Name:   policy.PolicyType.Name,
			Params: policy.Policy.Parameters,
			Owner:  policy.Policy.Owner,
		})
	}

	// grouping private policies into one
	if privatePolicy {
		response.Policies = append(response.Policies, common.PolicyDTO{
			Code:   "Other",
			Name:   "Internal Checks",
			Params: "",
			Owner:  config.ResolveMyGlobalIdentifier(),
		})
	}

	// sending policies over a p2p network back to the requesting bank
	_, err = h.P2PClient.Send(request.RequesterGlobalIdentifier, "policies", response, messageID)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	return nil
}

// ReceivePolicies p2p handler method handles the response to the "get-policies" request. It is invoked
// when a "policies" message arrives from a p2p network. Policies (that is, response on a "get-policies"
// request) can be sent by commercial bank to another, or by a central bank to commercial one.
func (h *P2PHandler) ReceivePolicies(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method ReceivePolicies failed to execute properly")

	channel, _ := messages.LoadChannel(messageID)
	defer messages.RemoveChannel(messageID)

	var messageData common.PolicyResponseDTO
	err := json.Unmarshal(payload, &messageData)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	channel <- messageData

	return nil
}

// ConfirmComplianceCheck p2p handler method, as the name suggests, confirms a selected compliance check and starts the
// rules engine for an originator bank and beneficiary central bank. It is invoked when a "compliance-check-confirmation"
// message arrives from a p2p network. Additionally, it also aligns the central bank with the rest of the system.
func (h *P2PHandler) ConfirmComplianceCheck(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method ConfirmComplianceCheck failed to execute properly")

	var complianceCheckConfirmation common.ComplianceCheckConfirmationDTO

	err := json.Unmarshal(payload, &complianceCheckConfirmation)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	complianceCheck := complianceCheckConfirmation.Data.ComplianceCheck
	policies := complianceCheckConfirmation.Data.Policies

	// since the central bank may not be aware of the compliance check and the information about it, it is necessary to carry
	// out its alignment with the rest of the system (commercial banks), the following needs to be done:
	// 1. potentially create (if needed) originator (sender)
	// 2. potentially create (if needed) beneficiary (receiver)
	// 3. potentially create (if needed) policy types
	// 4. potentially create (if needed) applicable policies
	// 5. create compliance check
	if config.ResolveIsCentralBank() {

		originatorId, err := h.DB.CreateOrGetBankClient(complianceCheck.OriginatorGlobalIdentifier, complianceCheck.OriginatorName, "", complianceCheck.OriginatorBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		beneficiaryId, err := h.DB.CreateOrGetBankClient(complianceCheck.BeneficiaryGlobalIdentifier, complianceCheck.BeneficiaryName, "", complianceCheck.BeneficiaryBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		originatorJurisdiction, err := h.DB.GetBankJurisdiction(complianceCheck.OriginatorBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		beneficiaryJurisdiction, err := h.DB.GetBankJurisdiction(complianceCheck.BeneficiaryBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		transactionType, err := h.DB.GetTransactionTypeByCode(complianceCheck.TransactionType)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}

		for _, policy := range policies {
			if policy.Owner == config.ResolveMyGlobalIdentifier() {
				continue
			}

			policyTypeId, err := h.DB.CreateOrGetPolicyType(policy.Code, policy.Name)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}

			_, _, err = h.DB.CreateOrUpdatePolicy(policyTypeId, policy.Owner, transactionType.Id, beneficiaryJurisdiction.Id, originatorJurisdiction.Id, beneficiaryJurisdiction.Id, policy.Params, 0)
			if err != nil {
				errlog.Println(err)
				return returnErr
			}
		}

		_, err = h.DB.AddComplianceCheck(models.ComplianceCheck{
			Id:                complianceCheck.ComplianceCheckId,
			OriginatorBankId:  complianceCheck.OriginatorBankGlobalIdentifier,
			BeneficiaryBankId: complianceCheck.BeneficiaryBankGlobalIdentifier,
			SenderId:          originatorId,
			ReceiverId:        beneficiaryId,
			Currency:          complianceCheck.Currency,
			Amount:            complianceCheck.Amount,
			TransactionTypeId: transactionType.Id,
			LoanId:            complianceCheck.LoanId,
		})
		if err != nil {
			errlog.Println(err)
			return returnErr
		}
	}

	go h.RulesEngine.Do(complianceCheckConfirmation.ComplianceCheckId, config.ResolveRuleEngineProofType())

	return nil
}

// ProcessPolicyCheckResult p2p handler method processes the result of the policy check. It is invoked
// when a "policy-check-result" message arrives from a p2p network.
func (h *P2PHandler) ProcessPolicyCheckResult(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method ProcessPolicyCheckResult failed to execute properly")

	var policyCheckResult common.PolicyCheckResultDTO
	err := json.Unmarshal(payload, &policyCheckResult)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	complianceCheck, err := h.DB.GetComplianceCheckById(policyCheckResult.ComplianceCheckId)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	policyType, err := h.DB.GetPolicyTypeByCode(policyCheckResult.Code)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	originatorJurisdiction, err := h.DB.GetBankJurisdiction(complianceCheck.OriginatorBankId)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	policy, err := h.DB.GetPolicyToProcessItsCheckResult(policyType.Id, policyCheckResult.Owner, complianceCheck.TransactionTypeId, originatorJurisdiction.Id, 0)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	err = h.DB.UpdatePolicyStatus(complianceCheck.Id, policy.Id, policyCheckResult.Result)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	if policyCheckResult.Proof != "" {
		h.DB.InsertTransactionProof(complianceCheck.Id, policyCheckResult.Proof)
	}

	if policyCheckResult.ForwardTo != "" {
		forwardToGlobalIdentifier := policyCheckResult.ForwardTo
		policyCheckResult.ForwardTo = ""
		_, err = h.P2PClient.Send(forwardToGlobalIdentifier, "policy-check-result", policyCheckResult, 0)
		if err != nil {
			errlog.Println(err)
			return returnErr
		}
	}

	return nil
}

// ProcessMpcStartSignal p2p handler method processes the MPC start signal. It is invoked
// when a "mpc-start-signal" message arrives from a p2p network.
func (h *P2PHandler) ProcessMpcStartSignal(messageID int, payload []byte) error {
	returnErr := errors.New("p2p handler method ProcessMpcStartSignal failed to execute properly")

	var signal common.MPCStartSignalDTO
	err := json.Unmarshal(payload, &signal)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	err = subscribe.StoreAndNotify(subscribe.SCLServerStarted, signal)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

	return nil
}
