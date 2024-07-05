package services_test

import (
	"context"
	jobs "encoder/application/jobs/accounts"
	"encoder/application/services"
	"encoder/infrastructure/config"
	"encoder/model/repository"
	"encoding/json"
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

func prepareFind() repository.AccountRepositoryDb {
	db := config.NewDbTest()

	repoAccount := repository.AccountRepositoryDb{
		Db: db,
	}

	return repoAccount
}

func TestInsertAccount(t *testing.T) {
	_, repo := prepare()

	accountService := services.NewAccountService()
	accountService.AccountRepository = repo

	_, err := accountService.InsertAccount(2, 10)
	require.Nil(t, err)

	ctx := context.Background()
	defer repo.Db.Client().Disconnect(ctx)
}

func TestFindAvailableAccounts(t *testing.T) {
	repoAccount := prepareFind()

	accountService := services.NewAccountService()
	accountService.AccountRepository = repoAccount

	accounts, err := accountService.FindAvailableAccounts()
	require.Nil(t, err)
	require.NotNil(t, accounts)

	accountsJSON, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal accounts: %v", err)
	}

	log.Println("DATAS: ", string(accountsJSON))

	ctx := context.Background()
	defer repoAccount.Db.Client().Disconnect(ctx)
}
