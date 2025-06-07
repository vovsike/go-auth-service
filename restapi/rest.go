package main

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"restapi/db"
	"restapi/sessions"
	"restapi/users"
	"strconv"
)

type usersServer struct {
	userStore      users.UserStore
	sessionService *sessions.SessionService
}

func newUserServer(conn *pgx.Conn) *usersServer {
	per := users.NewUserStoreDB(conn)
	ss := sessions.NewSessionService(sessions.NewSessionStoreDB(conn))
	return &usersServer{
		userStore:      per,
		sessionService: ss}
}

func (ts *usersServer) GetUserById(w http.ResponseWriter, r *http.Request) {
	v, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := ts.userStore.GetById(v)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func (ts *usersServer) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	usersList := ts.userStore.GetAll()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersList)
}

func (ts *usersServer) Ping(w http.ResponseWriter, r *http.Request) {
	ts.userStore.Ping()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}

func (ts *usersServer) CreateNewSession(w http.ResponseWriter, r *http.Request) {
	ts.sessionService.CreateNewSession(1)
}

func main() {
	dbConnection := db.CreateNewConnection()
	mux := http.NewServeMux()
	us := newUserServer(dbConnection)
	defer us.userStore.(*users.UserStoreDB).Close()

	mux.Handle("GET /{id}", Auth(http.HandlerFunc(us.GetUserById)))
	mux.HandleFunc("GET /", us.GetAllUsers)
	mux.HandleFunc("GET /ping", us.Ping)
	mux.HandleFunc("POST /session", us.CreateNewSession)

	handler := Logging(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
