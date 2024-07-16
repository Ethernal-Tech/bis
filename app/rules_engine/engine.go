package engine

import (
	"bisgo/app/DB"
	p2pclient "bisgo/app/P2P/client"
	"bisgo/app/P2P/subscribe"
	"bisgo/app/models"
	provingclient "bisgo/app/proving/client"
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RulesEngine struct {
	db            *DB.DBHandler
	provingClient *provingclient.ProvingClient
	p2pClient     *p2pclient.P2PClient
}

var engine RulesEngine

func init() {
	engine = RulesEngine{
		db:            DB.GetDBHandler(),
		provingClient: provingclient.GetProvingClient(),
		p2pClient:     p2pclient.GetP2PClient(),
	}
}

func GetRulesEngine() *RulesEngine {
	return &engine
}

func rulesEngineStartLog(complianceCheckId string) {
	dateTime := strings.Split(time.Now().String(), ".")[0]

	fmt.Printf("\033[34m%v INFO\033[0m: Rules engine for compliance check %v is started\n", dateTime, complianceCheckId)
}

func (e *RulesEngine) Do(complianceCheckId string, proofType string) {
	// due to the system running locally, there is almost no delay, so we
	// need to introduce it somehow thus it feels like a distributed system
	time.Sleep(2 * time.Second)

	rulesEngineStartLog(complianceCheckId)

	complianceCheck, err := e.db.GetComplianceCheckById(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return
	}

	// all policies that apply for a given compliance check
	allPolicies, err := e.db.GetPoliciesByComplianceCheckId(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return
	}

	// policies that are the responsibility of the current bank and that it should process
	var myPolicies []models.PolicyAndItsType

	for _, policy := range allPolicies {
		if policy.Policy.Owner == config.ResolveMyGlobalIdentifier() {
			myPolicies = append(myPolicies, policy)
			continue
		}

		// it may happen that the current (commercial) bank is not the owner of the policy, but it is an SCL policy
		if !config.ResolveIsCentralBank() && policy.PolicyType.Code == "SCL" {
			myPolicies = append(myPolicies, policy)
		}
	}

	if proofType == "interactive" {
		e.doInteractive(complianceCheck, myPolicies)
	} else {
		e.doNonInteractive(complianceCheck, myPolicies)
	}
}

func (e *RulesEngine) doInteractive(complianceCheck models.ComplianceCheck, policies []models.PolicyAndItsType) {
	var privatePolicies []models.PolicyAndItsType

	for _, policy := range policies {
		if policy.Policy.IsPrivate {
			privatePolicies = append(privatePolicies, policy)
		} else {
			if policy.PolicyType.Code == "SCL" {
				go e.interactiveSanctionCheckList(complianceCheck, policy)
			} else if policy.PolicyType.Code == "CFM" {
				go e.interactiveCapitalFlowManagement(complianceCheck, policy)
			}
		}
	}

	if len(privatePolicies) != 0 {
		go e.interactivePrivatePolicy(complianceCheck, privatePolicies)
	}
}

func (e *RulesEngine) interactivePrivatePolicy(complianceCheck models.ComplianceCheck, policies []models.PolicyAndItsType) {
	// loop through private policies and set their status to successful (1)
	// they are not checked in any way
	for _, policy := range policies {
		err := e.db.UpdatePolicyStatus(complianceCheck.Id, policy.Policy.Id, 1)
		if err != nil {
			errlog.Println(err)
			return
		}
	}

	policyCheckResult := common.PolicyCheckResultDTO{
		ComplianceCheckId: complianceCheck.Id,
		Code:              "Other",
		Name:              "Internal Checks",
		Owner:             config.ResolveMyGlobalIdentifier(),
		Result:            1,
	}

	// if the current bank is the central one, the p2p message should be sent to the beneficiary and originator bank
	// otherwise, the p2p message should be sent to the central and originator bank
	if config.ResolveIsCentralBank() {
		// send policy check result to the beneficiary bank
		_, err := e.p2pClient.Send(complianceCheck.BeneficiaryBankId, "policy-check-result", policyCheckResult, 0)
		if err != nil {
			errlog.Println(err)
			return
		}

		// send policy check result to the originator bank
		_, err = e.p2pClient.Send(complianceCheck.OriginatorBankId, "policy-check-result", policyCheckResult, 0)
		if err != nil {
			errlog.Println(err)
			return
		}
	} else {
		// send policy check result to the beneficiary central bank
		if config.ResolveCBGlobalIdentifier() != "" {
			_, err := e.p2pClient.Send(config.ResolveCBGlobalIdentifier(), "policy-check-result", policyCheckResult, 0)
			if err != nil {
				errlog.Println(err)
				return
			}
		}

		// send policy check result to the originator bank
		_, err := e.p2pClient.Send(complianceCheck.OriginatorBankId, "policy-check-result", policyCheckResult, 0)
		if err != nil {
			errlog.Println(err)
			return
		}
	}
}

func (e *RulesEngine) interactiveSanctionCheckList(complianceCheck models.ComplianceCheck, policy models.PolicyAndItsType) {
	// check whether the current bank is originator or beneficiary bank
	// depending on that, the client (if) or server (else) side of the SCL MPC protocol is executed
	if complianceCheck.OriginatorBankId == config.ResolveMyGlobalIdentifier() {
		// originator bank (client side of the SCL MPC protocol)

		subscription, err := subscribe.Subscribe(subscribe.SCLServerStarted, func(messages []any) (any, bool) {
			for _, message := range messages {
				if message.(common.MPCStartSignalDTO).ComplianceCheckId == complianceCheck.Id {
					return message, true
				}
			}

			return nil, false
		})
		if err != nil {
			errlog.Println(err)
			return
		}

		var signal common.MPCStartSignalDTO

		for {
			signal = (<-subscription.NotifyCh).(common.MPCStartSignalDTO)

			if signal.ComplianceCheckId == complianceCheck.Id {
				err = subscribe.Unsubscribe(subscribe.SCLServerStarted, subscription)
				if err != nil {
					errlog.Println(err)
					return
				}

				break
			}
		}

		err = e.provingClient.SendProofRequest("interactive", complianceCheck.Id, policy.Policy.Id, signal.VMAddress)
		if err != nil {
			errlog.Println(err)
			return
		}
	} else {
		// beneficiary bank (server side of the SCL MPC protocol)

		err := e.provingClient.SendProofRequest("interactive", complianceCheck.Id, policy.Policy.Id, "")
		if err != nil {
			errlog.Println(err)
			return
		}

		signal := common.MPCStartSignalDTO{
			ComplianceCheckId: complianceCheck.Id,
			VMAddress:         e.provingClient.GetVMAddress(),
		}

		_, err = e.p2pClient.Send(complianceCheck.OriginatorBankId, "mpc-start-signal", signal, 0)
		if err != nil {
			errlog.Println(err)
			return
		}
	}
}

func (e *RulesEngine) interactiveCapitalFlowManagement(complianceCheck models.ComplianceCheck, policy models.PolicyAndItsType) {
	originatorJurisdiction, err := e.db.GetBankJurisdiction(complianceCheck.OriginatorBankId)
	if err != nil {
		errlog.Println(err)
		return
	}

	allIncomingComplianceChecks, err := e.db.GetAllSuccessfulComplianceChecks(complianceCheck.ReceiverId, originatorJurisdiction.Id, 1)
	if err != nil {
		errlog.Println(err)
		return
	}

	allOutcomingComplianceChecks, err := e.db.GetAllSuccessfulComplianceChecks(complianceCheck.ReceiverId, originatorJurisdiction.Id, 2)
	if err != nil {
		errlog.Println(err)
		return
	}

	loanDrawdown, err := e.db.GetTransactionTypeByCode("DDWN")
	if err != nil {
		errlog.Println(err)
		return
	}

	loanRepayment, err := e.db.GetTransactionTypeByCode("PPAY")
	if err != nil {
		errlog.Println(err)
		return
	}

	var totalIncomingAmount float64

	// filter incoming compliance checks to those of the loan drawdown transaction type and take their amount
	for _, incomingComplianceCheck := range allIncomingComplianceChecks {
		if incomingComplianceCheck.TransactionTypeId == loanDrawdown.Id {
			totalIncomingAmount += float64(incomingComplianceCheck.Amount)
		}
	}

	var totalOutcomingAmount float64

	// filter outcoming compliance checks to those of the loan repayment transaction type and take their amount
	for _, outcomingComplianceCheck := range allOutcomingComplianceChecks {
		if outcomingComplianceCheck.TransactionTypeId == loanRepayment.Id {
			totalOutcomingAmount += float64(outcomingComplianceCheck.Amount)
		}
	}

	// add the amount of compliance check being checked
	totalIncomingAmount += float64(complianceCheck.Amount)

	// TODO: move ratio to the env variable
	ratio := 3.4

	// convert the amount from the currency of the originatory's country to the beneficiary's currency
	totalIncomingAmount *= ratio

	amount := totalIncomingAmount - totalOutcomingAmount

	limit, err := strconv.Atoi(policy.Policy.Parameters)
	if err != nil {
		errlog.Println(err)
		return
	}

	var policyCheckResult common.PolicyCheckResultDTO

	// (if) CFM policy check successful
	// (else) otherwise, unsuccessful
	if amount <= float64(limit) {
		err := e.db.UpdatePolicyStatus(complianceCheck.Id, policy.Policy.Id, 1)
		if err != nil {
			errlog.Println(err)
			return
		}

		policyCheckResult = common.PolicyCheckResultDTO{
			ComplianceCheckId: complianceCheck.Id,
			Code:              "CFM",
			Name:              "Capital Flow Management",
			Owner:             config.ResolveMyGlobalIdentifier(),
			Result:            1,
		}
	} else {
		err := e.db.UpdatePolicyStatus(complianceCheck.Id, policy.Policy.Id, 2)
		if err != nil {
			errlog.Println(err)
			return
		}

		policyCheckResult = common.PolicyCheckResultDTO{
			ComplianceCheckId: complianceCheck.Id,
			Code:              "CFM",
			Name:              "Capital Flow Management",
			Owner:             config.ResolveMyGlobalIdentifier(),
			Result:            2,
		}
	}

	// send CFM policy check result to the beneficiary bank
	_, err = e.p2pClient.Send(complianceCheck.BeneficiaryBankId, "policy-check-result", policyCheckResult, 0)
	if err != nil {
		errlog.Println(err)
		return
	}

	// send CFM policy check result to the originator bank
	_, err = e.p2pClient.Send(complianceCheck.OriginatorBankId, "policy-check-result", policyCheckResult, 0)
	if err != nil {
		errlog.Println(err)
		return
	}
}

func (e *RulesEngine) doNonInteractive(complianceCheck models.ComplianceCheck, policies []models.PolicyAndItsType) {

}
