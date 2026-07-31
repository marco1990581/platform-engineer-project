package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/marcosalbano/platform-api/internal/auth"
)

// contextKey evita colisiones con otras claves del contexto.
type contextKey string

const UserContextKey contextKey = "authenticatedUser"

// AuthMiddleware protege los endpoints utilizando
// HTTP Basic Authentication.
type AuthMiddleware struct {
	authenticator auth.Authenticator
}

// NewAuthMiddleware construye el middleware.
func NewAuthMiddleware(a auth.Authenticator) *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: a,
	}
}

// Handler envuelve cualquier http.Handler.
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorization := r.Header.Get("Authorization")

		if authorization == "" {

			w.Header().Set("WWW-Authenticate", `Basic realm="Platform API"`)

			http.Error(w, "missing authorization header", http.StatusUnauthorized)

			return
		}

		if !strings.HasPrefix(authorization, "Basic ") {

			http.Error(w, "invalid authentication method", http.StatusUnauthorized)

			return
		}

		encodedCredentials := strings.TrimPrefix(
			authorization,
			"Basic ",
		)

		decodedBytes, err := base64.StdEncoding.DecodeString(
			encodedCredentials,
		)

		if err != nil {

			http.Error(w, "invalid authorization header", http.StatusUnauthorized)

			return
		}

		credentials := strings.SplitN(
			string(decodedBytes),
			":",
			2,
		)

		if len(credentials) != 2 {

			http.Error(w, "invalid credentials", http.StatusUnauthorized)

			return
		}

		username := credentials[0]
		password := credentials[1]

		user, err := m.authenticator.Authenticate(
			username,
			password,
		)

		if err != nil {

			http.Error(w, "authentication failed", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserContextKey,
			user,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)

	})

}
