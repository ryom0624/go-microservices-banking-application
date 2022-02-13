package service

import (
	"auth/domain"
	"auth/dto"
	"local.packages/lib/errs"
)

type AuthService interface {
	Login(request dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo domain.AuthRepository
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

func NewAuthService(repo domain.AuthRepository) AuthService {
	return DefaultAuthService{repo: repo}
}
