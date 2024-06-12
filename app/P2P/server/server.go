package server

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/handler"
	"bisgo/app/P2P/messages"
	"encoding/json"
	"io"
	"net/http"
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

func (s *P2PServer) Mux() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}

		var message messages.P2PServerMessasge
		err = json.Unmarshal(body, &message)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		switch message.Method {
		case "create-transaction":
			go s.CreateTransaction(message.MessageID, message.Payload)
		case "get-policies":
			go s.GetPolicies(message.MessageID, message.Payload)
		case "send-policies":
			go s.SendPolicies(message.MessageID, message.Payload)
		case "check-confirmed":
			s.CheckConfirmed(message.MessageID, message.Payload)
		case "cfm-check-result":
			s.CheckConfirmed(message.MessageID, message.Payload)
		default:
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
