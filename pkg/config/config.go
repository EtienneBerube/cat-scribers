package config

import "os"

// config is a container of all necessary environment variables used in the code.
type config struct {
	Port       string
	MongoDBURL string
	JWTSecret  string
}

// Config is a globally accessible instance of the config struct
var Config config

// Init Initializes the global Config struct and fetches the data from the environment variables
func Init() {
	Config = config{
		Port:       os.Getenv("HTTP_PORT"),
		MongoDBURL: os.Getenv("MONGODB_URL"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}
