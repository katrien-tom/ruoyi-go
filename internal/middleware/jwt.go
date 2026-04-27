package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const defaultJWTSecret = "ruoyi-go-secret"

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

		claims, err := parseJWT(tokenString, jwtSecret())
		if err != nil {
			abortUnauthorized(c, err.Error())
			return
		}

		userID, ok := claimUserID(claims)
		if !ok {
			abortUnauthorized(c, "invalid user_id claim")
			return
		}

		c.Set("user_id", userID)
		c.Set("jwt_claims", claims)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"msg": msg,
	})
}

func jwtSecret() []byte {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		secret = defaultJWTSecret
	}

	return []byte(secret)
}

func parseJWT(token string, secret []byte) (map[string]any, error) {
	var claims map[string]any
	parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errJWT("unsupported token algorithm")
		}

		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errJWT("token expired")
		}
		return nil, errJWT("invalid token")
	}

	mapClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errJWT("invalid token")
	}

	claims = make(map[string]any, len(mapClaims))
	for key, value := range mapClaims {
		claims[key] = value
	}

	return claims, nil
}

func claimUserID(claims map[string]any) (int64, bool) {
	return claimInt64(claims["user_id"])
}

func claimInt64(v any) (int64, bool) {
	switch value := v.(type) {
	case float64:
		return int64(value), true
	case int64:
		return value, true
	case int:
		return int64(value), true
	case json.Number:
		n, err := value.Int64()
		return n, err == nil
	case string:
		var n json.Number = json.Number(value)
		i, err := n.Int64()
		return i, err == nil
	default:
		return 0, false
	}
}

type errJWT string

func (e errJWT) Error() string {
	return string(e)
}
