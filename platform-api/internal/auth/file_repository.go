package auth

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

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

	log.Printf("auth: loading user file path=%q", r.path)

	data, err := os.ReadFile(r.path)

	if err != nil {

		if os.IsNotExist(err) {

			log.Printf("auth: user file not found path=%q", r.path)

			return []models.User{}, nil

		}

		log.Printf("auth: failed to read user file path=%q: %v", r.path, err)

		return nil, err

	}

	var users []models.User

	err = json.Unmarshal(data, &users)

	if err != nil {

		log.Printf("auth: failed to parse user file path=%q: %v", r.path, err)

		return nil, err

	}

	log.Printf("auth: loaded %d user(s) from path=%q", len(users), r.path)

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

	// Obtiene el directorio del archivo.
	dir := filepath.Dir(r.path)

	// Lo crea si no existe.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Escribe el archivo con permisos 0600.
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

func (r *FileRepository) GetByUsername(username string) (*models.User, error) {

	users, err := r.LoadUsers()
	if err != nil {
		log.Printf("auth: user lookup failed username=%q: %v", username, err)

		return nil, err
	}

	for _, user := range users {
		if user.Username == username {
			log.Printf("auth: user lookup succeeded username=%q", username)

			return &user, nil
		}
	}

	log.Printf("auth: user lookup found no match username=%q", username)

	return nil, errors.New("user not found")

}
