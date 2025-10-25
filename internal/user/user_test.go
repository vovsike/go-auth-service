package user

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	validName := "valid"
	validEmail := "valid@email.test"
	invalidEmail := "invalid.email"
	validPassword := "password"
	invalidPassword := ""
	type args struct {
		name     string
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "valid user",
			args: args{
				name:     validName,
				email:    validEmail,
				password: validPassword,
			},
			want: &User{
				Name:  validName,
				Email: validEmail,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			args: args{
				name:     validName,
				email:    invalidEmail,
				password: validPassword,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid password",
			args: args{
				name:     validName,
				email:    validEmail,
				password: invalidPassword,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr) && (got.Name != tt.want.Name || got.Email != tt.want.Email) {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Activate(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		Name      string
		Email     string
		hash      []byte
		Joined    time.Time
		Activated bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "activate user",
			fields: fields{
				Activated: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Email:     tt.fields.Email,
				hash:      tt.fields.hash,
				Joined:    tt.fields.Joined,
				Activated: tt.fields.Activated,
			}
			u.Activate()
			if u.Activated != true {
				t.Errorf("Activate() = %v, want %v", u.Activated, true)
			}
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	validPassword := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
	invalidPassword := "fake"
	type fields struct {
		ID        uuid.UUID
		Name      string
		Email     string
		hash      []byte
		Joined    time.Time
		Activated bool
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "valid password",
			fields: fields{
				hash: hashedPassword,
			},
			args: args{
				password: "password",
			},
			want: true,
		},
		{
			name: "invalid password",
			fields: fields{
				hash: hashedPassword,
			},
			args: args{
				password: invalidPassword,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Email:     tt.fields.Email,
				hash:      tt.fields.hash,
				Joined:    tt.fields.Joined,
				Activated: tt.fields.Activated,
			}
			if got := u.CheckPassword(tt.args.password); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPasswordValid(t *testing.T) {
	validPassword := "password"
	passwordTooShort := "pw"
	passwordContainsSpace := "password with space"
	emptyPassword := ""
	spacePassword := " "
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid password",
			args: args{
				password: validPassword,
			},
			want: true,
		},
		{
			name: "password too short",
			args: args{
				password: passwordTooShort,
			},
			want: false,
		},
		{
			name: "password contains space",
			args: args{
				password: passwordContainsSpace,
			},
			want: false,
		},
		{
			name: "empty password",
			args: args{
				password: emptyPassword,
			},
			want: false,
		},
		{
			name: "space password",
			args: args{
				password: spacePassword,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPasswordValid(tt.args.password); got != tt.want {
				t.Errorf("isPasswordValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidEmail(t *testing.T) {
	validEmail := "valid@email.test"
	invalidEmail := "invalid.email"
	shortEmail := "e@"
	longemail := strings.Repeat("a", 256) + "@email.test"

	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid email",
			args: args{
				email: validEmail,
			},
			want: true,
		},
		{
			name: "invalid email",
			args: args{
				email: invalidEmail,
			},
			want: false,
		},
		{
			name: "short email",
			args: args{
				email: shortEmail,
			},
			want: false,
		},
		{
			name: "long email",
			args: args{
				email: longemail,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidEmail(tt.args.email); got != tt.want {
				t.Errorf("isValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
