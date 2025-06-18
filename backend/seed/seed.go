package seed

import (
	"go-vue-auth-crud/models"
	"github.com/jinzhu/gorm"
)

// Run seeds initial data
func Run(db *gorm.DB) error {
	// create admin user if not exists
	admin := models.User{ Name: "Administrator", Email: "admin@example.com", Password: "Admin123!", Role: "admin" }
	if err := db.Where("email = ?", admin.Email).FirstOrCreate(&admin).Error; err != nil {
		return err
	}
	// create sample user
	user := models.User{ Name: "Sample User", Email: "user@example.com", Password: "User123!", Role: "user" }
	if err := db.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
		return err
	}
	return nil
}