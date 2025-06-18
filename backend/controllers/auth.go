package controllers

import (
	"net/http"
	"time"

	"go-vue-auth-crud/config"
	"go-vue-auth-crud/models"
	"go-vue-auth-crud/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RegisterInput struct
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput struct
type LoginInput struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

// ForgotInput struct
type ForgotInput struct {
	Email string `json:"email" binding:"required,email"`
}

// Register new user
// @Summary Register a new user
// @Description Create a new user account and send notification emails
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Registration data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	// create user
	user := models.User{Name: input.Name, Email: input.Email, Password: input.Password}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// send email notifications
	cfg := config.GetConfig()
	emailSvc := services.NewEmailService(cfg)
	emailSvc.SendRegistration(user.Email, user.Name)
	emailSvc.SendRegistration(cfg.AdminEmail, user.Name+" (admin)")

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil"})
}

// Login existing user
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// generate JWT
	token, expiresAt, err := services.GenerateToken(user, input.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"expires_at": expiresAt,
		"user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
	})
}

// ForgotPassword send reset link
func ForgotPassword(c *gin.Context) {
	var input ForgotInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Jika email terdaftar, link reset akan dikirim."})
		return
	}

	// create reset token
	token, err := services.GenerateResetToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token reset"})
		return
	}

	cfg := config.GetConfig()
	emailSvc := services.NewEmailService(cfg)
	emailSvc.SendPasswordReset(user.Email, token)

	c.JSON(http.StatusOK, gin.H{"message": "Link reset password telah dikirim."})
}