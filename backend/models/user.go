package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
)

// User represents a system user
// @swagger:model
// User
type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`  
	Email     string    `gorm:"size:100;unique_index;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`       
	Role      string    `gorm:"size:20;not null;default:'user'" json:"role"` // 'user' or 'admin'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeSave GORM hook to hash password
func (u *User) BeforeSave(scope *gorm.Scope) error {
	if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost); err == nil {
		scope.SetColumn("Password", string(pw))
	} else {
		return err
	}
	return nil
}

// CheckPassword compares plain password with hashed
func (u *User) CheckPassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
}