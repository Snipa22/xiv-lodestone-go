package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

func RequiresAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the required X-Auth-Token header, make sure it's valid
		token := c.GetHeader("x-pool-auth")
		if len(token) == 0 || token != os.Getenv("API_AUTH_TOKEN") {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid auth token",
			})
			return
		}
		c.Next()
	}
}

func RequiresUserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the required X-Access-Token header, make sure it's valid
		token := c.GetHeader("x-access-token")
		if len(token) == 0 || token != os.Getenv("API_AUTH_TOKEN") {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid auth token",
			})
			return
		}
		c.Next()
	}
}
