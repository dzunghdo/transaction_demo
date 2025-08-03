package postgres

import (
	"context"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"

	"transaction_demo/app/domain/entity"
	"transaction_demo/app/domain/repository"
)

// transactionRepository is the implementation of the TransactionRepository interface
type transactionRepository struct {
	db       *gorm.DB           // The database connection
	txGetter *trmgorm.CtxGetter // The transaction manager context getter
}

func NewTransactionRepository(db *gorm.DB, txGetter *trmgorm.CtxGetter) repository.TransactionRepository {
	return &transactionRepository{
		db:       db,
		txGetter: txGetter,
	}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	// get the transaction if exists, otherwise use the default database connection
	db := r.txGetter.DefaultTrOrDB(ctx, r.db).WithContext(ctx)
	if err := db.Create(&transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}
