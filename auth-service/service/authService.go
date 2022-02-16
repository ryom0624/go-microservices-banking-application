package service

import (
	"auth/domain"
	"auth/dto"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"local.packages/lib/errs"
	"local.packages/lib/logger"
)

type AuthService interface {
	Login(request dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (d DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	jwtToken, err := jwtTokenFromString(urlParams["token"])
	if err != nil {
		return errs.NewAuthorizationError(err.Error())
	}
	if !jwtToken.Valid {
		return errs.NewAuthorizationError("Invalid token")
	}
	claims := jwtToken.Claims.(*domain.AccessTokenClaims)
	if claims.IsUserRole() {
		if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
			return errs.NewAuthorizationError("request not verified with the token claims")
		}
	}
	if !d.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"]) {
		return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))
	}

	return nil
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HmacSampleSecret), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func (d DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	login, appError := d.repo.FindBy(req.Username, req.Password)
	if appError != nil {
		return nil, appError
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	accessToken, appError := authToken.NewAccessToken()
	if appError != nil {
		return nil, appError
	}
	refreshToken, appError := d.repo.GenerateAndSaveRefreshTokenToStore(authToken)
	if appError != nil {
		return nil, appError
	}
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewAuthService(repo domain.AuthRepository, permissions domain.RolePermissions) AuthService {
	return DefaultAuthService{repo: repo, rolePermissions: permissions}
}
