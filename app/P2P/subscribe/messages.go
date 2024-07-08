package subscribe

import (
	"errors"
	"sync"
)

// MessageType indicates all message types that can be subscribed to in the system.
type MessageType int

// concrete message types
const (
	// SCLServerStarted represents a signaling message that the GPJC server is started.
	// Messages of this type contain the compliance check ID, so it is easy to identify
	// for which compliance check the signal has arrived.
	SCLServerStarted MessageType = iota
)

// Map that contains "all" received messages classified by its type.
// "all" in the context of possible message types that can be subscibed to.
var messages map[MessageType][]any
var mutMess sync.Mutex

// messageTypeExists, as its name suggests, returns whether the specified message type exists or not
func messageTypeExists(messageType MessageType) bool {
	if messageType != SCLServerStarted {
		return false
	}

	return true
}

// StoreMessage
func StoreMessage(messageType MessageType, message any) error {
	if !messageTypeExists(messageType) {
		return errors.New("unknown message type")
	}

	defer mutMess.Unlock()
	mutMess.Lock()

	messages[messageType] = append(messages[messageType], message)

	return nil
}

// FindMessage
func FindMessage(messageType MessageType, hitFunc func([]any) (any, bool)) (any, bool, error) {
	if !messageTypeExists(messageType) {
		return nil, false, errors.New("unknown message type")
	}

	defer mutMess.Unlock()
	mutMess.Lock()

	if len(messages[messageType]) == 0 {
		return nil, false, nil
	}

	message, exist := hitFunc(messages[messageType])

	return message, exist, nil
}

// LoadMessages
func LoadMessages(messageType MessageType) ([]any, bool, error) {
	if !messageTypeExists(messageType) {
		return nil, false, errors.New("unknown message type")
	}

	defer mutMess.Unlock()
	mutMess.Lock()

	if len(messages[messageType]) == 0 {
		return nil, false, nil
	}

	// Because of the way slices are implemented in golang, copy of the slice must be returned
	// Otherwise we would have a race condition (data race)
	retMessages := make([]any, len(messages[messageType]))
	copy(retMessages, messages[messageType])

	return retMessages, true, nil
}
