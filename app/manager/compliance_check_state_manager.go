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
	// ErrState represents the value that is returned when the transition didn't succeed,
	// it is not a real state (not part of the state machine) nor the real error
	ErrState ComplianceCheckState = iota - 1

	// internally used state (not a real state)
	zeroState

	ComplianceCheckCreated
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

// add return current state method

func (m *ComplianceCheckStateManager) Transition(complianceCheckId string) (ComplianceCheckState, bool, error) {
	returnErr := errors.New("unsuccessful transition of the compliance check state")

	// although the compliance check itself is not required and used in method, the given invocation
	// is necessary to check whether the compliance check exists in the system
	_, err := m.db.GetComplianceCheckById(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return ErrState, false, returnErr
	}

	states, err := m.db.GetAllComplianceCheckStates(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return ErrState, false, returnErr
	}

	currentState := zeroState
	for _, state := range states {
		if currentState < ComplianceCheckState(state.StateId) {
			currentState = ComplianceCheckState(state.StateId)
		}
	}

	var deleteDescriptions bool

	// state machine
	switch currentState {
	case zeroState, ComplianceCheckCreated, PoliciesApplied:
		m.mutDesc.Lock()
		description := m.descriptions[complianceCheckId][currentState+1]
		m.mutDesc.Unlock()

		err := m.db.AddComplianceCheckState(complianceCheckId, int(currentState+1), description)
		if err != nil {
			errlog.Println(err)
			return ErrState, false, returnErr
		}

	case ComplianceProofRequested:
		policies, err := m.db.GetPoliciesByComplianceCheckId(complianceCheckId)
		if err != nil {
			errlog.Println(err)
			return ErrState, false, returnErr
		}

		var allPassed = true

		for _, policy := range policies {
			status, err := m.db.GetPolicyStatus(complianceCheckId, policy.Policy.Id)
			if err != nil {
				errlog.Println(err)
				return ErrState, false, returnErr
			}

			if status == 0 {
				allPassed = false
			} else if status == 2 {
				deleteDescriptions = true

				m.mutDesc.Lock()
				description := m.descriptions[complianceCheckId][ComplianceProofGenerationFailed]
				m.mutDesc.Unlock()

				err := m.db.AddComplianceCheckState(complianceCheckId, int(ComplianceProofGenerationFailed), description)
				if err != nil {
					errlog.Println(err)
					return ErrState, false, returnErr
				}

				break
			}
		}

		if allPassed {
			m.mutDesc.Lock()
			description := m.descriptions[complianceCheckId][ComplianceProofGenerationSucceeded]
			m.mutDesc.Unlock()

			err := m.db.AddComplianceCheckState(complianceCheckId, int(ComplianceProofGenerationSucceeded), description)
			if err != nil {
				errlog.Println(err)
				return ErrState, false, returnErr
			}
		} else {
			return currentState, false, returnErr
		}
	case ComplianceProofGenerationFailed, AssetsReleased:
		return currentState, false, nil
	case ComplianceProofGenerationSucceeded:
	}

	if deleteDescriptions {
		m.mutDesc.Lock()
		delete(m.descriptions, complianceCheckId)
		m.mutDesc.Unlock()
	}

	return currentState + 1, true, nil
}
