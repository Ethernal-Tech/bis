package messages

import (
	"sync"
)

// request-response channel map
var rrChanMap sync.Map

func init() {
	// configuration settings
}

// TODO: add a function that will periodically (in time intervals) remove unused channels (start it as a goroutine from init)

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
