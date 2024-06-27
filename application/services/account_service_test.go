package services_test

import (
	"context"
	jobs "encoder/application/jobs/accounts"
	"encoder/application/services"
	"encoder/infrastructure/config"
	"encoder/model/repository"
	"log"
	"testing"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*jobs.Account, repository.AccountRepositoryDb) {
	db := config.NewDbTest()

	account := jobs.NewAccount(2)
	account.ID = uuid.NewV4().String()

	repo := repository.AccountRepositoryDb{
		Db: db,
	}

	return account, repo
}

func TestInsertAccount(t *testing.T) {
	_, repo := prepare()

	accountService := services.NewAccountService()
	accountService.AccountRepository = repo

	err := accountService.InsertAccount(2, 10)
	require.Nil(t, err)

	ctx := context.Background()
	defer repo.Db.Client().Disconnect(ctx)
}
