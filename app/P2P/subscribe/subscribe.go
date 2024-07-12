package subscribe

import (
	"bisgo/errlog"
	"errors"
	"sync"
)

// map containing all subscriptions classified by the type of message to which the subscription applies
var subscriptions map[MessageType][]chan<- any
var mutSub sync.Mutex

func init() {
	subscriptions = make(map[MessageType][]chan<- any)
}

// Subscribe function allows you to subscribe to messages of a given type. Since the looking message may have already be in the
// system, matchFunc parameter defines a function that will perform the check (read more about matchFunc in the [FindMessage]
// documentation). Return values ​​of the Subscribe function are the channel (to which a copy of the message will be sent upon
// arrival) and subId (used during unsubscribe). If the message is already in the system (checked via matchFunc), it will be
// in the channel buffer before the channel is returned.
func Subscribe(messageType MessageType, matchFunc func([]any) (any, bool)) (<-chan any, int, error) {
	subCh, subId, err := subscribe(messageType)
	if err != nil {
		errlog.Println(err)
		return nil, -1, err
	}

	message, exists, err := FindMessage(messageType, matchFunc)
	if err != nil {
		errlog.Println(err)
		return nil, -1, err
	}

	if exists {
		subCh <- message
	}

	return subCh, subId, nil
}

func subscribe(messageType MessageType) (chan any, int, error) {
	if !messageTypeExists(messageType) {
		return nil, -1, errors.New("unknown message type")
	}

	subCh := make(chan any, 1)

	defer mutSub.Unlock()
	mutSub.Lock()

	subscriptions[messageType] = append(subscriptions[messageType], subCh)

	return subCh, len(subscriptions[messageType]) - 1, nil
}

// Unsubscribe is used to remove subsribe (subId) from the subscription storage for a given message type.
func Unsubscribe(messageType MessageType, subId int) error {
	err := unsubscribe(messageType, subId)
	if err != nil {
		errlog.Println(err)
		return err
	}

	return nil
}

func unsubscribe(messageType MessageType, subId int) error {
	if !messageTypeExists(messageType) {
		return errors.New("unknown message type")
	}

	defer mutSub.Unlock()
	mutSub.Lock()

	subscriptions[messageType] = append(subscriptions[messageType][:subId], subscriptions[messageType][subId+1:]...)

	return nil
}

// StoreAndNotify stores the message and notifies (send a copy of the message) to all subscribers to the given message type.
func StoreAndNotify(messageType MessageType, message any) error {
	err := StoreMessage(messageType, message)
	if err != nil {
		errlog.Println(err)
		return err
	}

	defer mutSub.Unlock()
	mutSub.Lock()
	for _, ch := range subscriptions[messageType] {
		ch <- message
	}

	return nil
}
