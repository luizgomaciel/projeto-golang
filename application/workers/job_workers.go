package workers

import (
	jobs "encoder/application/jobs/accounts"
	"encoder/application/services"
	"encoding/json"

	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     *jobs.JobAccount
	Message *amqp.Delivery
	Error   error
}

type AccountMessage struct {
	typeAccount int
	number      int
}

func JobWorker(messageChannel chan amqp.Delivery, jobService services.AcccountService) error {

	for message := range messageChannel {
		am := AccountMessage{}

		err := json.Unmarshal(message.Body, am)
		if err != nil {
			return err
		}

		jobService.InsertAccount(am.typeAccount, am.number)
	}

	return nil
}

func Process(messageChannel chan amqp.Delivery, jobService services.AcccountService) error {
	go JobWorker(messageChannel, jobService)

	return nil
}
