package routes

import (
	"go-vue-auth-crud/config"
	"go-vue-auth-crud/controllers"
	"go-vue-auth-crud/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DBMiddleware injects *gorm.DB into context
func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

// RegisterRoutes sets up API routes and middleware
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// global middleware
	r.Use(DBMiddleware(db))

	// public routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/forgot-password", controllers.ForgotPassword)

	// protected routes
	auth := r.Group("/api")
	auth.Use(middlewares.JWTAuthMiddleware(cfg))
	{
		auth.GET("/users", controllers.ListUsers)
		auth.GET("/users/export", controllers.ExportUsersExcel)
		auth.GET("/users/:id", controllers.GetUser)
		auth.PUT("/users/:id", controllers.UpdateUser)
		// delete only admin
		auth.DELETE("/users/:id", middlewares.RoleMiddleware("admin"), controllers.DeleteUser)
	}
}