package sessions

import (
	"encoding/json"
	"net/http"
	"restapi/users"
)

type Controller struct {
	Service      *SessionService
	UsersService *users.UserService
}

func NewController(service *SessionService, usersService *users.UserService) *Controller {
	return &Controller{Service: service, UsersService: usersService}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	type upWrapper struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var up upWrapper

	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = c.UsersService.CheckUserPassword(up.Username, up.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	un, _ := c.UsersService.Store.FindByUsername(up.Username)
	s := c.Service.Authenticate(un.Id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}
