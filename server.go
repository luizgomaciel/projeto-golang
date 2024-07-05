package main

import (
	"encoder/application/services"
	"encoder/application/workers"
	"encoder/graph"
	"encoder/infrastructure/config"
	"encoder/model/repository"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

var db config.Database
var dbMongoConnection = make(chan *mongo.Database)

const defaultPort = "8080"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		initGraphQL()
	}()

	go func() {
		defer wg.Done()
		initTDM()
	}()

	wg.Wait()
}

func initTDM() {
	log.Println("Iniciando TDM")

	messageChannel := make(chan amqp.Delivery)

	dbConnection, err := db.Connect()
	dbMongoConnection <- dbConnection

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

func initGraphQL() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	repoAccount := repository.AccountRepositoryDb{
		Db: <-dbMongoConnection,
	}

	accountService := services.NewAccountService()
	accountService.AccountRepository = repoAccount

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		AcccountService: &accountService,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
