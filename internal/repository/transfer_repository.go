package repository

import (
	"Gobank/internal/domain"
	"context"

	"gorm.io/gorm"
)

type TransferRepository interface {
	Create(ctx context.Context, transfer *domain.Transfer) (*domain.Transfer, error)
	GetByID(ctx context.Context, id int64) (*domain.Transfer, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Transfer, error)
	ExecuteTransfer(ctx context.Context, tx *gorm.DB, fromAccountID, toAccountID, amount int64, currency string) (*domain.Transfer, error)
}
