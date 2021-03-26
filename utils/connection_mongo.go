package utils

import (
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoDBConnection(host, user, pass string) {
	// Initalize MongoDB driver
	err := mgm.SetDefaultConfig(
		nil,
		"mgm_lab",
		options.Client().ApplyURI("mongodb://"+user+":"+pass+"@"+host),
	)
	if err != nil {
		log.Fatalf("Connect to mongodb failed: %s", err)
	}
}
