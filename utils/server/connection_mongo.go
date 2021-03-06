package server

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDBClient creates a new mongodb client
func NewMongoDBClient(host, user, pass string) *mongo.Client {
	// Initialize MongoDB driver
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://"+user+":"+pass+"@"+host),
	)
	if err != nil {
		Error("Connect to mongodb failed: " + err.Error())
	}
	// Debug("MongoDB client created.")
	return client
}
