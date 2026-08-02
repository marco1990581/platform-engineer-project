package auth

import (
	"sync"

	"github.com/marcosalbano/platform-api/internal/models"
)

type MemoryRepository struct {
	users map[string]models.User
	mu    sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {

	return &MemoryRepository{

		users: make(map[string]models.User),
	}

}

func (m *MemoryRepository) GetByUsername(username string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[username]

	if !ok {

		return nil, ErrUserNotFound

	}

	return &user, nil

}

func (m *MemoryRepository) Add(user models.User) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.users[user.Username] = user

}
