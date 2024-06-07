package messages

import (
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

func LoadChannel(messageID int) (chan<- any, bool) {
	channel, ok := rrChanMap.Load(messageID)

	if !ok {
		return nil, false
	}

	return channel.(chan<- any), true
}

func RemoveChannel(messageID int) {
	rrChanMap.Delete(messageID)
}
