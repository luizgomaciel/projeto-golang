package repository

import (
	"context"
	jobs "encoder/application/jobs/loan"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoanRepository interface {
	Insert(ac *jobs.Loan) (*jobs.Loan, error)
}

type LoanRepositoryDb struct {
	Db *mongo.Database
}

func (repo LoanRepositoryDb) Insert(ac jobs.Loan) (*jobs.Loan, error) {
	collection := repo.Db.Collection(os.Getenv("MONGO_COLLECTION_LOAN"))

	_, err := collection.InsertOne(context.Background(), ac)
	if err != nil {
		log.Fatalf("failed to insert document into MongoDB: ")
		return nil, err
	}

	return &ac, nil
}

func (repo LoanRepositoryDb) FindAllByAccountNumber(accountNumber string) (*jobs.Loan, error) {
	collection := repo.Db.Collection(os.Getenv("MONGO_COLLECTION_LOAN"))
	filter := bson.M{"account_number": accountNumber}

	single := collection.FindOne(context.Background(), filter)
	if single.Err() != nil {
		return nil, single.Err()
	}

	var loan jobs.Loan
	err := single.Decode(&loan)
	if err != nil {
		log.Fatalf("failed to decode document: %v", err)
		return nil, err
	}

	return &loan, nil
}
