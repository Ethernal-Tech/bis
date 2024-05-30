package handler

import "bisgo/app/P2P/core"

type P2PHandler struct {
	*core.Core
}

func CreateP2PHandler(core *core.Core) *P2PHandler {
	return &P2PHandler{core}
}
