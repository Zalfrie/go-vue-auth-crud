package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"go-vue-auth-crud/models"
)

// ListUsers returns paginated user list (for datatable)
func ListUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User
	// simple pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	db.Offset(offset).Limit(limit).Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUser returns a single user by ID
func GetUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser modifies user data
func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	var input struct { Name, Role string }
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Model(&user).Updates(models.User{Name: input.Name, Role: input.Role})
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser removes a user (admin only)
func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User dihapus"})
}

// ExportUsersExcel generates an Excel file of all users
func ExportUsersExcel(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User
	db.Find(&users)

	f := excelize.NewFile()
	sheet := f.NewSheet("Users")
	// header
	f.SetCellValue("Users", "A1", "ID")
	f.SetCellValue("Users", "B1", "Name")
	f.SetCellValue("Users", "C1", "Email")
	f.SetCellValue("Users", "D1", "Role")
	// data
	for i, u := range users {
		row := i + 2
		f.SetCellValue("Users", "A"+strconv.Itoa(row), u.ID)
		f.SetCellValue("Users", "B"+strconv.Itoa(row), u.Name)
		f.SetCellValue("Users", "C"+strconv.Itoa(row), u.Email)
		f.SetCellValue("Users", "D"+strconv.Itoa(row), u.Role)
	}
	f.SetActiveSheet(sheet)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=users.xlsx")
	c.Header("File-Name", "users.xlsx")
	_ = f.Write(c.Writer)
}