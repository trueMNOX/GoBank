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
func (r *accountRepositoryImpl) GetByOwnerID(ctx context.Context, ownerID int64, limit, offset int) ([]*domain.Account, error) {
	var accountModel []models.AccountModel
	if err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Limit(limit).Offset(offset).Find(&accountModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	accounts := make([]*domain.Account, len(accountModel))
	for i, model := range accountModel {
		accounts[i] = &domain.Account{
			ID:        model.ID,
			OwnerID:   model.OwnerID,
			Balance:   model.Balance,
			Currency:  model.Currency,
			CreatedAt: model.CreatedAt,
		}
	}
	return accounts, nil
}
func (r *accountRepositoryImpl) UpdateBalance(ctx context.Context, id int64, amount int64) error {
	result := r.db.WithContext(ctx).Model(&models.AccountModel{}).Where("id = ?", id).Update("balance", gorm.Expr("balance + ?", amount))
	if result.Error != nil {
		return fmt.Errorf("failed to update balance: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}
func (r *accountRepositoryImpl) GetByIDForUpdate(ctx context.Context, id int64) (*domain.Account, error) {
	var accountModel models.AccountModel
	if err := r.db.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").First(&accountModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account for update: %w", err)
	}
	return &domain.Account{
		ID:        accountModel.ID,
		OwnerID:   accountModel.OwnerID,
		Balance:   accountModel.Balance,
		Currency:  accountModel.Currency,
		CreatedAt: accountModel.CreatedAt,
	}, nil
}
