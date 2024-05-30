package handler

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/messages"
)

type P2PHandler struct {
	*core.Core
}

func CreateP2PHandler(core *core.Core) *P2PHandler {
	return &P2PHandler{core}
}

func (h *P2PHandler) Method1(messageID int, payload []byte) {
	channel, err := messages.LoadChannel(messageID)

	if err != nil {
		// handle error
	}

	// handler logic

	channel <- nil // send data to the listener

	messages.RemoveChannel(messageID) // remove channel from the map

}

func (h *P2PHandler) Method2(messageID int, payload []byte) {
	// handler logic
}
