package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"restapi/db"
	"restapi/sessions"
	"restapi/users"
)

func main() {
	// Init DBs
	dbConnection := db.CreateNewConnection()

	// Init MUX
	mux := http.NewServeMux()

	// Init stores
	uStore := users.NewUserStoreDB(dbConnection)
	sStore := sessions.NewSessionStoreDB(dbConnection)

	// Init services
	uService := users.NewUserService(uStore)
	sService := sessions.NewSessionService(sStore)

	// Init controllers
	usersController := users.NewController(uService)
	sessionsController := sessions.NewController(sService, uService)

	defer func(dbConnection *pgx.Conn, ctx context.Context) {
		_ = dbConnection.Close(ctx)
	}(dbConnection, context.Background())

	// Register handlers
	mux.HandleFunc("GET /users", usersController.GetAllUsers)
	mux.HandleFunc("POST /users", usersController.CreateUser)
	mux.HandleFunc("POST /session", sessionsController.Login)

	// Register global middleware
	handler := Logging(mux)

	// Start
	log.Fatal(http.ListenAndServe(":8080", handler))
}
