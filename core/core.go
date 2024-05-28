package core

import (
	"bisgo/config"
	"bisgo/core/DB"
	"bisgo/core/manager"
	"bisgo/core/messaging"
)

// Core is the central component of the system.
// It encapsulates and manages all core functionalities of the system.
type Core struct {
	DB                  *DB.DBHandler
	SessionManager      *manager.SessionManager
	SanctionListManager *manager.SanctionListManager
	MessagingHandler    *messaging.MessagingHandler
	Config              *config.Config
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore(config *config.Config) *Core {
	return &Core{
		DB:                  DB.CreateDBHandler(config.ResolveDBAddress(), config.ResolveDBPassword(), config.ResolveDBPort(), config.ResolveDBName()),
		SessionManager:      manager.CreateSessionManager(),
		SanctionListManager: manager.CreateSanctionListManager(),
		MessagingHandler:    messaging.CreateMessagingHandler(config.P2PNodeAPIAddress),
		Config:              config,
	}
}
