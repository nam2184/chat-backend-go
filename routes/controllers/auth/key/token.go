package key

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	Username string
	ID       int64
	Type     TokenType
	Expiry   time.Duration
	Issuer   string
}

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

func NewClaims(cfg TokenConfig) CustomClaims {
	return CustomClaims{
		Username: cfg.Username,
		Tpe:      cfg.Type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatInt(cfg.ID, 10),
			Issuer:    cfg.Issuer,
		},
	}
}
