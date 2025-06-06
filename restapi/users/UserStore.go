package users

type UserStore interface {
	GetAll() []User
	FindByUsername(username string) (*User, bool)
	GetById(id int) (User, error)
	Ping()
}
