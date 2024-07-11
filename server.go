package main

import (
	"encoder/application/services"
	"encoder/application/workers"
	"encoder/graph"
	"encoder/grpc"

	"encoder/infrastructure/config"
	"encoder/model/repository"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"

	"go.mongodb.org/mongo-driver/mongo"
	gr "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var db config.Database
var dbMongoConnection = make(chan *mongo.Database)
var dbMongoConnection2 = make(chan *mongo.Database)

const defaultPort = "8080"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()
		initGraphQL()
	}()

	go func() {
		defer wg.Done()
		initTDM()
	}()

	go func() {
		defer wg.Done()
		initGRPC()
	}()

	go func() {
		defer wg.Done()
		initRest()
	}()

	wg.Wait()
}

func initGRPC() {

	repoAccount := repository.AccountRepositoryDb{
		Db: <-dbMongoConnection2,
	}

	accountService := services.NewAccountService()
	accountService.AccountRepository = repoAccount

	grpcServer := gr.NewServer()
	grpc.RegisterAccountServiceRequestServer(grpcServer, accountService)
	reflection.Register(grpcServer)
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	log.Printf("connect to tcp://localhost%s/ for gRPC", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func initTDM() {
	log.Println("Iniciando TDM")

	messageChannel := make(chan amqp.Delivery)

	dbConnection, err := db.Connect()
	dbMongoConnection <- dbConnection
	dbMongoConnection2 <- dbConnection

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

func initRest() {

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
