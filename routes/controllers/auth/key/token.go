package key

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Username string    `json:"username"`
	Tpe      TokenType `json:"tpe"`
	jwt.RegisteredClaims
}

type TokenType string

const (
	AccessToken  TokenType = "Access Token"
	RefreshToken TokenType = "Refresh Token"
)

func NewNormalClaims(username string, id int64) CustomClaims {
	return CustomClaims{
		Username: username,
		Tpe:      AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatInt(id, 10),
			Issuer:    "auth-service",
		},
	}
}

func NewRefreshClaims(username string, id int64) CustomClaims {
	return CustomClaims{
		Username: username,
		Tpe:      RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatInt(id, 10),
			Issuer:    "auth-service",
		},
	}
}
