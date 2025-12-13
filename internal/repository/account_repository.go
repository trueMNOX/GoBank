package repository

import (
	"Gobank/internal/domain"
	"context"
)

type AccountRepository interface {
    Create(ctx context.Context, account *domain.Account) (*domain.Account, error)
    GetByID(ctx context.Context, id int64) (*domain.Account, error)
    GetByOwnerID(ctx context.Context, ownerID int64, limit, offset int) ([]*domain.Account, error)
    UpdateBalance(ctx context.Context, id int64, amount int64) error
    GetByIDForUpdate(ctx context.Context, id int64) (*domain.Account, error)
}
