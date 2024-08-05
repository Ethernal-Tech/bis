package manager

import (
	"bisgo/app/DB"
	"bisgo/errlog"
	"errors"
	"sync"
)

// ComplianceCheckState indicates all possible states for the compliance check.
type ComplianceCheckState int

// concrete states
const (
	ComplianceCheckCreated ComplianceCheckState = iota
	PoliciesApplied
	ComplianceProofRequested
	ComplianceProofGenerationFailed
	ComplianceProofGenerationSucceeded
	ComplianceProofAttached
	SettlementAssetTransferred
	AssetsReleased
)

// ComplianceCheckStateManager manages the states of a compliance check during its lifecycle. All possible
// states of the compliance check are shown below. The order in which they are displayed also defines the
// possible transitions between states. States:
//  1. Compliance check created
//  2. Policies applied
//  3. Compliance proof requested
//  4. Compliance proof generation failed (END - last state!)
//  5. Compliance proof generation succeeded
//  6. Compliance proof attached to the selected settlement asset
//  7. Settlement asset transferred to the beneficiary bank
//  8. Assets released to the client
//
// From state 3, state manager can go to either state 4 or state 5. State 4 is entered if the compliance
// check fails and it denotes the last, final state for the given compliance checks. Otherwise, state
// manager moves from state 3 to state 5, and so on to all listed states (if possible).
//
// The manager is completely safe to use concurrently (simultaneously) for different compliance checks.
// It should not be used in concurrently manner for a single compliance check as it may cause undesired
// race conditions. However, it is data-race-free in both cases.
type ComplianceCheckStateManager struct {
	// database handler with all SQL wrapper methods
	db *DB.DBHandler

	// map containing all registered details for all "still active" compliance checks,
	// "still active" compliance checks are those that are not in state 4 or state 8
	details registeredDetails

	// synchronization primitive used for data-race-free access to the details map
	mutDet sync.Mutex
}

// registeredDetails type can be used to register details for the compliance check states, it consists of
// two nested maps with the following meanings:
//  1. map[string]map (key - compliance check id, value - map 2)
//  2. map[ComplianceCheckState]string (key - concrete state, value - text details)
type registeredDetails map[string]map[ComplianceCheckState]string

var complianceCheckStateManager ComplianceCheckStateManager

func init() {
	complianceCheckStateManager = ComplianceCheckStateManager{
		db:      DB.GetDBHandler(),
		details: registeredDetails{},
		mutDet:  sync.Mutex{},
	}
}

// GetComplianceCheckStateManager returns the compliance check state manager ([ComplianceCheckStateManager]).
// Since only one manager exists in the entire system, each call returns the same instance. For more information
// (especially related to race conditions) read the [ComplianceCheckStateManager] documentation.
func GetComplianceCheckStateManager() *ComplianceCheckStateManager {
	return &complianceCheckStateManager
}

func (m *ComplianceCheckStateManager) RegisterDetails(complianceCheckId string, state ComplianceCheckState, details string) {

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
