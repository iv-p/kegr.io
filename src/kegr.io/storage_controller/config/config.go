package config

import "os"

// Config is used to deserialize the yaml config file
// into an usable golang structure we can pass around
type Config struct {
	DataRoot         string
	Address          string
	APIPort          string
	InternalGrpcPort string
	ExternalGrpcPort string
	Other            string
	MachineName      string
	LiquidExtension  string
	KegFile          string
}

// C is the config instance
var C *Config

// Load initialises the config properties
func Load() {
	C = &Config{
		DataRoot:         getenv("DATA_ROOT", "./www"),
		Address:          getenv("ADDRESS", "localhost"),
		APIPort:          getenv("API_PORT", "8080"),
		InternalGrpcPort: getenv("INTERNAL_GRPC_PORT", "23471"),
		ExternalGrpcPort: getenv("EXTERNAL_GRPC_PORT", "24471"),
		Other:            getenv("OTHER", ""),
		MachineName:      getenv("MACHINE_NAME", "pesho"),
		LiquidExtension:  "liquid",
		KegFile:          ".keg",
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
