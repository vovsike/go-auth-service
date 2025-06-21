package main

import (
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	"restapi/accounts"
	"restapi/internal/db"
	"restapi/internal/jwtInternal"
	"restapi/internal/mq"
	"restapi/internal/observability"
)

func main() {
	// Load env
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	tp, err := observability.NewTraceProvider()
	if err != nil {
		log.Fatal("Failed to create trace provider", err)
	}

	otel.SetTracerProvider(tp)

	// Init DBs
	dbpool, err := db.CreateConnectionPool()
	if err != nil {
		log.Fatal("Error creating DB connection", err)
	}

	//setup mq
	mqClient, err := mq.NewClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to create MQ client:", err)
	}
	defer mqClient.Close()

	// Declare the exchange
	err = mqClient.DeclareExchange("accounts", "direct", false, false, false, false)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	// Init MUX
	mux := http.NewServeMux()

	// Init services
	aService := accounts.NewServiceDB(dbpool, mqClient)
	jwtService := jwtInternal.NewService()

	// Init controllers
	accountController := accounts.NewController(aService, jwtService)

	defer dbpool.Close()

	// Register handlers
	mux.HandleFunc("POST /account", accountController.CreateNewAccount)

	// Register global middleware
	handler := Logging(CountRequests(mux))

	// Start
	log.Fatal(http.ListenAndServe(":8080", handler))
}
