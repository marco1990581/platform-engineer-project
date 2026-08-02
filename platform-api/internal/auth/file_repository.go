package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/marcosalbano/platform-api/internal/models"
)

// FileRepository almacena usuarios en un archivo JSON.
//
// Más adelante podremos reemplazar esta implementación
// por SQLite, LDAP o Keycloak sin modificar el resto
// de la aplicación.
type FileRepository struct {
	path string
	mu   sync.Mutex
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
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.withExclusiveLock(func() error {
		return r.saveUsers(users)
	})
}

func (r *FileRepository) saveUsers(users []models.User) error {

	data, err := json.MarshalIndent(
		users,
		"",
		"    ",
	)
	if err != nil {
		return err
	}

	dir := filepath.Dir(r.path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create user store directory: %w", err)
	}

	tempFile, err := os.CreateTemp(dir, ".users-*.json")
	if err != nil {
		return fmt.Errorf("create temporary user store: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if err := tempFile.Chmod(0600); err != nil {
		tempFile.Close()
		return fmt.Errorf("set temporary user store permissions: %w", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("write temporary user store: %w", err)
	}
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return fmt.Errorf("sync temporary user store: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temporary user store: %w", err)
	}

	if err := os.Rename(tempPath, r.path); err != nil {
		return fmt.Errorf("replace user store: %w", err)
	}

	directory, err := os.Open(dir)
	if err != nil {
		return fmt.Errorf("open user store directory: %w", err)
	}
	defer directory.Close()

	if err := directory.Sync(); err != nil {
		return fmt.Errorf("sync user store directory: %w", err)
	}

	return nil

}

// AddUser agrega un usuario.
func (r *FileRepository) AddUser(user models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.withExclusiveLock(func() error {
		users, err := r.LoadUsers()

		if err != nil {

			return err

		}

		for _, existing := range users {

			if existing.Username == user.Username {

				return fmt.Errorf("user %q already exists", user.Username)

			}
		}

		users = append(users, user)

		return r.saveUsers(users)
	})

}

func (r *FileRepository) withExclusiveLock(operation func() error) error {
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create user store directory: %w", err)
	}

	lockFile, err := os.OpenFile(r.path+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("open user store lock: %w", err)
	}
	defer lockFile.Close()

	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("lock user store: %w", err)
	}
	defer syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)

	return operation()
}

func (r *FileRepository) GetByUsername(username string) (*models.User, error) {

	users, err := r.LoadUsers()
	if err != nil {
		log.Printf("auth: user lookup failed: %v", err)

		return nil, err
	}

	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}

	log.Print("auth: user lookup found no match")

	return nil, ErrUserNotFound

}
