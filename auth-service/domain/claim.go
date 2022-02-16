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

func (c AccessTokenClaims) IsUserRole() bool {
	return c.Role == "user"
}

func (c AccessTokenClaims) IsAdminRole() bool {
	return c.Role == "admin"
}

func (c AccessTokenClaims) IsValidCustomerID(customerID string) bool {
	return c.CustomerID == customerID
}

func (c AccessTokenClaims) IsValidAccountID(accountID string) bool {
	if accountID == "" {
		return true
	}

	for _, a := range c.Accounts {
		if a == accountID {
			return true
		}
	}
	return false
}

func (c AccessTokenClaims) IsRequestVerifiedWithTokenClaims(params map[string]string) bool {
	if c.CustomerID != params["customer_id"] {
		return false
	}
	if !c.IsValidAccountID(params["account_id"]) {
		return false
	}
	return true
}
