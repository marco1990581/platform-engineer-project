package auth

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/marcosalbano/platform-api/internal/models"
)

func TestFileRepositoryAddUserPersistsCanonicalStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.json")
	repository := NewFileRepository(path)
	user := models.User{
		Username:     "alice",
		PasswordHash: "hash",
		Role:         "admin",
	}

	if err := repository.AddUser(user); err != nil {
		t.Fatalf("add user: %v", err)
	}

	got, err := repository.GetByUsername("alice")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if *got != user {
		t.Errorf("get user = %#v, want %#v", *got, user)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat user store: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("user store permissions = %o, want 600", info.Mode().Perm())
	}
}

func TestFileRepositoryAddUserIsSafeAcrossRepositories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.json")
	const userCount = 20

	var group sync.WaitGroup
	for i := 0; i < userCount; i++ {
		group.Add(1)
		go func(index int) {
			defer group.Done()
			repository := NewFileRepository(path)
			user := models.User{
				Username:     string(rune('a' + index)),
				PasswordHash: "hash",
				Role:         "admin",
			}
			if err := repository.AddUser(user); err != nil {
				t.Errorf("add user %q: %v", user.Username, err)
			}
		}(i)
	}
	group.Wait()

	users, err := NewFileRepository(path).LoadUsers()
	if err != nil {
		t.Fatalf("load users: %v", err)
	}
	if len(users) != userCount {
		t.Errorf("stored users = %d, want %d", len(users), userCount)
	}
}
