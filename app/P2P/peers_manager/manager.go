package manager

import (
	"bisgo/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var peers sync.Map

func init() {
	queriedPeers, err := GetAvaiablePeers()
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range queriedPeers {
		peers.Store(k, v)
	}
}

func GetBankPeerID(bankID string) (string, error) {
	peerID, ok := peers.Load(bankID)
	if !ok {
		queriedPeers, err := GetAvaiablePeers()
		if err != nil {
			fmt.Println(err)
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

		return "", errors.New("error: missing peer for the selected bank")
	}

	return peerID.(string), nil
}

func GetAvaiablePeers() (map[string]string, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", strings.Replace(config.ResolveP2PNodeAPIAddress(), "passthrough", "peers", 1), nil)

	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
