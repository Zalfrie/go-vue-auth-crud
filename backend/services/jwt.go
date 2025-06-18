package services

import (
	"time"

	"go-vue-auth-crud/config"
	"go-vue-auth-crud/models"
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken creates a JWT token for given user
func GenerateToken(user models.User, rememberMe bool) (string, time.Time, error) {
	cfg := config.GetConfig()
	expiration := time.Hour * time.Duration(cfg.JWTExpirationHours)
	if rememberMe {
		expiration = expiration * 24 // extend for remember me (e.g. days)
	}
	expiresAt := time.Now().Add(expiration)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expiresAt.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	return signed, expiresAt, err
}

// GenerateResetToken creates a time-limited token for password reset
func GenerateResetToken(user models.User) (string, error) {
	cfg := config.GetConfig()
	expiration := time.Hour * 1 // 1 hour expiry
	expiresAt := time.Now().Add(expiration)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"type":    "reset",
		"exp":     expiresAt.Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString([]byte(cfg.JWTSecret))
}