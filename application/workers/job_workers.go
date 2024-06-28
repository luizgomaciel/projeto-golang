package workers

import (
	jobs "encoder/application/jobs/accounts"
	"encoder/application/services"
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
	TypeAccount int `json:"typeAccount"`
	Number      int `json:"number"`
}

func JobWorker(messageChannel <-chan amqp.Delivery, acccountService services.AcccountService, workerID int, jobWorkerStack chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()
	for message := range messageChannel {
		am := AccountMessage{}

		err := json.Unmarshal(message.Body, &am)
		if err != nil {
			log.Fatalf("erro no Unmarshal")
			return err
		}
		acccountService.InsertAccount(am.TypeAccount, am.Number)

		<-jobWorkerStack

	}

	return nil
}
