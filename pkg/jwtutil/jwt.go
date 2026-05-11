package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const defaultJWTSecret = "ruoyi-go-secret"
const defaultJWTTTL = 24 * time.Hour

var (
	ErrUnsupportedTokenAlgorithm = errors.New("unsupported token algorithm")
	ErrInvalidToken              = errors.New("invalid token")
	ErrTokenExpired              = errors.New("token expired")
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
	jwt.RegisteredClaims
}

func Sign(userID int64, userName, loginToken string, now time.Time) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: userName,
		Token:    loginToken,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(defaultJWTTTL)),
		},
	}

	signedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return signedToken.SignedString(secret())
}

func Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrUnsupportedTokenAlgorithm
		}

		return secret(), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, ErrUnsupportedTokenAlgorithm):
			return nil, ErrUnsupportedTokenAlgorithm
		default:
			return nil, ErrInvalidToken
		}
	}
	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func ExpiresAt(now time.Time) time.Time {
	return now.Add(defaultJWTTTL)
}

func secret() []byte {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		secret = defaultJWTSecret
	}

	return []byte(secret)
}
