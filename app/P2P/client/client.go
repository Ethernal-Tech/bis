package client

type P2PClient struct {
	p2pNodeAddress string
}

var client P2PClient

func init() {
	// configuration settings (e.g. for p2pNodeAddress)
}

func GetP2PClient() *P2PClient {
	return &client
}
