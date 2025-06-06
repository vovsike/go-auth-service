package users

import "fmt"

type UserStoreInMem struct {
	store  map[int]User
	nextId int
}

func NewUserStoreInMem() *UserStoreInMem {
	return &UserStoreInMem{
		store:  make(map[int]User),
		nextId: 0,
	}
}

func (u *UserStoreInMem) Get(id int) (User, error) {
	val, ok := u.store[id]
	if ok {
		return val, nil
	}
	return User{}, fmt.Errorf("no user with id %d exists", id)
}

func (u *UserStoreInMem) GetAll() []User {
	var users []User
	for _, value := range u.store {
		users = append(users, value)
	}
	return users
}
