package users

import "fmt"

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(username string, password string) User {
	u := User{
		Username: username,
		Password: password,
	}
	return u
}

type UserStore interface {
	Add(user User) User
	Get(id int) (User, error)
	GetAll() []User
	FindByUsername(username string) (*User, bool)
}

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

func (u *UserStoreInMem) FindByUsername(username string) (*User, bool) {
	for _, user := range u.store {
		if user.Username == username {
			return &user, true
		}
	}
	return nil, false
}

func (u *UserStoreInMem) Add(user User) User {
	id := u.nextId
	user.Id = id
	u.store[id] = user
	u.nextId++
	return user
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
