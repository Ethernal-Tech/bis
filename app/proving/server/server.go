package server

import (
	"bisgo/app/proving/core"
	"bisgo/app/proving/handler"
	"net/http"
	"strings"
)

type ProvingServer struct {
	*handler.ProvingHandler
	*core.Core
}

var server ProvingServer

func init() {
	var core *core.Core = core.CreateCore()

	server = ProvingServer{
		ProvingHandler: handler.CreateProvingHandler(core),
		Core:           core,
	}
}

func GetProvingServer() *ProvingServer {
	return &server
}

func (s *ProvingServer) Mux() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/proof/interactive"):
			s.HandleInteractiveProof()
		case strings.Contains(r.URL.Path, "/proof/noninteractive"):
			s.HandleNonInteractiveProof()
		default:
			http.Error(w, "Invalid proof type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
