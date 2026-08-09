package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcosalbano/platform-api/internal/models"
)

func TestAuthorizeAllowsPermittedRole(t *testing.T) {
	handler := Authorize("system", "read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := withUser(httptest.NewRequest(http.MethodGet, "/api/v1/system", nil), &models.User{Username: "alice", Role: "viewer"})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestAuthorizeDeniesUnpermittedRole(t *testing.T) {
	handler := Authorize("users", "read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := withUser(httptest.NewRequest(http.MethodGet, "/api/v1/users", nil), &models.User{Username: "alice", Role: "viewer"})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}

func TestAuthorizeDeniesMissingUser(t *testing.T) {
	handler := Authorize("system", "read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/api/v1/system", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}

func withUser(r *http.Request, user *models.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), UserContextKey, user))
}
