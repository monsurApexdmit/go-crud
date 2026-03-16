package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-crud/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token, claims, err := utils.ValidateJWT(parts[1])
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if utils.TokenBlacklist[parts[1]] {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token has been logged out",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("company_id", claims["company_id"])
		c.Set("email", claims["email"])
		c.Next()
	}
}

// GetCompanyID extracts the company_id from the Gin context
// Returns uint and bool (ok) for safe type assertion
func GetCompanyID(c *gin.Context) (uint, bool) {
	raw, exists := c.Get("company_id")
	if !exists {
		return 0, false
	}
	id, ok := raw.(float64)
	if !ok {
		return 0, false
	}
	return uint(id), true
}
