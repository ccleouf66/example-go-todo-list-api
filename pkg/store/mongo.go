package store

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoreOption struct {
	// Endpoints list of db in style ADDR:PORT
	Address string
	// User used with authentication procedure user:password
	User string
	// Password used with authentication procedure user:password
	Password string
	// TLS config object contenning client certificate for cert authentication procedure
	TLSConfig *tls.Config
	// DbName to use
	DbName string
	// RsName to connect in replicat set mode
	RsName string
	// AUTH_SOURCE the db where user is stored
	AuthSource string
}

type MongoStore struct {
	Client   *mongo.Client
	DataBase *mongo.Database
}

func NewMongoStore(c context.Context, conf *StoreOption) (*MongoStore, error) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Build creds ops
	creds := options.Credential{
		AuthSource: conf.AuthSource,
		Username:   conf.User,
		Password:   conf.Password,
	}

	// Create new mongo client with ops
	opts := options.Client().ApplyURI(conf.Address).SetAuth(creds)

	// Connect client to mongo instance
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Select and init the database
	db := client.Database(conf.DbName)
	err = initDb(c, db)
	if err != nil {
		return nil, err
	}

	ms := &MongoStore{
		Client:   client,
		DataBase: db,
	}

	return ms, nil
}

func initDb(c context.Context, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(c, time.Second*2)
	defer cancel()

	// Create user collection with index on email and name field as unique value
	tasksCollection := db.Collection("tasks")

	tasksNameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	indexNames, err := tasksCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{tasksNameIndex})
	if err != nil {
		return err
	}
	for _, indexName := range indexNames {
		log.Printf("Index %s created\n", indexName)
	}

	return nil
}
