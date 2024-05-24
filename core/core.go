package core

import (
	"bisgo/config"
	"bisgo/core/DB"
	"bisgo/core/manager"
)

// Core is the central component of the system.
// It encapsulates and manages all core functionalities of the system.
type Core struct {
	DB                  *DB.DBHandler
	SessionManager      *manager.SessionManager
	SanctionListManager *manager.SanctionListManager
	Config              *config.Config
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore(config *config.Config) *Core {
	return &Core{
		DB:                  DB.CreateDBHandler(config.ResolveDBAddress(), config.ResolveDBPassword(), config.ResolveDBPort(), config.ResolveDBName()),
		SessionManager:      manager.CreateSessionManager(),
		SanctionListManager: manager.CreateSanctionListManager(),
		Config:              config,
	}
}
