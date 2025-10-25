package user

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestInMemoryService_Authenticate(t *testing.T) {
	t.Setenv("SIGN_KEY", "secret")
	validEmail := "valid@email.test"
	invalidEmail := "invalid@email"
	validPassword := "validPassword"
	invalidPassword := "invalidPassword"

	validUser, err := NewUser("valid", validEmail, validPassword)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		users Store
	}
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{name: "valid user",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						validEmail: validUser,
					},
				},
			},
			args:    args{email: validEmail, password: validPassword},
			want:    "token",
			wantErr: false,
		},
		{name: "invalid user",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						validEmail: validUser,
					},
				},
			},
			args:    args{email: validEmail, password: invalidPassword},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid email",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						validEmail: validUser,
					},
				},
			},
			args:    args{email: invalidEmail, password: validPassword},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty email",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						validEmail: validUser,
					},
				},
			},
			args:    args{email: "", password: validPassword},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty password",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						validEmail: validUser,
					},
				},
			},
			args:    args{email: validEmail, password: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &InMemoryService{
				users: tt.fields.users,
			}
			got, err := us.Authenticate(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == "") != tt.wantErr {
				t.Errorf("Authenticate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryService_CreateNewUser(t *testing.T) {

	validEmail := "valid@email.test"
	invalidEmail := "invalid.email"

	type fields struct {
		users Store
	}
	type args struct {
		name     string
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "valid user",
			fields: fields{
				users: &InMemStore{
					usersByName:  map[string]*User{},
					usersByID:    map[uuid.UUID]*User{},
					usersByEmail: map[string]*User{},
				},
			},
			args: args{
				name:     "valid",
				email:    validEmail,
				password: "password",
			},
			want: &User{
				Name:  "valid",
				Email: validEmail,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			fields: fields{
				users: &InMemStore{
					usersByName:  map[string]*User{},
					usersByID:    map[uuid.UUID]*User{},
					usersByEmail: map[string]*User{},
				},
			},
			args: args{
				name:     "valid",
				email:    invalidEmail,
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty email",
			fields: fields{
				users: &InMemStore{
					usersByName:  map[string]*User{},
					usersByID:    map[uuid.UUID]*User{},
					usersByEmail: map[string]*User{},
				},
			},
			args: args{
				name:     "valid",
				email:    "",
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty password",
			fields: fields{
				users: &InMemStore{
					usersByName:  map[string]*User{},
					usersByID:    map[uuid.UUID]*User{},
					usersByEmail: map[string]*User{},
				},
			},
			args: args{
				name:     "valid",
				email:    validEmail,
				password: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "password made of spaces",
			fields: fields{
				users: &InMemStore{
					usersByName:  map[string]*User{},
					usersByID:    map[uuid.UUID]*User{},
					usersByEmail: map[string]*User{},
				},
			},
			args: args{
				name:     "valid",
				email:    validEmail,
				password: "          ",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "user already exists",
			fields: fields{
				users: &InMemStore{
					usersByName: map[string]*User{
						"valid": {
							Name:  "valid",
							Email: validEmail,
						},
					},
					usersByID: map[uuid.UUID]*User{},
					usersByEmail: map[string]*User{
						validEmail: {
							Name:  "valid",
							Email: validEmail,
						},
					},
				},
			},
			args: args{
				name:     "valid",
				email:    validEmail,
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &InMemoryService{
				users: tt.fields.users,
			}
			got, err := us.CreateNewUser(tt.args.name, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && (got == tt.want) {
				return
			}
			if got.Name != tt.want.Name {
				t.Errorf("CreateNewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryService_GetUserByEmail(t *testing.T) {

	existingEmail := "existing@email.test"
	nonExistingEmail := "fake@email.test"
	invalidEmail := "invalid@email"
	type fields struct {
		users Store
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "user exists",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						existingEmail: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				email: existingEmail,
			},
			want: &User{
				Name: "exists",
			},
			wantErr: false,
		},
		{
			name: "user does not exists",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						existingEmail: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				email: nonExistingEmail,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty email",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						existingEmail: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				email: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid email",
			fields: fields{
				users: &InMemStore{
					usersByEmail: map[string]*User{
						existingEmail: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				email: invalidEmail,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &InMemoryService{
				users: tt.fields.users,
			}
			got, err := us.GetUserByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryService_GetUserByID(t *testing.T) {
	existingID := uuid.New()
	nonExistingID := uuid.New()
	type fields struct {
		users Store
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "user exists",
			fields: fields{
				users: &InMemStore{
					usersByID: map[uuid.UUID]*User{
						existingID: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				id: existingID,
			},
			want: &User{
				Name: "exists",
			},
			wantErr: false,
		},
		{
			name: "user does not exists",
			fields: fields{
				users: &InMemStore{
					usersByID: map[uuid.UUID]*User{
						existingID: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				id: nonExistingID,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty id",
			fields: fields{
				users: &InMemStore{
					usersByID: map[uuid.UUID]*User{
						existingID: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				id: uuid.UUID{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &InMemoryService{
				users: tt.fields.users,
			}
			got, err := us.GetUserByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryService_GetUserByName(t *testing.T) {
	existingName := "existing"
	nonExistingName := "fake"
	type fields struct {
		users Store
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "user exists",
			fields: fields{
				users: &InMemStore{
					usersByName: map[string]*User{
						existingName: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				name: existingName,
			},
			want: &User{
				Name: "exists",
			},
			wantErr: false,
		},
		{
			name: "user does not exists",
			fields: fields{
				users: &InMemStore{
					usersByName: map[string]*User{
						existingName: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				name: nonExistingName,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty name",
			fields: fields{
				users: &InMemStore{
					usersByName: map[string]*User{
						existingName: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				name: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "name made of spaces",
			fields: fields{
				users: &InMemStore{
					usersByName: map[string]*User{
						existingName: {
							Name: "exists",
						},
					},
				},
			},
			args: args{
				name: "          ",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &InMemoryService{
				users: tt.fields.users,
			}
			got, err := us.GetUserByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInMemoryUserService(t *testing.T) {
	type args struct {
		users Store
	}
	tests := []struct {
		name string
		args args
		want *InMemoryService
	}{
		{
			name: "valid users",
			args: args{
				users: &InMemStore{
					usersByName:  map[string]*User{},
					usersByID:    map[uuid.UUID]*User{},
					usersByEmail: map[string]*User{},
				},
			},
			want: &InMemoryService{
				users: &InMemStore{
					usersByName:  map[string]*User{},
					usersByID:    map[uuid.UUID]*User{},
					usersByEmail: map[string]*User{},
				},
			},
		},
		{
			name: "nil users",
			args: args{
				users: nil,
			},
			want: &InMemoryService{
				users: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInMemoryUserService(tt.args.users); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemoryUserService() = %v, want %v", got, tt.want)
			}
			if got := NewInMemoryUserService(tt.args.users); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemoryUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_issueSignedToken(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		setEnv  bool
	}{
		{
			name: "valid user",
			args: args{
				user: &User{
					ID:   uuid.New(),
					Name: "valid",
				},
			},
			want:    "not empty	",
			wantErr: false,
			setEnv:  true,
		},
		{
			name: "nil user",
			args: args{
				user: nil,
			},
			want:    "",
			wantErr: true,
			setEnv:  true,
		},
		{
			name: "no secret",
			args: args{
				user: &User{
					ID:   uuid.New(),
					Name: "valid",
				},
			},
			want:    "",
			wantErr: true,
			setEnv:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv("SIGN_KEY", "not empty")
			}
			got, err := issueSignedToken(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("issueSignedToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) == 0 && !tt.wantErr {
				t.Errorf("issueSignedToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
