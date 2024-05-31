package core

import (
	"bisgo/app/DB"
	"bisgo/app/P2P/client"
	"bisgo/app/web/manager"
)

// Core is the central component of the system.
// It encapsulates and manages all core functionalities of the system.
type Core struct {
	DB                  *DB.DBHandler
	SessionManager      *manager.SessionManager
	SanctionListManager *manager.SanctionListManager
	P2PClient           *client.P2PClient
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore() *Core {
	return &Core{
		DB:                  DB.GetDBHandler(),
		SessionManager:      manager.CreateSessionManager(),
		SanctionListManager: manager.CreateSanctionListManager(),
		P2PClient:           client.GetP2PClient(),
	}
}
