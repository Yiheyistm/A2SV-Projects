package middleware

import (
	"net/http"
	"strings"
	"task_manager/config"
	"task_manager/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required on Authorization header"})
			c.Abort()
			return
		}
		jwtSecret := config.GetEnvString("JWT_SECRET", "my_secret_key")

		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}
		}
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
