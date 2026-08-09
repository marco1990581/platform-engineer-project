package main

import "testing"

func TestClearPassword(t *testing.T) {
	password := []byte("password")

	clearPassword(password)

	for _, value := range password {
		if value != 0 {
			t.Errorf("password byte = %d, want 0", value)
		}
	}
}
