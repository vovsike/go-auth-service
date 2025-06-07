package sessions

import "net/http"

type Controller struct {
	Service *SessionService
}

func NewController(service *SessionService) *Controller {
	return &Controller{Service: service}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	c.Service.CreateNewSession(1)
	w.WriteHeader(http.StatusCreated)
	return
}
