package middleware

import "github.com/gin-gonic/gin"

func Permission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// mock
		c.Next()
	}
}
