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
