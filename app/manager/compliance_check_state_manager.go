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
//  4. Compliance proof generation failed (END)
//  5. Compliance proof generation succeeded
//  6. Compliance proof attached to the selected settlement asset
//  7. Settlement asset transferred to the beneficiary bank
//  8. Assets released to the client (END)
//
// Final states are 4 (Compliance proof generation failed) and 8 (Assets released to the client).
//
// From state 3, state manager can go to either state 4 or state 5. State 4 is entered if the compliance
// check fails and it denotes the last, final state for the given compliance checks. Otherwise, state
// manager moves from state 3 to state 5, and so on through all listed states to the final state 8.
//
// The manager is completely safe to use concurrently (simultaneously) for different, but also, for the
// same compliance check.
type ComplianceCheckStateManager struct {
	// database handler with all SQL wrapper methods
	db *DB.DBHandler

	// map containing all registered descriptions for all "still active" compliance checks,
	// "still active" compliance checks are those that are not in state 4 or state 8
	descriptions registeredDescriptions

	// synchronization primitive used for data-race-free access to the descriptions map
	mutDesc sync.Mutex

	// channel to which requests for mutexes are sent
	mutReqCh chan mutReq

	// channel to which the signal is sent that the conformance check has reached the final
	// state and that the corresponding mutex can be deleted
	finSigCh chan string
}

// registeredDescriptions type can be used to register descriptions for the compliance check states, it
// consists of two nested maps with the following meanings:
//  1. map[string]map (key - compliance check id, value - map 2)
//  2. map[ComplianceCheckState]string (key - concrete state, value - text details)
type registeredDescriptions map[string]map[ComplianceCheckState]string

// mutReq structure type is used to create a request to obtain a mutex for a given compliance check
type mutReq struct {
	// unique identifier of the compliance check for which mutex is requested
	id string

	// channel to which the corresponding mutex should be sent
	resCh chan<- *sync.Mutex
}

var complianceCheckStateManager ComplianceCheckStateManager

func init() {
	complianceCheckStateManager = ComplianceCheckStateManager{
		db:           DB.GetDBHandler(),
		descriptions: registeredDescriptions{},
		mutDesc:      sync.Mutex{},
		mutReqCh:     make(chan mutReq),
		finSigCh:     make(chan string),
	}

	go func() {
		var mutCCs = map[string]*sync.Mutex{}

	inf_loop:
		for {
			select {
			case request := <-complianceCheckStateManager.mutReqCh:
				select {
				case signal := <-complianceCheckStateManager.finSigCh:
					if request.id == signal {
						delete(mutCCs, signal)
						request.resCh <- nil
						continue inf_loop
					}
				default:
				}

				if mut, ok := mutCCs[request.id]; ok {
					request.resCh <- mut
					continue
				}

				states, err := complianceCheckStateManager.db.GetAllComplianceCheckStates(request.id)
				if err != nil {
					errlog.Println(err)
					request.resCh <- nil
					continue
				}

				for _, state := range states {
					if state.StateId == 4 || state.StateId == 8 {
						request.resCh <- nil
						continue inf_loop
					}
				}

				var newMut *sync.Mutex
				mutCCs[request.id] = newMut

				request.resCh <- newMut
			case signal := <-complianceCheckStateManager.finSigCh:
				delete(mutCCs, signal)
			}
		}
	}()
}

// GetComplianceCheckStateManager returns the compliance check state manager ([ComplianceCheckStateManager]).
// Since only one manager exists in the entire system, each call returns the same instance. For more information
// (especially related to race conditions) read the [ComplianceCheckStateManager] documentation.
func GetComplianceCheckStateManager() *ComplianceCheckStateManager {
	return &complianceCheckStateManager
}

