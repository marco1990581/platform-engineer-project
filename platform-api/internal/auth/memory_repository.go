package auth

import (
	"errors"

	"github.com/marcosalbano/platform-api/internal/models"
)

type MemoryRepository struct {
	users map[string]models.User
}

func NewMemoryRepository() *MemoryRepository {

	return &MemoryRepository{

		users: make(map[string]models.User),
	}

}

func (m *MemoryRepository) GetByUsername(username string) (*models.User, error) {

	user, ok := m.users[username]

	if !ok {

		return nil, errors.New("user not found")

	}

	return &user, nil

}

func (m *MemoryRepository) Add(user models.User) {

	m.users[user.Username] = user

}
