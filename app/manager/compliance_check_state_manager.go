package manager

import (
	"bisgo/app/DB"
	"bisgo/errlog"
	"errors"
)

// ComplianceCheckStateManager manages the states of a compliance check during its lifecycle. All possible states of
// the compliance check are shown below. The order in which they are displayed also defines the possible transitions
// between states. States:
// 1. Compliance check created
// 2. Policies applied
// 3. Compliance proof requested
// 4.1. (actual 4.) Compliance proof generation failed (END - last state!)
// 4.2. (actual 5.) Compliance prrof generation succeeded
// 5. (actual 6.) Compliance proof attached to the selected settlement asset
// 6. (actual 7.) Settlement asset transferred to the beneficiary bank
// 7. (actual 8.) Assets released to the client
type ComplianceCheckStateManager struct{}

func CreateComplianceCheckStateManager() *ComplianceCheckStateManager {
	return &ComplianceCheckStateManager{}
}

func (*ComplianceCheckStateManager) UpdateComplianceCheckPolicyStatus(db *DB.DBHandler, complianceCheckID string, policyID int, isFailed bool, description string) error {
	returnErr := errors.New("failed to update compliance check policy status")
	if isFailed {
		db.UpdateTransactionPolicyStatus(complianceCheckID, policyID, 2, description)
	} else {
		db.UpdateTransactionPolicyStatus(complianceCheckID, policyID, 1, description)
	}

	statuses, err := db.GetTransactionPolicyStatuses(complianceCheckID)
	if err != nil {
		errlog.Println(err)
		return returnErr
	}

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

	return nil
}
