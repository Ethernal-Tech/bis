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
	// Messages of this type contain the compliance check id, so it is easy to identify
	// for which compliance check the signal has arrived.
	SCLServerStarted MessageType = iota
)

// map that contains "all" received messages classified by its type,
// "all" in the context of possible message types that can be subscibed to
var messages map[MessageType][]any
var mutMess sync.Mutex

func init() {
	messages = make(map[MessageType][]any)
}

// messageTypeExists, as its name suggests, returns whether the specified message type exists or not
func messageTypeExists(messageType MessageType) bool {
	if messageType != SCLServerStarted {
		return false
	}

	return true
}

// StoreMessage stores a message in the storage for the given message type.
func StoreMessage(messageType MessageType, message any) error {
	if !messageTypeExists(messageType) {
		return errors.New("unknown message type")
	}

	defer mutMess.Unlock()
	mutMess.Lock()

	messages[messageType] = append(messages[messageType], message)

	return nil
}

// FindMessage, as the name suggests, can be used to find a specific message of a given type.
// Matching logic is defined through the matchFunc parameter. Input to the matchFunc is a list
// of all messages of the given type that are known in the system, while the return values are the
// concrete message and a flag indicating whether it was found or not. It could be said that the
// FindMessage function is a wrapper around the matchFunc, thus the return values ​​(excluding error)
// of the FindMessage function are exactly the return values ​​of the matchFunc.
func FindMessage(messageType MessageType, matchFunc func([]any) (any, bool)) (any, bool, error) {
	if !messageTypeExists(messageType) {
		return nil, false, errors.New("unknown message type")
	}

	defer mutMess.Unlock()
	mutMess.Lock()

	if len(messages[messageType]) == 0 {
		return nil, false, nil
	}

	message, exist := matchFunc(messages[messageType])

	return message, exist, nil
}

// LoadMessages loads (returns) all messages of the given type
func LoadMessages(messageType MessageType) ([]any, bool, error) {
	if !messageTypeExists(messageType) {
		return nil, false, errors.New("unknown message type")
	}

	defer mutMess.Unlock()
	mutMess.Lock()

	if len(messages[messageType]) == 0 {
		return nil, false, nil
	}

	// because of the way slices are implemented in golang, copy of the slice must be returned
	// otherwise we would have a race condition (data race)
	retMessages := make([]any, len(messages[messageType]))
	copy(retMessages, messages[messageType])

	return retMessages, true, nil
}
