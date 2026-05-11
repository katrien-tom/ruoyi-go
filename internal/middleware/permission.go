package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/security"
	"github.com/banyejiu/ruoyi-go/pkg/response"
)

func Permission(code string) gin.HandlerFunc {
	code = strings.TrimSpace(code)

	return func(c *gin.Context) {
		if code == "" {
			c.Next()
			return
		}

		loginUserValue, ok := c.Get("login_user")
		if !ok {
			response.FailAbort(c, response.Unauthorized, "missing login user")
			return
		}

		loginUser, ok := loginUserValue.(*security.LoginUser)
		if !ok || loginUser == nil {
			response.FailAbort(c, response.Unauthorized, "invalid login user")
			return
		}

		for _, permission := range loginUser.Permissions {
			if permission == "*:*:*" || permission == code {
				c.Next()
				return
			}
		}

		response.FailAbort(c, response.PermissionDenied, "permission denied")
	}
}
