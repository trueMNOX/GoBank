package repository

import (
	"Gobank/internal/domain"
	"Gobank/internal/repository/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type accountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepositoryImpl{db: db}
}

func (r *accountRepositoryImpl) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	accountModel := &models.AccountModel{
		OwnerID:  account.OwnerID,
		Balance:  account.Balance,
		Currency: account.Currency,
	}
	if err := r.db.WithContext(ctx).Create(accountModel).Error; err != nil {
		return nil, err
	}
	return &domain.Account{
		ID:       accountModel.ID,
		OwnerID:  accountModel.OwnerID,
		Balance:  accountModel.Balance,
		Currency: accountModel.Currency,
	}, nil
}
func (r *accountRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	var accountModel models.AccountModel
	if err := r.db.WithContext(ctx).First(&accountModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &domain.Account{
		ID:        accountModel.ID,
		OwnerID:   accountModel.OwnerID,
		Balance:   accountModel.Balance,
		Currency:  accountModel.Currency,
		CreatedAt: accountModel.CreatedAt,
	}, nil
}
