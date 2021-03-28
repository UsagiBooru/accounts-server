package utils

import (
	"log"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	MongoHost   string
	MongoUser   string
	MongoPass   string
	ElasticHost string
	ElasticUser string
	ElasticPass string
	JwtSecret   string
}

func GetConfig() ConfigList {
	// Load ini
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("config.ini could not be loaded.")
	}

	// Parse to ConfigList struct
	mongoSection := cfg.Section("MongoDB")
	elasticSection := cfg.Section("ElasticSearch")
	secretSection := cfg.Section("Secret")
	return ConfigList{
		MongoHost:   mongoSection.Key("host").String(),
		MongoUser:   mongoSection.Key("user").String(),
		MongoPass:   mongoSection.Key("pass").String(),
		ElasticHost: elasticSection.Key("host").String(),
		ElasticUser: elasticSection.Key("user").String(),
		ElasticPass: elasticSection.Key("pass").String(),
		JwtSecret:   secretSection.Key("jwt").String(),
	}
}
