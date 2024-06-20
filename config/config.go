package config

import (
	"log"
	"os"
)

func ResolveP2PNodeAPIAddress() string {
	env_address := os.Getenv("P2P_NODE_ADDRESS")
	if env_address == "" {
		return "http://localhost:5000/passthrough"
	}

	return env_address
}

func ResolvePeerID() string {
	env_val := os.Getenv("PEER_ID")
	if env_val == "" {
		log.Fatal("PEER_ID env variable must be set")
	}

	return env_val
}

func ResolveGpjcApiAddress() string {
	env_address := os.Getenv("GPJC_API_ADDRESS")
	if env_address == "" {
		return "localhost"
	}

	return env_address
}

func ResolveGpjcApiPort() string {
	env_address := os.Getenv("GPJC_API_PORT")
	if env_address == "" {
		return "9090"
	}

	return env_address
}

func ResolveGpjcPort() string {
	env_address := os.Getenv("GPJC_PORT")
	if env_address == "" {
		return "10501"
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
		return "BIS1"
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
		return ":4001"
	}

	return ":" + env_port
}

func ResovleIsCentralBank() bool {
	if os.Getenv("IS_CENTRAL_BANK") == "" {
		return false
	}

	return true
}

func ResolveMyGlobalIdentifier() string {
	env_global_ident := os.Getenv("MY_GLOBAL_IDENTIFIER")
	if env_global_ident == "" {
		log.Fatal("MY_GLOBAL_IDENTIFIER environment variable must be set for every bank")
	}

	return env_global_ident
}

func ResolveCBGlobalIdentifier() string {
	env_global_ident := os.Getenv("CB_GLOBAL_IDENTIFIER")
	if env_global_ident == "" {
		log.Fatal("Expected to have defined CB_GLOBAL_IDENTIFIER for this bank")
	}

	return env_global_ident
}

func ResolveJurisdictionCode() string {
	env_code := os.Getenv("JURISDICTION_CODE")
	if env_code == "" {
		log.Fatal("Expected to have defined JURISDICTION_CODE for this bank")
	}

	return env_code
}
