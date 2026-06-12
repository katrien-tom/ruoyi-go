package middleware

import (
	"errors"
	"strings"

	"github.com/banyejiu/ruoyi-go/internal/security"
	"github.com/banyejiu/ruoyi-go/pkg/jwtutil"
	"github.com/banyejiu/ruoyi-go/pkg/response"
	"github.com/gin-gonic/gin"
)

func JWTAuth(sessionService *security.SessionService) gin.HandlerFunc {
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

		if claims.UserID == 0 || strings.TrimSpace(claims.Token) == "" {
			abortUnauthorized(c, "invalid login token")
			return
		}

		loginUser, _, err := sessionService.VerifyAndRefresh(c.Request.Context(), claims.Token)
		if err != nil {
			switch {
			case errors.Is(err, security.ErrTokenNotFound), errors.Is(err, security.ErrTokenInvalid):
				abortUnauthorized(c, "login state expired")
			default:
				response.FailAbort(c, response.InternalError, "failed to verify token")
			}
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("access_token", tokenString)
		c.Set("login_token", claims.Token)
		c.Set("login_user", loginUser)
		c.Set("jwt_claims", claims)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, msg string) {
	response.FailAbort(c, response.Unauthorized, msg)
}
