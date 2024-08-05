package subscribe

import (
	"bisgo/errlog"
	"errors"
	"math/rand"
	"sync"
)

// Subscription contains all the necessary parameters to define a single subscription
type Subscription struct {
	// unique subscription identifier
	id int

	// channel for sending a copy of a newly arrived message
	NotifyCh chan any

	// internal unsubscribe signal channel
	unsubCh chan struct{}
}

type subscriptionMap map[int]Subscription

// map containing all subscriptions classified by the type of message to which the subscription applies
var subscriptions map[MessageType]subscriptionMap

// each message type has its own mutex
var mutSub map[MessageType]*sync.Mutex

func init() {
	subscriptions = make(map[MessageType]subscriptionMap)
	subscriptions[SCLServerStarted] = make(subscriptionMap)

	mutSub = make(map[MessageType]*sync.Mutex)
	mutSub[SCLServerStarted] = &sync.Mutex{}
}

// Subscribe function allows you to subscribe to messages of a given type. Since the looking message may have
// already be in the system, matchFunc parameter defines a function that will perform the check (read more about
// matchFunc in the [FindMessage] documentation). Return value of the Subscribe function is the subscription itself.
// It contains a channel (NotifyCh) to which messages will be sent upon arrival in the system. If the message is
// already in the system (checked via matchFunc), it will be in the channel buffer before the return.
func Subscribe(messageType MessageType, matchFunc func([]any) (any, bool)) (Subscription, error) {
	subscription, err := subscribe(messageType)
	if err != nil {
		errlog.Println(err)
		return Subscription{}, err
	}

	message, exists, err := FindMessage(messageType, matchFunc)
	if err != nil {
		errlog.Println(err)
		return Subscription{}, err
	}

	if exists {
		subscription.NotifyCh <- message
	}

	return subscription, nil
}

func subscribe(messageType MessageType) (Subscription, error) {
	if !messageTypeExists(messageType) {
		return Subscription{}, errors.New("unknown message type")
	}

	subscription := Subscription{
		id:       rand.Int(),
		NotifyCh: make(chan any, 1),
		unsubCh:  make(chan struct{}),
	}

	defer mutSub[messageType].Unlock()
	mutSub[messageType].Lock()

	subscriptions[messageType][subscription.id] = subscription

	return subscription, nil
}

// Unsubscribe is used to remove a subscription from the subscription list for a given message type.
func Unsubscribe(messageType MessageType, subscription Subscription) error {
	err := unsubscribe(messageType, subscription)
	if err != nil {
		errlog.Println(err)
		return err
	}

	return nil
}

func unsubscribe(messageType MessageType, subscription Subscription) error {
	if !messageTypeExists(messageType) {
		return errors.New("unknown message type")
	}

	close(subscription.unsubCh)

	defer mutSub[messageType].Unlock()
	mutSub[messageType].Lock()

	delete(subscriptions[messageType], subscription.id)

	return nil
}

// StoreAndNotify stores the message and notifies (send a message copy to) all subscribers to the given message type.
func StoreAndNotify(messageType MessageType, message any) error {
	err := StoreMessage(messageType, message)
	if err != nil {
		errlog.Println(err)
		return err
	}

	defer mutSub[messageType].Unlock()
	mutSub[messageType].Lock()
	for _, subscription := range subscriptions[messageType] {
		select {
		case subscription.NotifyCh <- message:
		case <-subscription.unsubCh:
		}
	}

	return nil
}
