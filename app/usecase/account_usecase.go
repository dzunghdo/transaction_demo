package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/gin-gonic/gin"

	"transaction_demo/app/apperr"
	"transaction_demo/app/domain/entity"
	"transaction_demo/app/domain/repository"
	"transaction_demo/app/usecase/dto"
)

// AccountUC defines the interface for account-related business operations.
// Provides methods for account management and secure money transfers.
type AccountUC interface {
	// Create creates a new account with validation and duplicate checking.
	Create(ctx context.Context, account dto.AccountDTO) (dto.AccountDTO, error)

	// GetBalance retrieves the current balance of an account.
	GetBalance(ctx *gin.Context, id uint64) (dto.AccountDTO, error)

	// MakeTransaction performs atomic money transfer between accounts.
	MakeTransaction(c *gin.Context, req dto.TransactionDTO) error
}

type accountUsecase struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	txManager       trm.Manager
}

func NewAccountUsecase(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	txManager trm.Manager) AccountUC {
	return &accountUsecase{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		txManager:       txManager,
	}
}

// Create validates input, checks for duplicates, and creates a new account.
// Ensures no two accounts can have the same ID through database constraints.
func (uc accountUsecase) Create(ctx context.Context, account dto.AccountDTO) (dto.AccountDTO, error) {
	// Validate input data according to business rules
	err := account.Validate()
	if err != nil {
		fmt.Println("account validation failed", "error", err)
		return dto.AccountDTO{}, apperr.ErrInvalidInput.WithError(err)
	}

	// Check for existing account to provide clear error message
	existingAccount, err := uc.accountRepo.FindOne(ctx, account.AccountID)
	if err != nil {
		fmt.Println("failed to find account", "error", err)
		return dto.AccountDTO{}, apperr.ErrNotFound.WithMessage("account not found")
	}
	if existingAccount != nil {
		fmt.Println("account ID already exists", "account_id", account.AccountID)
		return dto.AccountDTO{}, apperr.ErrAlreadyExists.WithMessage("account ID already exists")
	}

	ent := entity.Account{
		ID:      account.AccountID,
		Balance: account.Balance,
	}
	createdAcc, err := uc.accountRepo.Create(ctx, &ent)
	if err != nil {
		fmt.Println("failed to create account", "error", err)
		return dto.AccountDTO{}, apperr.ErrInternalServer.WithError(err).WithMessage("failed to create account")
	}

	return dto.AccountDTO{
		AccountID: createdAcc.ID,
		Balance:   createdAcc.Balance,
	}, nil
}

// GetBalance returns account balance and details.
// Provides point-in-time snapshot without locking.
func (uc accountUsecase) GetBalance(ctx *gin.Context, id uint64) (dto.AccountDTO, error) {
	account, err := uc.accountRepo.FindOne(ctx, id)
	if err != nil {
		fmt.Println("failed to find account", "error", err)
		return dto.AccountDTO{}, apperr.ErrNotFound.WithMessage("account not found")
	}
	if account == nil {
		fmt.Println("account not found", "account_id", id)
		return dto.AccountDTO{}, apperr.ErrNotFound.WithMessage("account not found")
	}

	return dto.AccountDTO{
		AccountID: account.ID,
		Balance:   account.Balance,
	}, nil
}

