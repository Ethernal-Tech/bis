package engine

import (
	"bisgo/app/DB"
	"bisgo/app/proving/client"
	"log"
)

type RulesEngine struct {
	db            *DB.DBHandler
	provingClient *client.ProvingClient
}

var engine RulesEngine

func init() {
	engine = RulesEngine{
		db:            DB.GetDBHandler(),
		provingClient: client.GetProvingClient(),
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
			go e.sanctionCheckList(proofType, policy.Policy.Parameters, transactionID, policy.Policy.Id, data["vm_address"].(string))
		case "CFM":
			go e.capitalFlowManagement(proofType, policy.Policy.Parameters, transactionID)
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

	} else if proofType == "noninteractive" {

	}
}

func (e *RulesEngine) GetVMAddress() string {
	return e.provingClient.GetVMAddress()
}
