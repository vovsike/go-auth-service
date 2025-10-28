package user

import (
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
		name     string
		body     PasswordWrapper
		wantCode int
	}{
		{
			name:     "valid user",
			body:     PasswordWrapper{Password: "password", Email: "valid@email.test"},
			wantCode: http.StatusOK,
		},
		{
			name:     "wrong password",
			body:     PasswordWrapper{Password: "wrong", Email: "valid@email.test"},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "no user found",
			body:     PasswordWrapper{Password: "password", Email: "invalid@email.test"},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "empty password",
			body:     PasswordWrapper{Password: "", Email: "valid@email.test"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty email",
			body:     PasswordWrapper{Password: "password", Email: ""},
			wantCode: http.StatusBadRequest,
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

			assert.NoError(t, h.Authenticate(w, req))
			assert.Equal(t, tt.wantCode, w.Code)
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
		name     string
		body     CreationDTO
		wantCode int
	}{
		{
			name: "valid user",
			body: CreationDTO{
				Name:     validName,
				Email:    validEmail,
				Password: validPassword,
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "invalid password",
			body: CreationDTO{
				Name:     validName,
				Email:    validEmail,
				Password: invalidPassword,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "user already exists",
			body: CreationDTO{
				Name:     existingName,
				Email:    existingEmail,
				Password: validPassword,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "empty password",
			body: CreationDTO{
				Name:     validName,
				Email:    validEmail,
				Password: "",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "empty email",
			body: CreationDTO{
				Name:     validName,
				Email:    "",
				Password: validPassword,
			},
			wantCode: http.StatusBadRequest,
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

			assert.NoError(t, h.CreateUser(w, req))
			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode != http.StatusCreated {
				return
			}
			var dto DTO
			err = json.NewDecoder(w.Body).Decode(&dto)
			if err != nil {
				t.Fatal(err)
			}
			assert.NotEmpty(t, dto.ID)
			assert.Equal(t, dto.Name, tt.body.Name)
			assert.Equal(t, dto.Email, tt.body.Email)
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
		name     string
		id       string
		wantCode int
	}{
		{
			name:     "valid id",
			id:       validID.String(),
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid id",
			id:       invalidID.String(),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "empty id",
			id:       emptyID,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "id with spaces",
			id:       spaceID,
			wantCode: http.StatusBadRequest,
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

			assert.NoError(t, h.GetUser(w, req), "got error")
			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode != http.StatusOK {
				return
			}
			var dto DTO
			err := json.NewDecoder(w.Body).Decode(&dto)
			if err != nil {
				t.Fatal(err)
			}
			assert.NotEmpty(t, dto.ID)
			assert.Equal(t, dto.ID, uuid.MustParse(tt.id))
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
		name     string
		body     SearchDTO
		wantCode int
	}{
		{
			name:     "valid search",
			body:     SearchDTO{Name: validName, Email: validEmail},
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid name",
			body:     SearchDTO{Name: invalidName},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "invalid email",
			body:     SearchDTO{Email: invalidEmail},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "empty name",
			body:     SearchDTO{Name: emptyName},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty email",
			body:     SearchDTO{Email: emptyEmail},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "name with spaces",
			body:     SearchDTO{Name: spaceName},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "email with spaces",
			body:     SearchDTO{Email: spaceEmail},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "badly formated email",
			body:     SearchDTO{Email: badlyFormatedEmail},
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			req.URL.RawQuery = "name=" + tt.body.Name + "&email=" + tt.body.Email
			w := httptest.NewRecorder()

			assert.NoError(t, h.SearchUser(w, req), "got error")
			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode != http.StatusOK {
				return
			}
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
