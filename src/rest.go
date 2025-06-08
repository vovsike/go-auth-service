package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"restapi/db"
	"restapi/jwtInternal"
	"restapi/sessions"
	"restapi/users"
)

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Init DBs
	dbpool, err := db.CreateConnectionPool()
	if err != nil {
		log.Fatal("Error creating DB connection", err)
	}

	// Init MUX
	mux := http.NewServeMux()

	// Init stores
	uStore := users.NewUserStoreDB(dbpool)
	sStore := sessions.NewSessionStoreDB(dbpool)

	// Init services
	uService := users.NewUserService(uStore)
	sService := sessions.NewSessionService(sStore)
	jwtService := jwtInternal.NewService()

	// Init controllers
	usersController := users.NewController(uService)
	sessionsController := sessions.NewController(sService, uService, jwtService)

	defer dbpool.Close()

	// Register handlers
	mux.HandleFunc("POST /users", usersController.CreateUser)
	mux.HandleFunc("POST /session", sessionsController.Login)
	mux.HandleFunc("POST /session/token", sessionsController.GetToken)

	// Register global middleware
	handler := Logging(mux)

	// Start
	log.Fatal(http.ListenAndServe(":8080", handler))
}
