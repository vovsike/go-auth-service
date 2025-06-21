package accounts

import (
	"encoding/json"
	"log"
	"net/http"
	"restapi/internal/jwtInternal"
	"restapi/utils"
)

type Controller struct {
	Service    AccountService
	JwtService jwtInternal.Service
}

func NewController(service AccountService, jwtService jwtInternal.Service) *Controller {
	return &Controller{Service: service, JwtService: jwtService}
}

func (c *Controller) CreateNewAccount(w http.ResponseWriter, r *http.Request) {
	type upWrapper struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var up upWrapper
	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	account, err := c.Service.CreateNewAccount(r.Context(), up.Username, up.Password)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to create account")
		return
	}

	utils.RespondJSON(w, account, http.StatusCreated)
	return
}
