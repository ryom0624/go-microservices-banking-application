package domain

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type Login struct {
	Username   string         `db:"username"`
	CustomerID sql.NullString `db:"customer_id"`
	Accounts   sql.NullString `db:"account_numbers"`
	Role       string         `db:"role"`
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	if l.Accounts.Valid && l.CustomerID.Valid {
		return l.claimForUser()
	} else {
		return l.claimForAdmin()
	}
}

func (l Login) claimForUser() AccessTokenClaims {
	accounts := strings.Split(l.Accounts.String, ",")
	return AccessTokenClaims{
		CustomerID: l.CustomerID.String,
		Accounts:   accounts,
		Username:   l.Username,
		Role:       l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenDuration).Unix(),
		},
	}
}

func (l Login) claimForAdmin() AccessTokenClaims {
	return AccessTokenClaims{
		Username: l.Username,
		Role:     l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenDuration).Unix(),
		},
	}
}
