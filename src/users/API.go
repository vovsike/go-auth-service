package users

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore interface {
	FindByUsername(username string) (*User, bool)
	GetById(id int) (User, error)
	Add(username string, password string) (User, error)
}
