package middleware

import (
	"net/http"

	"github.com/marcosalbano/platform-api/internal/authz"
	"github.com/marcosalbano/platform-api/internal/models"
)

// Authorize returns middleware that enforces role-based access to resource
// for action.
//
// It must run after AuthMiddleware.Handler, since it reads the authenticated
// user from the request context populated there.
func Authorize(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*models.User)
			if !ok || user == nil {
				writeForbidden(w)
				return
			}

			if !authz.Enforce(user.Role, resource, action) {
				writeForbidden(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeForbidden(w http.ResponseWriter) {
	http.Error(w, "authorization failed", http.StatusForbidden)
}
