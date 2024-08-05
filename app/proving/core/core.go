package core

import (
	"bisgo/app/DB"
	p2pclient "bisgo/app/P2P/client"
	"bisgo/app/manager"
	"bisgo/app/proving/client"
)

type Core struct {
	DB                          *DB.DBHandler
	Client                      *client.ProvingClient
	P2PClient                   *p2pclient.P2PClient
	ComplianceCheckStateManager *manager.ComplianceCheckStateManager
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore() *Core {
	return &Core{
		DB:                          DB.GetDBHandler(),
		Client:                      client.GetProvingClient(),
		P2PClient:                   p2pclient.GetP2PClient(),
		ComplianceCheckStateManager: manager.GetComplianceCheckStateManager(),
	}
}
