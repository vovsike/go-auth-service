package users

import (
	"encoding/json"
	"net/http"
)

type Controller struct {
	Service *UserService
}

func NewController(service *UserService) *Controller {
	return &Controller{Service: service}
}

func (c *Controller) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := c.Service.GetAllUsers()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	type upWrapper struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var up upWrapper
	_ = json.NewDecoder(r.Body).Decode(&up)

	createddUser, err := c.Service.AddUser(up.Username, up.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createddUser)
}
