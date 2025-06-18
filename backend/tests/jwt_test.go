package tests

import (
	"testing"
	"time"

	"go-vue-auth-crud/models"
	"go-vue-auth-crud/services"
)

func TestGenerateToken(t *testing.T) {
	user := models.User{ID: 1, Email: "test@example.com", Role: "user"}
	token, exp, err := services.GenerateToken(user, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// expiry should be roughly now + configured hours
	if time.Until(exp) < time.Hour*70 || time.Until(exp) > time.Hour*74 {
		t.Errorf("Unexpected expiration duration: %v", time.Until(exp))
	}
	if token == "" {
		t.Error("Expected token string, got empty")
	}
}