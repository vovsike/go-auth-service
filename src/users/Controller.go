package users

import (
	"encoding/json"
	"net/http"
	"restapi/utils"
)

type Controller struct {
	Service *UserService
}

func NewController(service *UserService) *Controller {
	return &Controller{Service: service}
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type upWrapper struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var up upWrapper
	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	if up.Username == "" || up.Password == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Username or password is empty")
		return
	}

	createdUser, err := c.Service.CreateNewUser(ctx, up.Username, up.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.RespondJSON(w, createdUser, http.StatusCreated)
}
