package core

import (
	"bisgo/app/DB"
	"bisgo/app/P2P/client"
)

type Core struct {
	DB     *DB.DBHandler
	Client *client.P2PClient
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore() *Core {
	return &Core{
		DB:     DB.GetDBHandler(),
		Client: client.GetP2PClient(),
	}
}
