package auth

import (
	"errors"
	"testing"

	"github.com/marcosalbano/platform-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func TestBasicAuthenticatorAuthenticate(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	repository := NewMemoryRepository()
	repository.Add(models.User{Username: "alice", PasswordHash: string(hash)})
	authenticator := NewBasicAuthenticator(repository)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{name: "valid credentials", username: "alice", password: "correct-password"},
		{name: "invalid password", username: "alice", password: "incorrect", wantErr: ErrInvalidCredentials},
		{name: "unknown user", username: "unknown", password: "incorrect", wantErr: ErrInvalidCredentials},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := authenticator.Authenticate(test.username, test.password)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Authenticate() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && user.Username != "alice" {
				t.Errorf("Authenticate() user = %#v, want alice", user)
			}
		})
	}
}
