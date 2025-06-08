package sessions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/jwtInternal"
	"restapi/users"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = c.UsersService.CheckUserPassword(up.Username, up.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	un, _ := c.UsersService.Store.FindByUsername(up.Username)
	s := c.Service.Authenticate(un.ID)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func (c *Controller) GetToken(w http.ResponseWriter, r *http.Request) {
	type sWrap struct {
		SessionId string `json:"sessionId"`
	}

	var s sWrap

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sesh, ok := c.Service.VerifySession(s.SessionId)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtToken, err := c.JwtService.GenerateToken(sesh.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type jwtWrapper struct {
		Token []byte `json:"token"`
	}
	jwtWrapped := jwtWrapper{Token: jwtToken}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(jwtWrapped)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
