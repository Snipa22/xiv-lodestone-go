package middleware

import (
	"github.com/Snipa22/xiv-lodestone-go/support"
	"github.com/gin-gonic/gin"
)

func SetupMilieu(milieu support.Milieu) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("MILIEU", milieu)
		c.Next()
	}
}
