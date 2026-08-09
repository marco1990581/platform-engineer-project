package config

import (
	"os"
	"strings"
)

const (
	UsersFilePathEnv     = "PLATFORM_API_USERS_FILE"
	DefaultUsersFilePath = "configs/users.json"
)

// UsersFilePath returns the configured user store or the local development default.
func UsersFilePath() string {
	if path := strings.TrimSpace(os.Getenv(UsersFilePathEnv)); path != "" {
		return path
	}

	return DefaultUsersFilePath
}
