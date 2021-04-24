package tests

import (
	"context"
	"errors"

	"github.com/UsagiBooru/accounts-server/models/mongo_models"
	"github.com/UsagiBooru/accounts-server/utils/server"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func ReGenerateTestDatabase() error {
	conf := server.GetConfig()
	m := server.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass)
	// Drop database
	drops := []string{"users", "invites", "sequence"}
	for _, d := range drops {
		col := m.Database("accounts").Collection(d)
		err := col.Drop(context.Background())
		if err != nil {
			return err
		}
	}
	// Get password hash
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(PASSWORD),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("password hash create failed")
	}
}

func DestroyMongoTestContainer(pool *dockertest.Pool, resource *dockertest.Resource) {
	// When you're done, kill and remove the container
	if err := pool.Purge(resource); err != nil {
		server.Fatal(err.Error())
	}
}

func ReGenerateDatabase(m *mongo.Client) error {
	// Drop database
	drops := []string{"users", "invites", "sequence"}
	for _, d := range drops {
		col := m.Database("accounts").Collection(d)
		err := col.Drop(context.Background())
		if err != nil {
			return err
		}
	}
	if err := InitAccountDatabase(m); err != nil {
		return err
	}
	return nil
}
