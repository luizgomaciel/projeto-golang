package producer

import (
	"encoder/graph/model"
	"encoder/infrastructure/config"
	"encoding/json"
	"os"
)

type AccountProducer struct {
	Rabbit *config.RabbitMQ
}

type Message struct {
	TypeAccount int
	Qtd         int
	Products    []string
}

func NewAccountProducer() AccountProducer {
	rabbit := config.NewRabbitMQ()
	return AccountProducer{
		Rabbit: rabbit,
	}
}

func (ap *AccountProducer) Produce(request model.JobQueue) (*model.JobQueueResponse, error) {
	ch := ap.Rabbit.Connect()
	defer ch.Close()

	exchange := os.Getenv("RABBITMQ_NOTIFICATION_EX")
	routingKey := os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY")
	contentType := "GOLANG"
	message := Message{
		TypeAccount: request.TypeAccount,
		Qtd:         request.Quantity,
		Products:    request.Products,
	}

	jsonStr, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	err = ap.Rabbit.Notify(string(jsonStr), contentType, exchange, routingKey)
	if err != nil {
		return nil, err
	}

	resp := true
	return &model.JobQueueResponse{
		IsStarted: &resp,
	}, nil
}
