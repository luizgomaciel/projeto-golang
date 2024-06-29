package repository

import (
	"context"
	jobs "encoder/application/jobs/loan"
	"log"
	"os"

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
