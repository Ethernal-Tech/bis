package core

import (
	"bisgo/app/DB"
	p2pclient "bisgo/app/P2P/client"
	provingclient "bisgo/app/proving/client"
)

type Core struct {
	DB            *DB.DBHandler
	P2PClient     *p2pclient.P2PClient
	ProvingClient *provingclient.ProvingClient
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore() *Core {
	return &Core{
		DB:            DB.GetDBHandler(),
		P2PClient:     p2pclient.GetP2PClient(),
		ProvingClient: provingclient.GetProvingClient(),
	}
}
