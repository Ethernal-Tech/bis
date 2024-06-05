package core

import (
	"bisgo/app/DB"
	p2pclient "bisgo/app/P2P/client"
	provingclient "bisgo/app/proving/client"
	"bisgo/app/web/manager"
)

// Core is the central component of the system.
// It encapsulates and manages all core functionalities of the system.
type Core struct {
	DB                  *DB.DBHandler
	SessionManager      *manager.SessionManager
	SanctionListManager *manager.SanctionListManager
	P2PClient           *p2pclient.P2PClient
	ProvingClient       *provingclient.ProvingClient
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore() *Core {
	return &Core{
		DB:                  DB.GetDBHandler(),
		SessionManager:      manager.CreateSessionManager(),
		SanctionListManager: manager.CreateSanctionListManager(),
		P2PClient:           p2pclient.GetP2PClient(),
		ProvingClient:       provingclient.GetProvingClient(),
	}
}
