package workers

import (
	jobs "encoder/application/jobs/accounts"
	"encoder/application/services"
	"encoder/infrastructure/config"
	"encoder/model/repository"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobManager struct {
	Db             *mongo.Database
	Domain         jobs.JobAccount
	MessageChannel chan amqp.Delivery
	JobWorkerStack chan struct{}
	RabbitMQ       *config.RabbitMQ
}

func NewJobManager(db *mongo.Database, rabbitMQ *config.RabbitMQ, messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		Db:             db,
		MessageChannel: messageChannel,
		JobWorkerStack: make(chan struct{}, 10), // Limita o número de workers ativos
		RabbitMQ:       rabbitMQ,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS.")
	}

	acccountService := services.NewAccountService()
	acccountService.AccountRepository = repository.AccountRepositoryDb{Db: j.Db}

	var wg sync.WaitGroup

	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		log.Println("Rodando process:", strconv.Itoa(qtdProcesses))
		wg.Add(1)
		go JobWorker(j.MessageChannel, acccountService, qtdProcesses, j.JobWorkerStack, &wg)
	}

	go func() {
		for msg := range j.MessageChannel {
			log.Println("restart process:", msg)
			j.JobWorkerStack <- struct{}{} // Adiciona uma tarefa ao stack
			go JobWorker(j.MessageChannel, acccountService, -1, j.JobWorkerStack, &wg)
			log.Println(msg)
		}
	}()

	wg.Wait() // Aguarda até que todas as goroutines terminem
}
