package auth

import (
	"errors"

	"github.com/marcosalbano/platform-api/internal/models"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

const dummyPasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

type BasicAuthenticator struct {
	repository UserRepository
}

func NewBasicAuthenticator(repo UserRepository) *BasicAuthenticator {
	return &BasicAuthenticator{
		repository: repo,
	}
}

func (a *BasicAuthenticator) Authenticate(
	username string,
	password string,
) (*models.User, error) {
	if a == nil || a.repository == nil {
		return nil, errors.New("authenticate: repository is not configured")
	}

	user, err := a.repository.GetByUsername(username)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Equalize lookup failures with password failures to avoid user enumeration by timing.
			_ = ComparePassword(dummyPasswordHash, password)
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := ComparePassword(
		user.PasswordHash,
		password,
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
