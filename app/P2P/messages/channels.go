package messages

import (
	"errors"
	"sync"
)

// request-response channel map
var rrChanMap sync.Map

func init() {
	// configuration settings
}

func StoreChannel(messageID int, channel chan<- any) {
	rrChanMap.Store(messageID, channel)
}

func LoadChannel(messageID int) (chan<- any, error) {
	channel, ok := rrChanMap.Load(messageID)

	if !ok {
		return nil, errors.New("error: missing peer for the message ID")
	}

	return channel.(chan<- any), nil
}

func RemoveChannel(messageID int) {
	rrChanMap.Delete(messageID)
}
