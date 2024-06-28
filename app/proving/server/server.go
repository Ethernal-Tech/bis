package server

import (
	"bisgo/app/proving/core"
	"bisgo/app/proving/handler"
	"bisgo/errlog"
	"io"
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errlog.Println(err)
			return
		}

		switch {
		case strings.Contains(r.URL.Path, "/proof/interactive"):
			s.HandleInteractiveProof(body)
		case strings.Contains(r.URL.Path, "/proof/noninteractive"):
			s.HandleNonInteractiveProof()
		default:
			http.Error(w, "Invalid proof type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
