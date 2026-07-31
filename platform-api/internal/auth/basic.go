package auth

import (
	"errors"

	"github.com/marcosalbano/platform-api/internal/models"
)

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

	user, err := a.repository.GetByUsername(username)

	if err != nil {
		return nil, err
	}

	err = ComparePassword(
		user.PasswordHash,
		password,
	)

	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
