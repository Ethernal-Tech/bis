package manager

import (
	"errors"
	"sync"
)

var peers sync.Map

func init() {
	// configuration settings
}

func GetBankPeerID(bankID string) (string, error) {
	peerID, ok := peers.Load(bankID)
	if !ok {
		return "", errors.New("error: missing peer for the selected bank")
	}

	return peerID.(string), nil
}
