package manager

import (
	"bisgo/common"
	"bisgo/config"
	"bisgo/errlog"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
)

var peers sync.Map

func init() {
	queriedPeers, err := getAvailablePeers()
	if err != nil {
		errlog.Println(err)
	}

	for k, v := range queriedPeers {
		peers.Store(k, v)
	}
}

func GetBankPeerID(bankID string) (string, error) {
	peerID, ok := peers.Load(bankID)
	if !ok {
		queriedPeers, err := getAvailablePeers()
		if err != nil {
			errlog.Println(err)
			return "", errors.New("cannot get available peers")
		}

		peerID := ""
		for k, v := range queriedPeers {
			if k == bankID {
				peerID = v
			}
			peers.Store(k, v)
		}

		if peerID != "" {
			return peerID, nil
		}

		return "", errors.New("cannot find requested peer in a set of available peers")
	}

	return peerID.(string), nil
}

func getAvailablePeers() (map[string]string, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", strings.Replace(config.ResolveP2PNodeAPIAddress(), "passthrough", "peers", 1), nil)
	if err != nil {
		errlog.Println(err)
		return nil, errors.New("cannot create a p2p wrapping (HTTP POST) message due to internal system problems; not caused by incorrect parameters")
	}

	response, err := client.Do(request)
	if err != nil {
		errlog.Println(err)
		return nil, errors.New("cannot send a p2p wrapping (HTTP POST) message due to internal system problems; not caused by incorrect parameters")
	}

	if response.StatusCode != 200 {
		returnErr := errors.New("p2p node rejected the message")
		body, err := io.ReadAll(response.Body)
		if err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		var response common.ErrorResponse
		if err := json.Unmarshal(body, &response); err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		errlog.Println(errors.New("p2p node response: " + response.Message))

		return nil, returnErr
	}

	returnErr := errors.New("p2p node sent malformed response")

	body, err := io.ReadAll(response.Body)
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	var result map[string]string = map[string]string{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	return result, nil
}
