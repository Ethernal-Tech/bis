package common

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

func KeyGen() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)

	if err != nil {
		panic(err)
	}

	return key
}

func Sign(message string, key *ecdsa.PrivateKey) []byte {
	// Turn the message into a 32-byte hash
	hash := solsha3.SoliditySHA3(solsha3.String(message))
	sig, err := crypto.Sign(hash, key)

	if err != nil {
		panic(err)
	}

	signatureLastByte := sig[len(sig)-1]
	if signatureLastByte == 0 || signatureLastByte == 1 {
		sig[len(sig)-1] += 27
	}

	return sig
}
