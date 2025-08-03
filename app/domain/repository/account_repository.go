package repository

import (
	"context"

	"transaction_demo/app/domain/entity"
)

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

// AccountRepository represents the repository interface for the account entity
type AccountRepository interface {
	FindOne(ctx context.Context, id uint64) (*entity.Account, error)
	FindForUpdate(ctx context.Context, ids []uint64) ([]*entity.Account, error)
	Create(ctx context.Context, account *entity.Account) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
}
