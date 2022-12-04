package middleware

import (
	"github.com/gin-gonic/gin"
	"xiv-lodestone-go/support"
)

func SetupMilieu(milieu support.Milieu) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("MILIEU", milieu)
		c.Next()
	}
}
