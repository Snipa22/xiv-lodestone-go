package support

import "github.com/gin-gonic/gin"

func Err500(c *gin.Context) {
	c.AbortWithStatusJSON(500, gin.H{
		"error": "Internal Server Error",
	})
}

func Err403(c *gin.Context) {
	c.AbortWithStatusJSON(403, gin.H{
		"error": "Access Forbidden",
	})
}

func Err400(c *gin.Context, errorStr string) {
	c.AbortWithStatusJSON(403, gin.H{
		"error": errorStr,
	})
}

func JsonSuccess(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
	})
}

func NoResult200(c *gin.Context) {
	c.Status(200)
}
