package repository

import (
	"Gobank/internal/domain"
	"context"

	"gorm.io/gorm"
)

type entryRepositoryImpl struct {
	db *gorm.DB
}

func NewEntryRepositoryImpl(db *gorm.DB) EntryRepository {
	return &entryRepositoryImpl{db: db}
}
func (r *entryRepositoryImpl) Create(ctx context.Context, entry *domain.Entry) (*domain.Entry, error) {
	entryModel := &domain.Entry{
		AccountID: entry.AccountID,
		Amount:    entry.Amount,
	}
	if err := r.db.WithContext(ctx).Create(entryModel).Error; err != nil {
		return nil, err
	}
	return entryModel, nil
}
func (r *entryRepositoryImpl) GetByAccountID(ctx context.Context, accountID int64, limit, offset int) ([]*domain.Entry, error) {
	var entries []domain.Entry
	if err := r.db.WithContext(ctx).Where("account_id", accountID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&entries).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	entriesDomain := make([]*domain.Entry, len(entries))
	for i, entry := range entries {
		entriesDomain[i] = &domain.Entry{
			ID:        entry.ID,
			AccountID: entry.AccountID,
			Amount:    entry.Amount,
			CreatedAt: entry.CreatedAt,
		}
	}
	return entriesDomain, nil
}
