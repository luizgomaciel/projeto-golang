package repository

import (
	"context"
	jobs "encoder/application/jobs/accounts"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (repo AccountRepositoryDb) FindAllAccounts() (*[]jobs.Account, error) {
	collection := repo.Db.Collection(os.Getenv("MONGO_COLLECTION_ACCOUNT"))

	var startDate, endDate time.Time
	endDate = time.Now()
	startDate = endDate.AddDate(0, 0, -10)

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("failed to find document into MongoDB: FindAllAccounts: jobs.Account")
		return nil, err
	}
	defer cursor.Close(context.Background())

	accounts := make([]jobs.Account, 0)
	for cursor.Next(context.Background()) {
		var loan jobs.Account
		err := cursor.Decode(&loan)
		if err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, loan)
	}

	return &accounts, nil
}