// RegisterDescription registers the description for a given compliance check and passed state. The registered
// description will be used (and recorded in the database) during the transition of the compliance check to the
// given state (using [*ComplianceCheckStateManager.Transition]). A list of available states can be found in the
// documentation for [ComplianceCheckStateManager]. If the compliance check is in state 4 or 8, registration is
// rejected and false is returned. In the remaining states, registration will succeed and method returns true.
// However, be careful that registering a description for a state that the compliance check has already reached
// will have no effect. If the passed compliance check (complianceCheckId) or state does not exist, an error will
// be returned. The method is data-race-free. It can be used concurrently by multiple goroutines. However, some
// race conditions may occur if the method is used for the same compliance check without synchronization. It can
// be especially noticeable if the method is used without synchronization (concurrently) with [Transition]. This
// is not the case for any two different compliance checks. In that case, everything is race-free.
func (m *ComplianceCheckStateManager) RegisterDescription(complianceCheckId string, state ComplianceCheckState, description string) (bool, error) {
	returnErr := errors.New("failed to register compliance check description")

	// input state must be one defined by the [ComplianceCheckStateManager] doc
	if state < ComplianceCheckCreated || state > AssetsReleased {
		errlog.Println(errors.New("unknown compliance check state"))
		return false, returnErr
	}

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

// Transition, as the name suggests, performs the transition, if possible, of a given compliance check from the
// previous state to the next state. The method behaves like a black box. The complete logic of transitions and
// state machine is internally implemented. The only thing required is to call the method in the "right" places.
// If the transition is possible, the new state is returned and the flag (second return value) is set to true.
// Otherwise, it returns false and the previous (current) state. If there is an error, a special [ErrState] is
// returned instead of the state. The method is fully free of race conditions. In other words, it can be used
// concurrently both for different compliance checks and for the same one.
func (m *ComplianceCheckStateManager) Transition(complianceCheckId string) (ComplianceCheckState, bool, error) {
	returnErr := errors.New("unsuccessful transition of the compliance check state")

	// although the compliance check itself is not required and used in method, the given invocation
	// is necessary to check whether the compliance check exists in the system
	_, err := m.db.GetComplianceCheckById(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return ErrState, false, returnErr
	}

	var ch = make(chan *sync.Mutex)

	request := mutReq{
		id:    complianceCheckId,
		resCh: ch,
	}

	m.mutReqCh <- request

	mut := <-ch

	// if no mutex is obtained, it means that the compliance check is already in one of the final states
	// thus state machine will get into final state cases and no race condition can occur
	if mut != nil {
		mut.Lock()
		defer mut.Unlock()
	}

	states, err := m.db.GetAllComplianceCheckStates(complianceCheckId)
	if err != nil {
		errlog.Println(err)
		return ErrState, false, returnErr
	}

	// taking the current state of the compliance check
	currentState := zeroState
	for _, state := range states {
		if currentState < ComplianceCheckState(state.StateId) {
			currentState = ComplianceCheckState(state.StateId)
		}
	}

	// flag indicating whether it is necessary to delete the description for a given compliance check,
	// it is activated only if a transition is made to final states (4 or 8)
	var deleteDescriptions bool

	// state machine
	switch currentState {
	case zeroState, ComplianceCheckCreated, PoliciesApplied:
		// for the first three initial states, the only thing that needs to be done
		// is the transition to the next state (no additional logic is present)
		{
			m.mutDesc.Lock()
			description := m.descriptions[complianceCheckId][currentState+1]
			m.mutDesc.Unlock()

			err := m.db.AddComplianceCheckState(complianceCheckId, int(currentState+1), description)
			if err != nil {
				errlog.Println(err)
				return ErrState, false, returnErr
			}
		}
	case ComplianceProofRequested:
		// this state represents the one in which the policies checks are executed for the
		// given compliance check; therefore, if the compliance check is in this state, it
		// is necessary to take all its policies and see their statuses; there are three
		// possible scenarios:
		// 1. some of policies are still being checked (policy status 0 - pending)
		//    - scenario 2 should be checked
		//    - nothing additional should be done
		// 2. some policy has failed
		//    - there is no need to further check other policies
		//    - transition to the failure state (state 4)
		//    - raise deleteDescriptions flag
		// 3. all policies have passed
		//    - transition to the succeed state (state 5)
		{
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

					m.finSigCh <- complianceCheckId

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
				return currentState, false, nil
			}
		}
	case ComplianceProofGenerationFailed, AssetsReleased:
		// final states, nothing needs to be done
		{
			return currentState, false, nil
		}
	case ComplianceProofGenerationSucceeded:
	case SettlementAssetTransferred:
		m.finSigCh <- complianceCheckId

		deleteDescriptions = true
	}

	if deleteDescriptions {
		m.mutDesc.Lock()
		delete(m.descriptions, complianceCheckId)
		m.mutDesc.Unlock()
	}

	return currentState + 1, true, nil
}
