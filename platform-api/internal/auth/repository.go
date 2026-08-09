package auth

import (
	"errors"

	"github.com/marcosalbano/platform-api/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
}
