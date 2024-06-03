package handler

import (
	"bisgo/app/proving/core"
)

type ProvingHandler struct {
	*core.Core
}

func CreateProvingHandler(core *core.Core) *ProvingHandler {
	return &ProvingHandler{core}
}

func (h *ProvingHandler) HandleInteractiveProof() {

}

func (h *ProvingHandler) HandleNonInteractiveProof() {

}
