package main

import (
	"context"
	"log"
	"net/http"
	"restapi/db"
	"restapi/sessions"
	"restapi/users"
)

func main() {
	dbConnection := db.CreateNewConnection()
	mux := http.NewServeMux()
	usersController := users.NewController(users.NewUserService(users.NewUserStoreDB(dbConnection)))
	sessionsController := sessions.NewController(sessions.NewSessionService(sessions.NewSessionStoreDB(dbConnection)))
	defer dbConnection.Close(context.Background())

	mux.HandleFunc("GET /users", usersController.GetAllUsers)
	mux.HandleFunc("POST /session", sessionsController.Login)

	handler := Logging(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
