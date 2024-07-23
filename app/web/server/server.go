package server

import (
	"bisgo/app/web/controller"
	"bisgo/app/web/core"
	"net/http"
)

// WebServer represents a wrapper around the standard library's http.Server type.
// It integrates custom handlers grouped into controllers and embeds a wide range of additional core functionalities.
type WebServer struct {
	*controller.HomeController
	*controller.ComplianceCheckController
	*controller.TransactionController
	*controller.PolicyController
	*controller.CBController
	*controller.APIController
	*core.Core
	*http.Server
}

var server WebServer

func init() {
	var core *core.Core = core.CreateCore()

	server = WebServer{
		HomeController:            controller.CreateHomeController(core),
		ComplianceCheckController: controller.CreateComplianceCheckController(core),
		TransactionController:     controller.CreateTransactionController(core),
		PolicyController:          controller.CreatePolicyController(core),
		CBController:              controller.CreateCBController(core),
		APIController:             controller.CreateAPIController(core),
		Core:                      core,
	}
}

func GetWebServer() *WebServer {
	return &server
}
