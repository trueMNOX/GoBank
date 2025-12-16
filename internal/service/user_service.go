package service

import (
	"Gobank/internal/repository"
	"Gobank/internal/transport/http/dto"
	"context"
	"fmt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
