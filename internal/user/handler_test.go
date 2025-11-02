package user

import (
	"awesomeProject/internal/apperror"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_Authenticate(t *testing.T) {
	t.Setenv("SIGN_KEY", "secret")

	service := createTestService(t, "valid", "password", "valid@email.test")
	h := &Handler{Service: service}

	tests := []struct {
		name      string
		body      PasswordWrapper
		wantError *apperror.HTTPError
	}{
		{
			name:      "valid user",
			body:      PasswordWrapper{Password: "password", Email: "valid@email.test"},
			wantError: nil,
		},
		{
			name:      "wrong password",
			body:      PasswordWrapper{Password: "wrong", Email: "valid@email.test"},
			wantError: &apperror.HTTPError{StatusCode: http.StatusUnauthorized},
		},
		{
			name:      "no user found",
			body:      PasswordWrapper{Password: "password", Email: "invalid@email.test"},
			wantError: &apperror.HTTPError{StatusCode: http.StatusUnauthorized},
		},
		{
			name:      "empty password",
			body:      PasswordWrapper{Password: "", Email: "valid@email.test"},
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
		{
			name:      "empty email",
			body:      PasswordWrapper{Password: "password", Email: ""},
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewReader(body))
			w := httptest.NewRecorder()

			err = h.Authenticate(w, req)

			if tt.wantError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				var httpError *apperror.HTTPError
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, tt.wantError.StatusCode, httpError.StatusCode)
			}
		})
	}
}

