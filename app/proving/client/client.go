package client

import (
	"bisgo/config"
)

type ProvingClient struct {
	gpjcApiAddress string
}

var client ProvingClient

func init() {
	client = ProvingClient{config.ResolveGpjcApiAddress()}
}

func GetProvingClient() *ProvingClient {
	return &client
}

func (c *ProvingClient) SendProofRequest(proofType string /* additional parameters */) error {

	// logic for sending proof requests
	// different logic based on whether the request is for interactive or non-interactive proof

	return nil
}
