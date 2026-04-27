package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/pkg/logger"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				// 获取 request_id
				var requestID any
				if rid, exists := c.Get(RequestIDKey); exists {
					requestID = rid
				}

				logger.Log.Error("Panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"request_id", requestID,
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"msg":        "internal server error",
					"request_id": requestID,
				})
			}
		}()

		c.Next()
	}
}
