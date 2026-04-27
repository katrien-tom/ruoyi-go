package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/pkg/logger"
	"github.com/banyejiu/ruoyi-go/pkg/response"
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

				response.FailAbort(c, response.InternalError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
