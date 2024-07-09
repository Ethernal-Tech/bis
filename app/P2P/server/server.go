package server

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/handler"
	"bisgo/app/P2P/messages"
	"bisgo/errlog"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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

		// TODO: describe why it is needed
		payload := make([]byte, len(message.Payload))
		for i, v := range message.Payload {
			payload[i] = byte(v)
		}

		go func() {
			var err error

			messageTypeLog(r, message.Method)

			switch message.Method {
			case "create-transaction":
				err = s.CreateTransaction(message.MessageID, payload)
			case "get-policies":
				err = s.GetPolicies(message.MessageID, payload)
			case "policies":
				err = s.ReceivePolicies(message.MessageID, payload)
			case "check-confirmed":
				err = s.CheckConfirmed(message.MessageID, payload)
			case "cfm-result-beneficiary":
				err = s.CFMResultBeneficiary(message.MessageID, payload)
			case "cfm-result-originator":
				err = s.CFMResultOriginator(message.MessageID, payload)
			default:
				err = errors.New("invalid p2p method received")
			}

			if err != nil {
				errlog.Println(err)
			}
		}()
	})
}

func messageTypeLog(r *http.Request, method string) {
	dateTime := strings.Split(time.Now().String(), ".")[0]
	sender := r.RemoteAddr

	fmt.Printf("\033[34m%v INFO\033[0m: [%v] %v\n", dateTime, sender, method)
}
