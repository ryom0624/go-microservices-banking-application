package domain

import (
	"github.com/dgrijalva/jwt-go"
	"local.packages/lib/errs"
	"local.packages/lib/logger"
)

type AuthToken struct {
	token *jwt.Token
}

func (t AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := t.token.SignedString([]byte(HmacSampleSecret))
	if err != nil {
		logger.Error("Failed while signing access token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate access token")
	}
	return signedString, nil
}

func (t AuthToken) newRefreshToken() (string, *errs.AppError) {
	claims := t.token.Claims.(AccessTokenClaims)
	refreshClaims := claims.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedString, err := token.SignedString([]byte(HmacSampleSecret))
	if err != nil {
		logger.Error("Failed while signing refresh token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate refresh token")
	}
	return signedString, nil
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}
