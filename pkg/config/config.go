package config

import "os"

type config struct {
	Port       string
	MongoDBURL string
	JWTSecret  string
}

var Config config

func Init() {
	Config = config{
		Port:       os.Getenv("HTTP_PORT"),
		MongoDBURL: os.Getenv("MONGODB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
