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
	per := users.NewUserStoreInMem()
	return &usersServer{store: per}
}

func (ts *usersServer) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	u := users.User{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&u)
	ts.store.Add(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (ts *usersServer) GetUserById(w http.ResponseWriter, r *http.Request) {
	v, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := ts.store.Get(v)
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

func main() {
	mux := http.NewServeMux()
	ts := newUserServer()
	mux.HandleFunc("POST /", ts.CreateNewUser)
	mux.HandleFunc("GET /{id}", ts.GetUserById)
	mux.HandleFunc("GET /", ts.GetAllUsers)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
