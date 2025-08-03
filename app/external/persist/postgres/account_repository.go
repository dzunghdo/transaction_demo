package postgres

import (
	"context"
	"errors"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"transaction_demo/app/domain/entity"
	"transaction_demo/app/domain/repository"
)

// accountRepository is the implementation of the AccountRepository interface
type accountRepository struct {
	db       *gorm.DB           // The database connection
	txGetter *trmgorm.CtxGetter // The transaction manager context getter
}

func NewAccountRepository(db *gorm.DB, txGetter *trmgorm.CtxGetter) repository.AccountRepository {
	return &accountRepository{db: db, txGetter: txGetter}
}

func (r accountRepository) FindOne(ctx context.Context, id uint64) (*entity.Account, error) {
	var ent entity.Account
	// get the transaction if exists, otherwise use the default database connection
	db := r.txGetter.DefaultTrOrDB(ctx, r.db).WithContext(ctx).Where("id = ?", id)

	err := db.First(&ent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &ent, err
}

func (r accountRepository) FindForUpdate(ctx context.Context, ids []uint64) ([]*entity.Account, error) {
	var ents []*entity.Account
	// get the transaction if exists, otherwise use the default database connection
	// Use SELECT FOR UPDATE to lock the rows for the duration of the transaction
	err := r.txGetter.DefaultTrOrDB(ctx, r.db).WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id IN ?", ids).
		Find(&ents).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return ents, err
}

func (r accountRepository) Create(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	// get the transaction if exists, otherwise use the default database connection
	db := r.txGetter.DefaultTrOrDB(ctx, r.db).WithContext(ctx)

	if err := db.Create(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (r accountRepository) Update(ctx context.Context, account *entity.Account) error {
	// get the transaction if exists, otherwise use the default database connection
	db := r.txGetter.DefaultTrOrDB(ctx, r.db).WithContext(ctx)
	return db.Save(account).Error
}
