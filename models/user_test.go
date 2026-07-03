package models

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "my-secret-password"
	hash, err := HashPassword(password)
	
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if hash == "" {
		t.Fatal("Expected hash, got empty string")
	}

	if !CheckPassword(hash, password) {
		t.Error("Password check failed for correct password")
	}

	if CheckPassword(hash, "wrong-password") {
		t.Error("Password check succeeded for wrong password")
	}
}
