package repository

import (
	"context"
	jobs "encoder/application/jobs/accounts"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type AccountRepository interface {
	Insert(ac *jobs.Account) (*jobs.Account, error)
}

type AccountRepositoryDb struct {
	Db *mongo.Database
}

func (repo AccountRepositoryDb) Insert(ac *jobs.Account) (*jobs.Account, error) {
	collection := repo.Db.Collection(os.Getenv("MONGO_COLLECTION_ACCOUNT"))

	_, err := collection.InsertOne(context.Background(), ac)
	if err != nil {
		log.Fatalf("failed to insert document into MongoDB: ")
		return nil, err
	}

	return ac, nil
}
