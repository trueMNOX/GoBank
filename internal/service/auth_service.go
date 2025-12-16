package service

import (
	"Gobank/internal/domain"
	"Gobank/internal/repository"
	"Gobank/internal/token"
	"Gobank/internal/transport/http/dto"
	"Gobank/pkg/config"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo             repository.UserRepository
	tokenMaker           token.TokenMaker
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) (*AuthService, error) {
	tokenMaker, err := token.NewJwtMaker(cfg.JWTSecret)
	if err != nil {
		return nil, err
	}
	return &AuthService{
		userRepo:             userRepo,
		tokenMaker:           tokenMaker,
		accessTokenDuration:  time.Duration(cfg.AccessTokenDuration),
		refreshTokenDuration: time.Duration(cfg.RefreshTokenDuration),
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	existingUser, _, err := s.userRepo.GetByUserName(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}
	existingEmail, _, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return nil, fmt.Errorf("email already exists")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
	}
	createdUser, err := s.userRepo.Create(ctx, user, string(hashPassword))
	if err != nil {
		return nil, err
	}
	accessToken, err := s.tokenMaker.CreateToken(createdUser.Username, createdUser.ID, s.accessTokenDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokenMaker.CreateToken(createdUser.Username, createdUser.ID, s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       createdUser.ID,
			Username: createdUser.Username,
			FullName: createdUser.FullName,
			Email:    createdUser.Email,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, hashPassword, err := s.userRepo.GetByUserName(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}
	accessToken, err := s.tokenMaker.CreateToken(user.Username, user.ID, s.accessTokenDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokenMaker.CreateToken(user.Username, user.ID, s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
	}, nil
}
func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (string, error) {
	payload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return "", err
	}
	accessToken, err := s.tokenMaker.CreateToken(payload.Username, payload.UserID, s.accessTokenDuration)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
