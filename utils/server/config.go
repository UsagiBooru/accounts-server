package server

import (
	"os"

	"github.com/joho/godotenv"
)

// ConfigList stores credentials
type ConfigList struct {
	MongoHost   string
	MongoUser   string
	MongoPass   string
	ElasticHost string
	ElasticUser string
	ElasticPass string
	JwtSecret   string
}

// GetConfig creates ConfigList from environment variables
func GetConfig() ConfigList {
	// Load ini
	if envFilePath := os.Getenv("GO_ENV"); envFilePath != "" {
		if err := godotenv.Load(envFilePath); err != nil {
			Fatal("envFile " + envFilePath + "could not be loaded")
		}
	} else {
		_ = godotenv.Load(".env")
	}
	// Parse to ConfigList struct
	return ConfigList{
		MongoHost:   os.Getenv("MONGO_HOST"),
		MongoUser:   os.Getenv("MONGO_USER"),
		MongoPass:   os.Getenv("MONGO_PASS"),
		ElasticHost: os.Getenv("ELASTIC_HOST"),
		ElasticUser: os.Getenv("ELASTIC_USER"),
		ElasticPass: os.Getenv("ELASTIC_PASS"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
	}
}
