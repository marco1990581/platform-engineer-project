package models

import "encoding/json"

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

// UnmarshalJSON supports user stores written before the password_hash schema
// was introduced.
func (u *User) UnmarshalJSON(data []byte) error {
	var payload struct {
		Username           string `json:"username"`
		PasswordHash       string `json:"password_hash"`
		LegacyPasswordHash string `json:"PasswordHash"`
		Role               string `json:"role"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	if payload.PasswordHash == "" {
		payload.PasswordHash = payload.LegacyPasswordHash
	}

	*u = User{
		Username:     payload.Username,
		PasswordHash: payload.PasswordHash,
		Role:         payload.Role,
	}

	return nil
}
