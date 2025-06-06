package main

import (
	"encoding/json"
	"log"
	"net/http"
	"restapi/users"
	"strconv"
)

type usersServer struct {
	store users.UserStore
}

func newUserServer() *usersServer {
	per := users.NewUserStoreDB()
	return &usersServer{store: per}
}

func (ts *usersServer) GetUserById(w http.ResponseWriter, r *http.Request) {
	v, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := ts.store.GetById(v)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func (ts *usersServer) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	usersList := ts.store.GetAll()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersList)
}

func (ts *usersServer) Ping(w http.ResponseWriter, r *http.Request) {
	ts.store.Ping()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}

func main() {
	mux := http.NewServeMux()
	ts := newUserServer()
	defer ts.store.(*users.UserStoreDB).Close()
	mux.Handle("GET /{id}", Auth(http.HandlerFunc(ts.GetUserById)))
	mux.HandleFunc("GET /", ts.GetAllUsers)
	mux.HandleFunc("GET /ping", ts.Ping)

	handler := Logging(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
