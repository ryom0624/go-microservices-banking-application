package domain

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	HmacSampleSecret     = "hmacSampleSecret"
	AccessTokenDuration  = time.Hour
	RefreshTokenDuration = 24 * 30 * time.Hour
)

type AccessTokenClaims struct {
	CustomerID string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	TokenType  string   `json:"token_type"`
	CustomerId string   `json:"cid"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"un"`
	Role       string   `json:"role"`
	jwt.StandardClaims
}

func (c AccessTokenClaims) RefreshTokenClaims() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType:  "refresh_token",
		CustomerId: c.CustomerID,
		Accounts:   c.Accounts,
		Username:   c.Username,
		Role:       c.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshTokenDuration).Unix(),
		},
	}
}