package manager

import (
	"errors"
	"sync"
)

var peers sync.Map

func init() {
	var mocks map[string]string = map[string]string{
		// JPM Chase
		"984500653R409CC5AB28": "hash1",
		// MAS
		"54930035WQZLGC45RZ35": "hash2",
		// HLB
		"549300BUPYUQGB5BFX94": "hash3",
		// BNM
		"549300NROGNBV2T1GS07": "hash4",
	}

	for k, v := range mocks {
		peers.Store(k, v)
	}
}

func GetBankPeerID(bankID string) (string, error) {
	peerID, ok := peers.Load(bankID)
	if !ok {
		return "", errors.New("error: missing peer for the selected bank")
	}

	return peerID.(string), nil
}
