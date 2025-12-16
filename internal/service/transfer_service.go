package service

import (
	"Gobank/internal/cache"
	"Gobank/internal/queue/rabbitmq"
	"Gobank/internal/repository"
	"Gobank/internal/transport/http/dto"
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type TransferService struct {
	transferRepo repository.TransferRepository
	db           *gorm.DB
	accountRepo  repository.AccountRepository
	cache        cache.Cache
	producer     rabbitmq.Producer
}

func NewTransferService(transferRepo repository.TransferRepository, db *gorm.DB, accountRepo repository.AccountRepository, cache cache.Cache, producer rabbitmq.Producer) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		db:           db,
		accountRepo:  accountRepo,
		cache:        cache,
		producer:     producer,
	}
}
func (s *TransferService) CreateTransfer(ctx context.Context, req *dto.TransferRequest, requesterID int64, idempotencyKey string) (*dto.TransferResponse, error) {
	if idempotencyKey != "" {
		val, _ := s.cache.Get(ctx, "idemp:"+idempotencyKey)
		if val != "" {
			return nil, fmt.Errorf("transfer already processed with this key")
		}
	}
	if req.FromAccountID == req.ToAccountID {
		return nil, fmt.Errorf("cannot transfer to the same account")
	}
	fromAccount, err := s.accountRepo.GetByID(ctx, req.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get from account: %w", err)
	}
	if fromAccount.OwnerID != requesterID {
		return nil, fmt.Errorf("account doesn't belong to the authenticated user")
	}
	_, err = s.accountRepo.GetByID(ctx, req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get to account: %w", err)
	}
	var createdTransfer *dto.TransferResponse

	err = s.db.Transaction(func(tx *gorm.DB) error {
		transfer, terr := s.transferRepo.ExecuteTransfer(ctx, tx, req.FromAccountID, req.ToAccountID, req.Amount, req.Currency)
		if terr != nil {
			return terr
		}
		createdTransfer = &dto.TransferResponse{
			ID:            transfer.ID,
			FromAccountID: transfer.FromAccountID,
			ToAccountID:   transfer.ToAccountID,
			Amount:        transfer.Amount,
			Currency:      transfer.Currency,
			CreatedAt:     transfer.CreatedAt,
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transfer failed: %w", err)
	}
	if idempotencyKey != "" {
		s.cache.SetIdempotencyKey(ctx, idempotencyKey, "processed", 24*time.Hour)
	}

	go func() {
		err := s.producer.PublishTransferEvent(context.Background(), rabbitmq.TransferEvent{
			TransferID:    createdTransfer.ID,
			FromAccountID: createdTransfer.FromAccountID,
			ToAccountID:   createdTransfer.ToAccountID,
			Amount:        createdTransfer.Amount,
			Currency:      createdTransfer.Currency,
		})
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
		}
	}()
	return createdTransfer, nil
}
