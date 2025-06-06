package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
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

	type createUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	u := createUser{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&u)
	bhash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	createdUser := ts.store.Add(users.User{
		Username: u.Username,
		Password: string(bhash),
	})
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (ts *usersServer) VerifyUserPassword(w http.ResponseWriter, r *http.Request) {
	type login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	l := login{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&l)

	u, found := ts.store.FindByUsername(l.Username)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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
	mux.Handle("GET /{id}", Auth(http.HandlerFunc(ts.GetUserById)))
	mux.HandleFunc("GET /", ts.GetAllUsers)
	mux.HandleFunc("POST /auth", ts.VerifyUserPassword)

	handler := Logging(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
