package repository

import (
	"Gobank/internal/domain"
	"context"
)

type EntryRepository interface {
	Create(ctx context.Context, entry *domain.Entry) (*domain.Entry, error)
	GetByAccountID(ctx context.Context, accountID int64, limit, offset int) ([]*domain.Entry, error)
}
