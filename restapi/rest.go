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

type usersServer struct {
	usersService   *users.UserService
	sessionService *sessions.SessionService
}

func newUserServer(conn *pgx.Conn) *usersServer {
	per := users.NewUserService(users.NewUserStoreDB(conn))
	ss := sessions.NewSessionService(sessions.NewSessionStoreDB(conn))
	return &usersServer{
		usersService:   per,
		sessionService: ss}
}

func (ts *usersServer) CreateNewSession(w http.ResponseWriter, r *http.Request) {
	ts.sessionService.CreateNewSession(1)
}

func main() {
	dbConnection := db.CreateNewConnection()
	mux := http.NewServeMux()
	us := newUserServer(dbConnection)
	defer dbConnection.Close(context.Background())

	mux.HandleFunc("POST /session", us.CreateNewSession)

	handler := Logging(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
