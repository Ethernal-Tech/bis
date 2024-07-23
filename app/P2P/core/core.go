package core

import (
	"bisgo/app/DB"
	p2pclient "bisgo/app/P2P/client"
	engine "bisgo/app/rules_engine"
)

type Core struct {
	DB          *DB.DBHandler
	P2PClient   *p2pclient.P2PClient
	RulesEngine *engine.RulesEngine
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore() *Core {
	return &Core{
		DB:          DB.GetDBHandler(),
		P2PClient:   p2pclient.GetP2PClient(),
		RulesEngine: engine.GetRulesEngine(),
	}
}
