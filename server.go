package main

import (
	"encoder/application/workers"
	"encoder/infrastructure/config"
	"log"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db config.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	log.Println("Iniciando TDM")

	messageChannel := make(chan amqp.Delivery)

	dbConnection, err := db.Connect()

	if err != nil {
		log.Fatalf("error connecting to DB")
	}

	rabbitMQ := config.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := workers.NewJobManager(dbConnection, rabbitMQ, messageChannel)
	jobManager.Start(ch)

	log.Println("Iniciou TDM")

}
