package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetAvailablePeers(p2pNodeAddress string) ([]Peer, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", strings.Join([]string{p2pNodeAddress, "v1", "p2p", "peers"}, "/"), nil)
	if err != nil {
		return []Peer{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []Peer{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Peer{}, err
	}

	if resp.StatusCode != 200 {
		return []Peer{}, handleError(body)
	}

	var ret []Peer
	err = json.Unmarshal(body, &ret)

	return ret, err
}

func SendPassthruMessage(p2pNodeAddress string, requestData []byte) error {
	client := &http.Client{}

	req, err := http.NewRequest("POST", strings.Join([]string{"http:/", p2pNodeAddress, "v1", "p2p", "passthru"}, "/"), bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return handleError(body)
	}

	return nil
}

func handleError(responseBody []byte) error {
	var resp ErrorResponse
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return err
	}

	return fmt.Errorf("call to p2p failed with error: " + resp.Message)
}
