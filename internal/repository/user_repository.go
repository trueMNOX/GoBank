package repository

import (
	"Gobank/internal/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User, hashedPassword string) (*domain.User, error)
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetByUserName(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
