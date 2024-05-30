package client

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/handler"
)

type P2PServer struct {
	*handler.P2PHandler
	*core.Core
}

var server P2PServer

func init() {
	var core *core.Core = core.CreateCore()

	server = P2PServer{
		P2PHandler: handler.CreateP2PHandler(core),
		Core:       core,
	}
}

func GetP2PServer() *P2PServer {
	return &server
}

func (s *P2PServer) Mux() {

}
