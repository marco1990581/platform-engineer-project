package auth

import (
	"sync"
	"testing"

	"github.com/marcosalbano/platform-api/internal/models"
)

func TestMemoryRepositoryConcurrentAccess(t *testing.T) {
	repository := NewMemoryRepository()
	const userCount = 100

	var group sync.WaitGroup
	for i := 0; i < userCount; i++ {
		group.Add(1)
		go func(index int) {
			defer group.Done()
			username := string(rune('a' + index))
			repository.Add(models.User{Username: username})
			if _, err := repository.GetByUsername(username); err != nil {
				t.Errorf("get user %q: %v", username, err)
			}
		}(i)
	}
	group.Wait()
}
