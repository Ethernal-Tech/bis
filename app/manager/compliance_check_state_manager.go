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

	// map containing all registered descriptions for all "still active" compliance checks,
	// "still active" compliance checks are those that are not in state 4 or state 8
	descriptions registeredDescriptions

	// synchronization primitive used for data-race-free access to the descriptions map
	mutDesc sync.Mutex
}

// registeredDescriptions type can be used to register descriptions for the compliance check states, it
// consists of two nested maps with the following meanings:
//  1. map[string]map (key - compliance check id, value - map 2)
//  2. map[ComplianceCheckState]string (key - concrete state, value - text details)
type registeredDescriptions map[string]map[ComplianceCheckState]string

var complianceCheckStateManager ComplianceCheckStateManager

func init() {
	complianceCheckStateManager = ComplianceCheckStateManager{
		db:           DB.GetDBHandler(),
		descriptions: registeredDescriptions{},
		mutDesc:      sync.Mutex{},
	}
}

// GetComplianceCheckStateManager returns the compliance check state manager ([ComplianceCheckStateManager]).
// Since only one manager exists in the entire system, each call returns the same instance. For more information
// (especially related to race conditions) read the [ComplianceCheckStateManager] documentation.
func GetComplianceCheckStateManager() *ComplianceCheckStateManager {
	return &complianceCheckStateManager
}

// RegisterDescription registers the description for a given compliance check and passed state. The registered
// description will be used (and recorded in the database) during the transition of the compliance check to the
// given state (using [*ComplianceCheckStateManager.Transition]). If the compliance check is in state 4 or 8,
// registration is rejected and false is returned. In the remaining states, registration will succeed and method
// returns true. However, be careful that registering a description for a state that the compliance check has
// already reached will have no effect. If the passed compliance check (complianceCheckId) does not exist, an
// error will be returned. The method is data-race-free. It can be used concurrently by multiple goroutines.
// However, a race condition may occur if the method is used for the same compliance check without synchronization.
// This is not the case for any two different compliance checks. In that case, everything is race-free.
func (m *ComplianceCheckStateManager) RegisterDescription(complianceCheckId string, state ComplianceCheckState, description string) (bool, error) {
	returnErr := errors.New("failed to register compliance check description")

	// although the compliance check itself is not required and used in method, the given invocation
	// is necessary to check whether the compliance check exists in the system
	_, err := m.db.GetComplianceCheckById(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return false, returnErr
	}

	states, err := m.db.GetAllComplianceCheckStates(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return false, returnErr
	}

	// registration can't be performed for compliance checks that are in state (4) or (8)
	for _, state := range states {
		if state.StateId == 4 || state.StateId == 8 {
			return false, nil
		}
	}

	m.mutDesc.Lock()

	// if there were no previous registrations for the given compliance check, it is necessary
	// to create a new map
	_, ok := m.descriptions[complianceCheckId]
	if !ok {
		m.descriptions[complianceCheckId] = map[ComplianceCheckState]string{}
	}

	m.descriptions[complianceCheckId][state] = description

	m.mutDesc.Unlock()

	return true, nil
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
