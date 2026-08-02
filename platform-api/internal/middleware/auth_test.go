package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcosalbano/platform-api/internal/auth"
	"github.com/marcosalbano/platform-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthMiddlewareChallengesEveryUnauthorizedRequest(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	repository := auth.NewMemoryRepository()
	repository.Add(models.User{Username: "alice", PasswordHash: string(hash)})
	middleware := NewAuthMiddleware(auth.NewBasicAuthenticator(repository))
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	tests := []struct {
		name          string
		authorization string
		wantStatus    int
	}{
		{name: "missing header", wantStatus: http.StatusUnauthorized},
		{name: "wrong scheme", authorization: "Bearer token", wantStatus: http.StatusUnauthorized},
		{name: "invalid encoding", authorization: "Basic !", wantStatus: http.StatusUnauthorized},
		{name: "wrong credentials", authorization: basicAuthorization("alice", "wrong"), wantStatus: http.StatusUnauthorized},
		{name: "valid credentials", authorization: basicAuthorization("alice", "password"), wantStatus: http.StatusNoContent},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.authorization != "" {
				request.Header.Set("Authorization", test.authorization)
			}
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Errorf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if test.wantStatus == http.StatusUnauthorized && response.Header().Get("WWW-Authenticate") == "" {
				t.Error("missing WWW-Authenticate header")
			}
		})
	}
}

func basicAuthorization(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}
