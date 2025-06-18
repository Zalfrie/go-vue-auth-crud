package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if user role is allowed
type contextKey string

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			return
		}
		userRole := role.(string)
		for _, r := range allowedRoles {
			if userRole == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	}
}