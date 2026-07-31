package auth

import "github.com/marcosalbano/platform-api/internal/models"

type Authenticator interface {
	Authenticate(username, password string) (*models.User, error)
}
