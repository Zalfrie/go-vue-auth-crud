package tests

import (
	"testing"

	"go-vue-auth-crud/models"
)

func TestPasswordHashingAndCheck(t *testing.T) {
	user := &models.User{Password: "Secret123"}
	// simulate GORM hook
	if err := user.BeforeSave(nil); err != nil {
		t.Fatalf("Hashing failed: %v", err)
	}
	// Password field should not equal plain text
	if user.Password == "Secret123" {
		t.Error("Password was not hashed")
	}
	// CheckPassword
	if err := user.CheckPassword("Secret123"); err != nil {
		t.Errorf("Expected password to match, got error: %v", err)
	}
	if err := user.CheckPassword("Wrong"); err == nil {
		t.Error("Expected mismatch error, got nil")
	}
}