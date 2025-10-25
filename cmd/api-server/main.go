package main

import (
	"awesomeProject/internal/user"
	"fmt"

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
	fmt.Println(service.Authenticate("admin@example.com", "password"))
}
