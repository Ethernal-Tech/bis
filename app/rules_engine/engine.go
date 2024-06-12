package engine

import (
	"bisgo/app/DB"
	p2pclient "bisgo/app/P2P/client"
	provingclient "bisgo/app/proving/client"
	"bisgo/common"
	"bisgo/config"
	"log"
	"strconv"
	"strings"
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

func (e *RulesEngine) Do(transactionID string, proofType string, data map[string]any) {
	policies := e.db.GetPoliciesForTransaction(transactionID)
	for _, policy := range policies {
		switch policy.PolicyType.Code {
		case "SCL":
			if !config.ResovleIsCentralBank() {
				go e.sanctionCheckList(proofType, policy.Policy.Parameters, transactionID, policy.Policy.Id, data["vm_address"].(string))
			}
		case "CFM":
			if config.ResovleIsCentralBank() {
				go e.capitalFlowManagement(proofType, policy.Policy.Parameters, transactionID)
			}
		}
	}
}

func (e *RulesEngine) sanctionCheckList(proofType string, parameters string, transactionID string, policyId int, vmAdress string) {
	if proofType == "interactive" {
		check := e.db.GetComplianceCheckByID(transactionID)
		sender := e.db.GetClientNameByID(check.SenderId)
		err := e.provingClient.SendProofRequest("interactive", transactionID, policyId, sender, vmAdress)
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else if proofType == "noninteractive" {

	}
}

func (e *RulesEngine) capitalFlowManagement(proofType string, parameters string, transactionID string) {
	if proofType == "interactive" {
		check := e.db.GetComplianceCheckByID(transactionID)
		receiver := e.db.GetClientByID(check.ReceiverId)
		amount := e.db.CheckCFM(receiver.GlobalIdentifier, e.db.GetCountryByCode(config.ResolveCountryCode()).Id)

		policyAmountStr := strings.ReplaceAll(parameters, ",", "")
		poliyAmountInt, err := strconv.Atoi(policyAmountStr)

		if err != nil {
			log.Println(err)
		}

		ratio := 3.4

		// TODO: update transaction CFM status for central bank (TransactionPolicyDB); ATM it will be updated only in commercial bank

		if float64(amount) >= float64(poliyAmountInt)*ratio {
			e.p2pClient.Send(check.BeneficiaryBankId, "cfm-result-beneficiary", any(common.CFMCheckDTO{transactionID, 2}), 0)
		} else {
			e.p2pClient.Send(check.BeneficiaryBankId, "cfm-result-beneficiary", any(common.CFMCheckDTO{transactionID, 1}), 0)
		}
	} else if proofType == "noninteractive" {

	}
}

func (e *RulesEngine) GetVMAddress() string {
	return e.provingClient.GetVMAddress()
}
