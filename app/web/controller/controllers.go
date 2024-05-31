package controller

import (
	"bisgo/app/web/core"
)

type HomeController struct {
	*core.Core
}

type TransactionController struct {
	*core.Core
}

type PolicyController struct {
	*core.Core
}

type APIController struct {
	*core.Core
}

type CBController struct {
	*core.Core
}

func CreateHomeController(core *core.Core) *HomeController {
	return &HomeController{core}
}

func CreateTransactionController(core *core.Core) *TransactionController {
	return &TransactionController{core}
}

func CreatePolicyController(core *core.Core) *PolicyController {
	return &PolicyController{core}
}

func CreateAPIController(core *core.Core) *APIController {
	return &APIController{core}
}

func CreateCBController(core *core.Core) *CBController {
	return &CBController{core}
}
