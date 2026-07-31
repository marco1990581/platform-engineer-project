package auth

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/marcosalbano/platform-api/internal/models"
)

// FileRepository almacena usuarios en un archivo JSON.
//
// Más adelante podremos reemplazar esta implementación
// por SQLite, LDAP o Keycloak sin modificar el resto
// de la aplicación.
type FileRepository struct {
	path string
}

// NewFileRepository crea un nuevo repositorio basado
// en un archivo JSON.
func NewFileRepository(path string) *FileRepository {

	return &FileRepository{
		path: path,
	}

}

// LoadUsers carga todos los usuarios desde disco.
func (r *FileRepository) LoadUsers() ([]models.User, error) {

	data, err := os.ReadFile(r.path)

	if err != nil {

		if os.IsNotExist(err) {

			return []models.User{}, nil

		}

		return nil, err

	}

	var users []models.User

	err = json.Unmarshal(data, &users)

	if err != nil {

		return nil, err

	}

	return users, nil

}

// SaveUsers guarda todos los usuarios.
func (r *FileRepository) SaveUsers(users []models.User) error {

	data, err := json.MarshalIndent(
		users,
		"",
		"    ",
	)

	if err != nil {

		return err

	}

	return os.WriteFile(
		r.path,
		data,
		0600,
	)

}

// AddUser agrega un usuario.
func (r *FileRepository) AddUser(user models.User) error {

	users, err := r.LoadUsers()

	if err != nil {

		return err

	}

	for _, existing := range users {

		if existing.Username == user.Username {

			return errors.New("user already exists")

		}

	}

	users = append(users, user)

	return r.SaveUsers(users)

}
