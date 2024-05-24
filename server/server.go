package server

import (
	"bisgo/config"
	"bisgo/controller"
	"bisgo/core"
	"fmt"
	"log"
	"net/http"
)

// Server represents a wrapper around the standard library's http.Server type.
// It integrates custom handlers grouped into controllers and embeds a wide range of additional core functionalities.
type Server struct {
	*controller.HomeController
	*controller.TransactionController
	*controller.PolicyController
	*controller.CBController
	*controller.APIController
	*core.Core
	*http.Server
}

func Run() {
	// configuration settings
	var config = config.CreateConfig()

	var core *core.Core = core.CreateCore(config)
	defer core.DB.Close()

	server := Server{
		HomeController:        controller.CreateHomeController(core),
		TransactionController: controller.CreateTransactionController(core),
		PolicyController:      controller.CreatePolicyController(core),
		CBController:          controller.CreateCBController(core),
		APIController:         controller.CreateAPIController(core),
		Core:                  core,
	}

	server.Server = &http.Server{
		Addr:    config.ResolveServerPort(),
		Handler: server.routes(),
	}

	fmt.Println("The server starts at", server.Addr)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("\033[31m" + "Failed to start the server!" + "\033[31m")
	}
}