// MakeTransaction performs atomic money transfer with deadlock prevention.
//
// Uses atomic multi-row locking strategy to handle high concurrency:
// - Locks both accounts simultaneously with single SELECT FOR UPDATE
// - Prevents deadlocks that occur with sequential account locking
// - Uses default READ COMMITTED isolation for optimal performance
// - Validates business rules within transaction boundary
// - Creates audit trail for all money movements
func (uc accountUsecase) MakeTransaction(ctx *gin.Context, req dto.TransactionDTO) error {
	// Validate transaction data
	err := req.Validate()
	if err != nil {
		fmt.Println("transaction validation failed", "error", err)
		return apperr.ErrInvalidInput.WithError(err).WithMessage(err.Error())
	}

	// Prevent self-transfers (business rule)
	if req.SourceAccountID == req.DestinationAccountID {
		fmt.Println("source and destination accounts have the same ID")
		return apperr.ErrInvalidInput.WithMessage("source and destination account IDs cannot be the same")
	}

	// Execute transaction with READ COMMITTED isolation
	// SERIALIZABLE is not needed since we explicitly lock required rows in a single operation
	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		// Lock both accounts atomically to prevent deadlocks
		sourceAcc, destAcc, err := uc.retrieveAccounts(ctx, req.SourceAccountID, req.DestinationAccountID)
		if err != nil {
			return err
		}

		// Validate business rules within transaction boundary
		if sourceAcc.Balance < req.Amount {
			fmt.Println("insufficient balance", "account_id", sourceAcc.ID, "balance", sourceAcc.Balance, "required", req.Amount)
			return apperr.ErrInvalidInput.WithMessage("insufficient balance")
		}

		// Execute the money transfer
		return uc.doTransaction(ctx, sourceAcc, destAcc, req.Amount)
	})

	if err != nil {
		fmt.Println("transaction failed", "error", err)
		return err
	}

	return nil
}

// retrieveAccounts locks both accounts atomically to prevent deadlocks.
//
// Deadlock Prevention Strategy:
// - Uses single SELECT FOR UPDATE with IN clause: WHERE id IN (x,y) FOR UPDATE
// - Avoids sequential locking which can cause circular wait conditions
// - Database locks both rows in consistent order regardless of parameter order
// - Returns error if either account doesn't exist
func (uc accountUsecase) retrieveAccounts(ctx context.Context, sourceAccID uint64, destAccID uint64,
) (*entity.Account, *entity.Account, error) {
	var (
		sourceAccount, destAccount *entity.Account
	)

	// Atomic locking prevents deadlocks that occur with sequential locking:
	// Instead of: LOCK(A) then LOCK(B) which can deadlock with LOCK(B) then LOCK(A)
	// We use: LOCK(A,B) atomically which eliminates circular wait conditions
	accounts, err := uc.accountRepo.FindForUpdate(ctx, []uint64{sourceAccID, destAccID})
	if err != nil {
		fmt.Println("failed to query accounts for update", "error", err)
		return nil, nil, apperr.ErrInternalServer.WithError(err).WithMessage("failed to find accounts for update")
	}

	// Ensure both accounts exist before proceeding
	if len(accounts) < 2 {
		fmt.Println("accounts not found for update")
		return nil, nil, apperr.ErrInternalServer.WithMessage("accounts not found for update")
	}

	// Map accounts by ID since database doesn't guarantee IN clause order
	for _, acc := range accounts {
		if acc.ID == sourceAccID {
			sourceAccount = acc
		} else if acc.ID == destAccID {
			destAccount = acc
		}
	}

	return sourceAccount, destAccount, nil
}

// doTransaction updates account balances and creates transaction log record.
// Operations performed atomically within the same database transaction:
// - Debits/credits accounts
// - Creates transaction record for audit trail
func (uc accountUsecase) doTransaction(
	ctx context.Context,
	sourceAccount *entity.Account,
	destinationAccount *entity.Account,
	amount float64,
) error {
	// Update account balances in memory
	sourceAccount.Balance -= amount
	destinationAccount.Balance += amount

	// Create transaction record for audit trail
	transaction := entity.Transaction{
		SourceAccountID:      sourceAccount.ID,
		DestinationAccountID: destinationAccount.ID,
		Amount:               amount,
		TransactionTime:      time.Now(),
	}

	// Save transaction record first for audit trail
	_, err := uc.transactionRepo.Create(ctx, &transaction)
	if err != nil {
		fmt.Println("transaction failed", "error", err)
		return apperr.ErrInternalServer.WithError(err).WithMessage("failed to create transaction")
	}

	// Persist account balance changes
	// Both updates occur within same DB transaction ensuring atomicity
	if err = uc.accountRepo.Update(ctx, sourceAccount); err != nil {
		fmt.Println("failed to update source account", "error", err)
		return apperr.ErrInternalServer.WithError(err).WithMessage("failed to update source account")
	}

	if err = uc.accountRepo.Update(ctx, destinationAccount); err != nil {
		fmt.Println("failed to update destination account", "error", err)
		return apperr.ErrInternalServer.WithError(err).WithMessage("failed to update destination account")
	}

	return nil
}
