package users

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore interface {
	GetAll() []User
	FindByUsername(username string) (*User, bool)
	GetById(id int) (User, error)
	Ping()
	Add(username string, password string) User
}
