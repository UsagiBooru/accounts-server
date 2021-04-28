package tests

import (
	"context"
	"errors"
	"flag"
	"time"

	"github.com/UsagiBooru/accounts-server/utils/server"
	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var parallelFlag = flag.Bool("docker", false, "Set true to use parallel test(Local), otherwise to simple test(CI)")

// GetDatabaseConnection creates connection to database
func GetDatabaseConnection() (db *mongo.Client, shutdown func(), isParallel bool) {
	var err error
	if *parallelFlag {
		server.Debug("Using mongo container")
		db, shutdown, err = generateMongoTestContainer()
		if err != nil {
			server.Fatal(err.Error())
		}
		isParallel = true
	} else {
		server.Debug("Using mongo server")
		conf := server.GetConfig()
		db = server.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass)
		shutdown = func() {}
		err = nil
		isParallel = false
	}
	if err := reGenerateDatabase(db); err != nil {
		server.Fatal(err.Error())
	}
	return
}

func generateMongoTestContainer() (*mongo.Client, func(), error) {
	var db *mongo.Client
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, errors.New("Could not connect to docker: " + err.Error())
	}
	resource, err := pool.Run("mongo", "4.4.5", nil)
	// Force delete after 5 minutes
	resource.Expire(300)
	if err != nil {
		return nil, nil, errors.New("Could not start resource: " + err.Error())
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		db, err = mongo.Connect(
			ctx,
			options.Client().ApplyURI("mongodb://localhost:"+string(resource.GetPort("27017/tcp"))),
		)
		if err != nil {
			return err
		}
		return db.Ping(ctx, readpref.Primary())
	}); err != nil {
		return nil, nil, err
	}
	return db, func() { destroyMongoTestContainer(pool, resource) }, nil
}

func destroyMongoTestContainer(pool *dockertest.Pool, resource *dockertest.Resource) {
	// When you're done, kill and remove the container
	if err := pool.Purge(resource); err != nil {
		server.Fatal(err.Error())
	}
}

func reGenerateDatabase(m *mongo.Client) error {
	// Drop database
	drops := []string{"users", "invites", "mutes", "sequence"}
	for _, d := range drops {
		col := m.Database("accounts").Collection(d)
		err := col.Drop(context.Background())
		if err != nil {
			return err
		}
	}
	if err := initAccountDatabase(m); err != nil {
		return err
	}
	if err := initMuteDatabase(m); err != nil {
		return err
	}
	return nil
}
