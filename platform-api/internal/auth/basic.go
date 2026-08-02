package auth

import (
	"errors"
	"log"

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
		log.Printf("auth: authentication failed during user lookup username=%q: %v", username, err)

		return nil, err
	}

	err = ComparePassword(
		user.PasswordHash,
		password,
	)

	if err != nil {
		log.Printf("auth: bcrypt password comparison failed username=%q", username)

		return nil, errors.New("invalid password")
	}

	log.Printf("auth: bcrypt password comparison succeeded username=%q", username)

	return user, nil
}
