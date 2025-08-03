package repository

import (
	"context"

	"transaction_demo/app/domain/entity"
)

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

// TransactionRepository represents the repository interface for the transaction entity
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
}
