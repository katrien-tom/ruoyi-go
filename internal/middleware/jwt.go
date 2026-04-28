package middleware

import (
	"strings"

	"github.com/banyejiu/ruoyi-go/pkg/jwtutil"
	"github.com/banyejiu/ruoyi-go/pkg/response"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortUnauthorized(c, "missing authorization header")
			return
		}

		tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || tokenString == "" {
			abortUnauthorized(c, "invalid authorization header")
			return
		}

		claims, err := jwtutil.Parse(tokenString)
		if err != nil {
			abortUnauthorized(c, err.Error())
			return
		}

		if claims.UserID == 0 {
			abortUnauthorized(c, "invalid user_id claim")
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("jwt_claims", claims)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, msg string) {
	response.FailAbort(c, response.Unauthorized, msg)
}
