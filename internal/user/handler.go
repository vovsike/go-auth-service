package user

import (
	"awesomeProject/internal/apperror"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	Service Service
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	nu := &CreationDTO{}
	err := json.NewDecoder(r.Body).Decode(nu)
	if err != nil {
		return err
	}
	user, err := h.Service.CreateNewUser(nu.Name, nu.Email, nu.Password)
	if err != nil {
		return apperror.BadRequest(err)
	}
	dto := DTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Activated: user.Activated,
		Joined:    user.Joined,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(dto)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest(err)
	}
	id = parsedId.String()
	u, err := h.Service.GetUserByID(parsedId)
	if err != nil {
		return apperror.NotFound(err)
	}
	dto := DTO{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Activated: u.Activated,
		Joined:    u.Joined,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dto)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) SearchUser(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	if (name == "") && (email == "") {
		return apperror.NewHTTPError(errors.New("name or email must be provided"), http.StatusBadRequest)
	}
	users := make([]User, 0)
	if name != "" {
		if u, err := h.Service.GetUserByName(name); err == nil {
			users = append(users, *u)
		}
	}
	if email != "" {
		if u, err := h.Service.GetUserByEmail(email); err == nil {
			users = append(users, *u)
		}
	}
	if len(users) == 0 {
		return apperror.NotFound(errors.New("user not found"))
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) error {
	pw := &PasswordWrapper{}
	err := json.NewDecoder(r.Body).Decode(pw)
	if err != nil {
		return err
	}
	if pw.Email == "" || pw.Password == "" {
		return apperror.NewHTTPError(errors.New("email and password must be provided"), http.StatusBadRequest)
	}
	token, err := h.Service.Authenticate(pw.Email, pw.Password)
	if err != nil {
		return apperror.Unauthorized(err)
	}
	tw := TokenWrapper{
		Token: token,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tw)
	if err != nil {
		return err
	}
	return nil
}

type SearchDTO struct {
	Name  string
	Email string
	ID    string
}
