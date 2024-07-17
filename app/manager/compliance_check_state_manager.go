package manager

import (
	"bisgo/app/DB"
)

type ComplianceCheckStateManager struct{}

func CreateComplianceCheckStateManager() *ComplianceCheckStateManager {
	return &ComplianceCheckStateManager{}
}

func (*ComplianceCheckStateManager) UpdateComplianceCheckPolicyStatus(db *DB.DBHandler, complianceCheckID string, policyID int, isFailed bool) {
	if isFailed {
		db.UpdateTransactionPolicyStatus(complianceCheckID, policyID, 2)
	} else {
		db.UpdateTransactionPolicyStatus(complianceCheckID, policyID, 1)
	}

	statuses := db.GetTransactionPolicyStatuses(complianceCheckID)
	noOfPassed := 0
	noOfCompleted := 0
	for _, status := range statuses {
		if status.Status == 1 {
			noOfPassed += 1
			noOfCompleted += 1
		} else if status.Status == 2 {
			noOfCompleted += 1
			db.UpdateTransactionState(complianceCheckID, 5)
			db.UpdateTransactionState(complianceCheckID, 8)
		}
	}

	if noOfPassed == len(statuses) {
		db.UpdateTransactionState(complianceCheckID, 4)
		db.UpdateTransactionState(complianceCheckID, 7)
		// TODO: Notify CB about the compliance check completion
		// if !config.ResovleIsCentralBank() {
		// }
	} else if noOfCompleted == len(statuses) {
		// TODO: Notify CB about the compliance check completion
		// if !config.ResovleIsCentralBank() {
		// }
	}
}