func TestHandler_CreateUser(t *testing.T) {
	t.Setenv("SIGN_KEY", "secret")
	validName := "valid"
	validEmail := "valid@email.test"
	validPassword := "password"
	existingName := "exists"
	existingEmail := "existing@email.test"
	invalidPassword := "pw"

	service := createTestService(t, existingName, validPassword, existingEmail)
	h := &Handler{Service: service}

	tests := []struct {
		name      string
		body      CreationDTO
		wantError *apperror.HTTPError
	}{
		{
			name: "valid user",
			body: CreationDTO{
				Name:     validName,
				Email:    validEmail,
				Password: validPassword,
			},
			wantError: nil,
		},
		{
			name: "invalid password",
			body: CreationDTO{
				Name:     validName,
				Email:    validEmail,
				Password: invalidPassword,
			},
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
		{
			name: "user already exists",
			body: CreationDTO{
				Name:     existingName,
				Email:    existingEmail,
				Password: validPassword,
			},
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
		{
			name: "empty password",
			body: CreationDTO{
				Name:     validName,
				Email:    validEmail,
				Password: "",
			},
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
		{
			name: "empty email",
			body: CreationDTO{
				Name:     validName,
				Email:    "",
				Password: validPassword,
			},
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
			w := httptest.NewRecorder()

			err = h.CreateUser(w, req)
			if tt.wantError == nil {
				assert.NoError(t, err)
				var dto DTO
				err = json.NewDecoder(w.Body).Decode(&dto)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotEmpty(t, dto.ID)
				assert.Equal(t, dto.Name, tt.body.Name)
				assert.Equal(t, dto.Email, tt.body.Email)
			} else {
				assert.Error(t, err)
				var httpError *apperror.HTTPError
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, tt.wantError.StatusCode, httpError.StatusCode)
			}
		})
	}
}

func TestHandler_GetUser(t *testing.T) {
	validID := uuid.New()
	invalidID := uuid.New()
	emptyID := ""
	spaceID := "   "

	service := createTestService(t, "valid", "password", "valid@email.test", validID)
	h := &Handler{Service: service}

	tests := []struct {
		name      string
		id        string
		wantError *apperror.HTTPError
	}{
		{
			name:      "valid id",
			id:        validID.String(),
			wantError: nil,
		},
		{
			name:      "invalid id",
			id:        invalidID.String(),
			wantError: &apperror.HTTPError{StatusCode: http.StatusNotFound},
		},
		{
			name:      "empty id",
			id:        emptyID,
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
		{
			name:      "id with spaces",
			id:        spaceID,
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user/", nil)
			w := httptest.NewRecorder()

			chiCtx := &chi.Context{
				URLParams: chi.RouteParams{
					Keys:   []string{"id"},
					Values: []string{tt.id},
				},
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			err := h.GetUser(w, req)
			if tt.wantError == nil {
				assert.NoError(t, err, "got error")
				var dto DTO
				err := json.NewDecoder(w.Body).Decode(&dto)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotEmpty(t, dto.ID)
				assert.Equal(t, dto.ID, uuid.MustParse(tt.id))
			} else {
				assert.Error(t, err)
				var httpError *apperror.HTTPError
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, tt.wantError.StatusCode, httpError.StatusCode)
			}
		})
	}
}

func TestHandler_SearchUser(t *testing.T) {

	validName := "valid"
	validEmail := "valid@email.test"

	invalidName := "invalid"
	invalidEmail := "invalid@email.test"

	emptyName := ""
	emptyEmail := ""

	spaceName := "    "
	spaceEmail := "   "

	badlyFormatedEmail := "invalid@email"

	service := createTestService(t, validName, "password", validEmail)
	h := &Handler{Service: service}

	tests := []struct {
		name      string
		body      SearchDTO
		wantError *apperror.HTTPError
	}{
		{
			name:      "valid search",
			body:      SearchDTO{Name: validName, Email: validEmail},
			wantError: nil,
		},
		{
			name:      "invalid name",
			body:      SearchDTO{Name: invalidName},
			wantError: &apperror.HTTPError{StatusCode: http.StatusNotFound},
		},
		{
			name:      "invalid email",
			body:      SearchDTO{Email: invalidEmail},
			wantError: &apperror.HTTPError{StatusCode: http.StatusNotFound},
		},
		{
			name:      "empty name",
			body:      SearchDTO{Name: emptyName},
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
		{
			name:      "empty email",
			body:      SearchDTO{Email: emptyEmail},
			wantError: &apperror.HTTPError{StatusCode: http.StatusBadRequest},
		},
		{
			name:      "name with spaces",
			body:      SearchDTO{Name: spaceName},
			wantError: &apperror.HTTPError{StatusCode: http.StatusNotFound},
		},
		{
			name:      "email with spaces",
			body:      SearchDTO{Email: spaceEmail},
			wantError: &apperror.HTTPError{StatusCode: http.StatusNotFound},
		},
		{
			name:      "badly formated email",
			body:      SearchDTO{Email: badlyFormatedEmail},
			wantError: &apperror.HTTPError{StatusCode: http.StatusNotFound},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			req.URL.RawQuery = "name=" + tt.body.Name + "&email=" + tt.body.Email
			w := httptest.NewRecorder()

			err := h.SearchUser(w, req)
			if tt.wantError == nil {
				assert.NoError(t, err, "got error")
				var dto []DTO
				err := json.NewDecoder(w.Body).Decode(&dto)
				if err != nil {
					t.Fatal(err)
				}
				for _, d := range dto {
					assert.NotEmpty(t, d.ID)
					assert.Equal(t, d.Name, tt.body.Name)
					assert.Equal(t, d.Email, tt.body.Email)
				}
			} else {
				assert.Error(t, err)
				var httpError *apperror.HTTPError
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, tt.wantError.StatusCode, httpError.StatusCode)
			}
		})
	}
}

func createTestService(t *testing.T, username, password, email string, id ...uuid.UUID) Service {
	t.Helper()
	var userID uuid.UUID
	if len(id) > 0 {
		userID = id[0]
	} else {
		userID = uuid.New()
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	return &InMemoryService{users: &InMemStore{
		usersByID: map[uuid.UUID]*User{
			userID: {
				ID:    userID,
				Name:  username,
				Email: email,
				hash:  hash,
			},
		},
		usersByEmail: map[string]*User{
			email: {
				ID:    userID,
				Name:  username,
				Email: email,
				hash:  hash,
			},
		},
		usersByName: map[string]*User{
			username: {
				ID:    userID,
				Name:  username,
				Email: email,
			},
		},
	}}
}
