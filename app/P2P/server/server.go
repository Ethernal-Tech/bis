package server

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/handler"
	"bisgo/app/P2P/messages"
	"bisgo/errlog"
	"encoding/json"
	"errors"
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

		// returning a 200 OK response regardless of message correctness, handling of invalid messages is internal
		w.WriteHeader(http.StatusOK)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			errlog.Println(err)
			return
		}

		var message messages.P2PServerMessasge
		err = json.Unmarshal(body, &message)
		if err != nil {
			errlog.Println(err)
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
			go s.CheckConfirmed(message.MessageID, message.Payload)
		case "cfm-result-beneficiary":
			go s.CFMResultBeneficiary(message.MessageID, message.Payload)
		case "cfm-result-originator":
			go s.CFMResultOriginator(message.MessageID, message.Payload)
		default:
			errlog.Println(errors.New("invalid p2p method received"))
		}
	})
}
