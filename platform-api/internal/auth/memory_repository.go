package auth


import (

	"errors"

	"platform-api/internal/models"


)

type MemoryRepository struct {

	users map[string]models.User
}

func NewMemoryRepository() *MemoryRepository {

	return &MemoryRepository{

		users: map[string]models.User{};
	}

}

func (m *MemoryRepository) GetByUserName(username string) (*models.User, error) {
	
	user, ok := m.users[username]


	if !ok {

		return nil, errors.New("user not found")


	}

	return &user, nil


}


func (m *MemoryRepository) Add(user models.User) {

	m.users[user.Username] = user

}


