package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Db *mongo.Database
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *mongo.Database {
	dbInstance := NewDb()

	connection, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (d *Database) Connect() (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_CONNECT")))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := client.Database(os.Getenv("MONGO_DATABASE"))

	return db, nil
}
