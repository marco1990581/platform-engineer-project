package models

import (
	"encoding/json"
	"testing"
)

func TestUserJSONSchema(t *testing.T) {
	user := User{
		Username:     "alice",
		PasswordHash: "hash",
		Role:         "admin",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("marshal user: %v", err)
	}

	const want = `{"username":"alice","password_hash":"hash","role":"admin"}`
	if string(data) != want {
		t.Errorf("marshal user = %s, want %s", data, want)
	}
}

func TestUserUnmarshalLegacyPasswordHash(t *testing.T) {
	var user User
	if err := json.Unmarshal([]byte(`{"Username":"alice","PasswordHash":"hash","Role":"admin"}`), &user); err != nil {
		t.Fatalf("unmarshal legacy user: %v", err)
	}

	if user != (User{Username: "alice", PasswordHash: "hash", Role: "admin"}) {
		t.Errorf("unmarshal legacy user = %#v", user)
	}
}
