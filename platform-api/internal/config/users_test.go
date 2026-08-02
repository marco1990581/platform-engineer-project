package config

import "testing"

func TestUsersFilePathUsesDefaultWhenUnset(t *testing.T) {
	t.Setenv(UsersFilePathEnv, "")

	if got := UsersFilePath(); got != DefaultUsersFilePath {
		t.Fatalf("UsersFilePath() = %q, want %q", got, DefaultUsersFilePath)
	}
}

func TestUsersFilePathUsesEnvironmentValue(t *testing.T) {
	const configuredPath = "/run/secrets/platform-api-users.json"
	t.Setenv(UsersFilePathEnv, configuredPath)

	if got := UsersFilePath(); got != configuredPath {
		t.Fatalf("UsersFilePath() = %q, want %q", got, configuredPath)
	}
}
