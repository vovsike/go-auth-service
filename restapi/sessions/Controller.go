package sessions

import (
	"encoding/json"
	"net/http"
	"restapi/jwtInternal"
	"restapi/users"
	"restapi/utils"
)

type Controller struct {
	Service      *SessionService
	UsersService *users.UserService
	JwtService   jwtInternal.Service
}

func NewController(service *SessionService, usersService *users.UserService, jwtService jwtInternal.Service) *Controller {
	return &Controller{Service: service, UsersService: usersService, JwtService: jwtService}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	type upWrapper struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var up upWrapper

	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to decode request body")
	}

	_, err = c.UsersService.CheckUserPassword(up.Username, up.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid username or password")
	}
	un, _ := c.UsersService.Store.FindByUsername(up.Username)
	s := c.Service.Authenticate(un.ID)
	utils.RespondJSON(w, s, http.StatusCreated)
}

func (c *Controller) GetToken(w http.ResponseWriter, r *http.Request) {
	type sWrap struct {
		SessionId string `json:"sessionId"`
	}

	var s sWrap

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to decode request body")
	}
	sesh, ok := c.Service.VerifySession(s.SessionId)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid session")
	}

	jwtToken, err := c.JwtService.GenerateToken(sesh.UserID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
	}

	type jwtWrapper struct {
		Token []byte `json:"token"`
	}
	jwtWrapped := jwtWrapper{Token: jwtToken}

	utils.RespondJSON(w, jwtWrapped, http.StatusOK)

}
