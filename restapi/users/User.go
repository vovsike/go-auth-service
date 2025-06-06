package users

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
