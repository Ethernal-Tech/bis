package config

import "os"

func ResolveP2PNodeAPIAddress() string {
	env_address := os.Getenv("P2P_NODE_ADDRESS")
	if env_address == "" {
		return "localhost:5000"
	}

	return env_address
}

func ResolveGpjcApiAddress() string {
	env_address := os.Getenv("GPJC_API")
	if env_address == "" {
		return "localhost"
	}

	return env_address
}

func ResolveGpjcClientUrl() string {
	env_url := os.Getenv("GPJC")
	if env_url == "" {
		return "0.0.0.0"
	}

	return env_url
}

func ResolveDBAddress() string {
	env_address := os.Getenv("DB_ADDRESS")
	if env_address == "" {
		return "localhost"
	}

	return env_address
}

func ResolveDBPort() string {
	env_port := os.Getenv("DB_PORT")
	if env_port == "" {
		return "1433"
	}

	return env_port
}

func ResolveDBName() string {
	env_name := os.Getenv("DB_NAME")
	if env_name == "" {
		return "BIS"
	}

	return env_name
}

func ResolveDBPassword() string {
	env_password := os.Getenv("DB_PASSWORD")
	if env_password == "" {
		return "Ethernal123"
	}

	return env_password
}

func ResolveServerPort() string {
	env_port := os.Getenv("SERVER_PORT")
	if env_port == "" {
		return ":4000"
	}

	return ":" + env_port
}
