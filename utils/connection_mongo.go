package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient(host, user, pass string) *mongo.Client {
	// Initalize MongoDB driver
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://"+user+":"+pass+"@"+host),
	)
	if err != nil {
		log.Fatalf("Connect to mongodb failed: %s", err)
	}
	return client
}
