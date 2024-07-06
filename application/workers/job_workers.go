package workers

import (
	jobs "encoder/application/jobs/accounts"
	"encoder/application/services"
	"encoder/application/utils"
	"encoder/model/repository"
	"encoding/json"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type JobWorkerStack struct {
	Job     *jobs.JobAccount
	Message *amqp.Delivery
	Error   error
}

type AccountMessage struct {
	TypeAccount      int      `json:"typeAccount"`
	NumberOfAccounts int      `json:"qtd"`
	Products         []string `json:"products"`
}

func JobWorker(messageChannel <-chan amqp.Delivery, acccountService services.AcccountService, workerID int, jobWorkerStack chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()
	for message := range messageChannel {
		am := AccountMessage{}

		err := json.Unmarshal(message.Body, &am)
		if err != nil {
			log.Fatalf("erro no Unmarshal:", err)
			return err
		}

		err = executeJobs(acccountService, am)
		if err != nil {
			return err
		}

		<-jobWorkerStack

	}

	return nil
}

func executeJobs(acccountService services.AcccountService, am AccountMessage) error {
	accounts, err := createAccount(acccountService, am)
	if err != nil {
		return err
	}

	if utils.Contains(am.Products, "LOAN_PRODUCT") {
		accounts, err = insertLoan(accounts, acccountService, am)
		if err != nil {
			return err
		}
	}

	return nil
}

func createAccount(acccountService services.AcccountService, am AccountMessage) (*[]jobs.Account, error) {
	return acccountService.InsertAccount(am.TypeAccount, am.NumberOfAccounts)
}

func insertLoan(accounts *[]jobs.Account, acccountService services.AcccountService, am AccountMessage) (*[]jobs.Account, error) {
	loanService := services.NewLoanService()
	loanService.LoanRepository = repository.LoanRepositoryDb{Db: acccountService.AccountRepository.Db}

	_, err := loanService.InsertLoan(accounts)
	return accounts, err
}
