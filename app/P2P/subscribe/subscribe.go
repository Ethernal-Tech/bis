package subscribe

import (
	"bisgo/errlog"
	"errors"
	"sync"
)

var subscriptions map[MessageType][]chan<- any
var mutSub sync.Mutex

func init() {
	subscriptions = make(map[MessageType][]chan<- any)
}

func Subscribe(messageType MessageType, hitFunc func([]any) (any, bool)) (<-chan any, int, error) {
	subCh, subID, err := subscribe(messageType)
	if err != nil {
		errlog.Println(err)
		return nil, -1, err
	}

	message, exists, err := FindMessage(messageType, hitFunc)
	if err != nil {
		errlog.Println(err)
		return nil, -1, err
	}

	if exists {
		subCh <- message
	}

	return subCh, subID, nil
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

func Unsubscribe(messageType MessageType, subID int) error {
	err := unsubscribe(messageType, subID)
	if err != nil {
		errlog.Println(err)
		return err
	}

	return nil
}

func unsubscribe(messageType MessageType, subID int) error {
	if !messageTypeExists(messageType) {
		return errors.New("unknown message type")
	}

	defer mutSub.Unlock()
	mutSub.Lock()

	subscriptions[messageType] = append(subscriptions[messageType][:subID], subscriptions[messageType][subID+1:]...)

	return nil
}

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
