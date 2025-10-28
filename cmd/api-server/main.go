package main

import (
	"awesomeProject/internal/user"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Dependency Injection
	var store user.Store = user.NewInMemStore()
	var service user.Service = user.NewInMemoryUserService(store)
	var handler user.Handler = user.Handler{Service: service}

	r := chi.NewRouter()
	r.Post("/user", ErrorHandler(handler.CreateUser))
	r.Post("/authenticate", ErrorHandler(handler.Authenticate))
	r.Get("/user", ErrorHandler(handler.SearchUser))
	r.Get("/user/{id}", ErrorHandler(handler.GetUser))

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
