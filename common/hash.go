package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func GenerateHash(structure interface{}) (string, error) {
	// Get the current time
	currentTime := time.Now().UnixNano()

	// Generate a random number
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Float64()

	// Convert the structure to a JSON string
	structureBytes, err := json.Marshal(structure)
	if err != nil {
		return "", err
	}
	structureStr := string(structureBytes)

	// Concatenate the structure string, current time, and random number
	inputStr := fmt.Sprintf("%s%d%f", structureStr, currentTime, randomNumber)

	// Create a SHA-256 hash object
	hash := sha256.New()

	// Update the hash object with the input string
	hash.Write([]byte(inputStr))

	// Get the hexadecimal representation of the hash
	hashHex := hex.EncodeToString(hash.Sum(nil))

	return hashHex, nil
}
