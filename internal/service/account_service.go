package service

import (
	"Gobank/internal/domain"
	"Gobank/internal/repository"
	"Gobank/internal/transport/http/dto"
	"context"
	"fmt"
)

type AccountService struct {
	accountRepo repository.AccountRepository
	userRepo    repository.UserRepository
}

func NewAccountService(accountRepo repository.AccountRepository, userRepo repository.UserRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo, userRepo: userRepo}
}
func (s *AccountService) CreateAccount(ctx context.Context, req *dto.CreateAccountRequest, ownerID int64) (*dto.AccountResponse, error) {
	_, err := s.userRepo.GetById(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	account := &domain.Account{
		OwnerID:  ownerID,
		Balance:  0,
		Currency: req.Currency,
	}
	createdAccount, err := s.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return &dto.AccountResponse{
		ID:        createdAccount.ID,
		OwnerID:   createdAccount.OwnerID,
		Balance:   createdAccount.Balance,
		Currency:  createdAccount.Currency,
		CreatedAt: createdAccount.CreatedAt,
	}, nil
}
func (s *AccountService) GetAccountByID(ctx context.Context, accountID int64, requesterID int64) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	if account.OwnerID != requesterID {
		return nil, fmt.Errorf("account does not belong to the requester")
	}
	return &dto.AccountResponse{
		ID:        account.ID,
		OwnerID:   account.OwnerID,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}, nil
}
func (s *AccountService) ListAccounts(ctx context.Context, ownerID int64, req *dto.ListAccountsRequest) ([]*dto.AccountResponse, error) {
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	account, err := s.accountRepo.GetByOwnerID(ctx, ownerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	response := make([]*dto.AccountResponse, len(account))
	for i, v := range account {
		response[i] = &dto.AccountResponse{
			ID:        v.ID,
			OwnerID:   v.OwnerID,
			Balance:   v.Balance,
			Currency:  v.Currency,
			CreatedAt: v.CreatedAt,
		}
	}
	return response, nil
}
