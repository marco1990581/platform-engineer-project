package auth

import "github.com/marcosalbano/platform-api/internal/models"

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
}
