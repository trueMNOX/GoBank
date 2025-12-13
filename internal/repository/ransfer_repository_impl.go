package repository

import (
	"Gobank/internal/domain"
	"Gobank/internal/repository/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type transferRepositoryImpl struct {
	db          *gorm.DB
	entryRepo   EntryRepository
	accountRepo AccountRepository
}

func NewTransferRepository(db *gorm.DB, entryRepo EntryRepository, accountRepo AccountRepository) TransferRepository {
	return &transferRepositoryImpl{db: db, entryRepo: entryRepo, accountRepo: accountRepo}
}

func (r *transferRepositoryImpl) Create(ctx context.Context, transfer *domain.Transfer) (*domain.Transfer, error) {
	transferModel := &models.TransferModel{
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Amount:        transfer.Amount,
		Currency:      transfer.Currency,
	}
	if err := r.db.WithContext(ctx).Create(transferModel).Error; err != nil {
		return nil, err
	}
	return &domain.Transfer{
		ID:            transferModel.ID,
		FromAccountID: transferModel.FromAccountID,
		ToAccountID:   transferModel.ToAccountID,
		Amount:        transferModel.Amount,
		Currency:      transferModel.Currency,
		CreatedAt:     transfer.CreatedAt,
	}, nil
}
func (r *transferRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Transfer, error) {
	var transferModel models.TransferModel
	if err := r.db.WithContext(ctx).First(&transferModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("transfer not found")
		}
		return nil, fmt.Errorf("failed to get transfer: %w", err)
	}
	return &domain.Transfer{
		ID:            transferModel.ID,
		FromAccountID: transferModel.FromAccountID,
		ToAccountID:   transferModel.ToAccountID,
		Amount:        transferModel.Amount,
		Currency:      transferModel.Currency,
		CreatedAt:     transferModel.CreatedAt,
	}, nil
}
func (r *transferRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*domain.Transfer, error) {
	var transferModels []models.TransferModel
	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&transferModels).Error; err != nil {
		return nil, fmt.Errorf("failed to list transfers: %w", err)
	}
	transfers := make([]*domain.Transfer, len(transferModels))
	for i, model := range transferModels {
		transfers[i] = &domain.Transfer{
			ID:            model.ID,
			FromAccountID: model.FromAccountID,
			ToAccountID:   model.ToAccountID,
			Amount:        model.Amount,
			Currency:      model.Currency,
			CreatedAt:     model.CreatedAt,
		}
	}
	return transfers, nil
}
func (r *transferRepositoryImpl) getAccountForUpdateInTx(ctx context.Context, tx *gorm.DB, id int64) (*domain.Account, error) {
	var accountModel models.AccountModel
	if err := tx.WithContext(ctx).Clauses(gorm.Expr("FOR UPDATE")).First(&accountModel, id).Error; err != nil {
		return nil, err
	}
	return &domain.Account{
		ID:        accountModel.ID,
		OwnerID:   accountModel.OwnerID,
		Balance:   accountModel.Balance,
		Currency:  accountModel.Currency,
		CreatedAt: accountModel.CreatedAt,
	}, nil
}
func (r *transferRepositoryImpl) UpdateBalanceInTx(ctx context.Context, tx *gorm.DB, id int64, amount int64) error {
	result := tx.WithContext(ctx).Model(&models.AccountModel{}).Where("id = ?", id).Update("balance", gorm.Expr("balance + ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("account not found or no rows affected")
	}
	return nil
}
func (r *transferRepositoryImpl) ExecuteTransfer(ctx context.Context, tx *gorm.DB, fromAccountID, toAccountID, amount int64, curreny string) (*domain.Transfer, error) {
	fromAccount, err := r.getAccountForUpdateInTx(ctx, tx, fromAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get from account: %w", err)
	}
	toAccount, err := r.getAccountForUpdateInTx(ctx, tx, toAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get to account: %w", err)
	}
	if fromAccount.Currency != curreny || toAccount.Currency != curreny {
		return nil, fmt.Errorf("currency mismatch: from account currency is %s, to account currency is %s", fromAccount.Currency, toAccount.Currency)
	}
	if fromAccount.Balance < amount {
		return nil, fmt.Errorf("insufficient balance in from account")
	}
	if err := r.UpdateBalanceInTx(ctx, tx, fromAccountID, -amount); err != nil {
		return nil, fmt.Errorf("failed to update from account balance: %w", err)
	}
	if err := r.UpdateBalanceInTx(ctx, tx, toAccountID, amount); err != nil {
		return nil, fmt.Errorf("failed to update to account balance: %w", err)
	}
	transferModel := &models.TransferModel{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		Currency:      curreny,
	}
	if err := tx.WithContext(ctx).Create(transferModel).Error; err != nil {
		return nil, fmt.Errorf("failed to create transfer: %w", err)
	}
	fromEntry := &models.EntryModel{
		AccountID: fromAccountID,
		Amount:    -amount,
	}
	if err := tx.WithContext(ctx).Create(fromEntry).Error; err != nil {
		return nil, fmt.Errorf("failed to create from entry: %w", err)
	}
	toEntry := &models.EntryModel{
		AccountID: toAccountID,
		Amount:    amount,
	}
	if err := tx.WithContext(ctx).Create(toEntry).Error; err != nil {
		return nil, fmt.Errorf("failed to create to entry: %w", err)
	}
	return &domain.Transfer{
		ID:            transferModel.ID,
		FromAccountID: transferModel.FromAccountID,
		ToAccountID:   transferModel.ToAccountID,
		Amount:        transferModel.Amount,
		Currency:      transferModel.Currency,
		CreatedAt:     transferModel.CreatedAt,
	}, nil
}
