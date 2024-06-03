package core

import (
	"bisgo/app/DB"
	"bisgo/app/proving/client"
)

type Core struct {
	DB     *DB.DBHandler
	Client *client.ProvingClient
}

// CreateCore function initializes and returns a new instance of the Core component.
func CreateCore() *Core {
	return &Core{
		DB:     DB.GetDBHandler(),
		Client: client.GetProvingClient(),
	}
}
