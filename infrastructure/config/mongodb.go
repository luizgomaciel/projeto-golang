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
		log.Fatalf("Test db error: ")
	}

	return connection
}

func (d *Database) Connect() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_CONNECT"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(os.Getenv("MONGO_DATABASE"))

	return db, nil
}
