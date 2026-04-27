package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/pkg/logger"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// 先执行后续逻辑
		c.Next()

		latency := time.Since(start)

		// 基础字段
		fields := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
		}

		// request_id
		if rid, exists := c.Get(RequestIDKey); exists {
			fields = append(fields, "request_id", rid)
		}

		// 错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, "errors", c.Errors.String())
		}

		// 根据状态码分级
		switch {
		case c.Writer.Status() >= 500:
			logger.Log.Error("HTTP Request", fields...)
		case c.Writer.Status() >= 400:
			logger.Log.Warn("HTTP Request", fields...)
		default:
			logger.Log.Info("HTTP Request", fields...)
		}
	}
}
