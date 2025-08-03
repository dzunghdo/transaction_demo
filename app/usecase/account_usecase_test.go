package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"transaction_demo/app/domain/entity"
	mock2 "transaction_demo/cmd/shared/db/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	"transaction_demo/app/domain/repository/mock"
	"transaction_demo/app/usecase/dto"
)

type fields struct {
	accountRepo     *mock.MockAccountRepository
	transactionRepo *mock.MockTransactionRepository
	txManager       *mock2.MockTxManager
}

func Test_accountUsecase_Create(t *testing.T) {
	type args struct {
		ctx     context.Context
		account dto.AccountDTO
	}
	tests := []struct {
		name    string
		args    args
		setup   func(fields fields)
		want    dto.AccountDTO
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				account: dto.AccountDTO{AccountID: 111, Balance: 1000},
			},
			setup: func(fields fields) {
				fields.accountRepo.EXPECT().FindOne(gomock.Any(), uint64(111)).Return(nil, nil)
				fields.accountRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.Account{
					ID:      111,
					Balance: 1000,
				}, nil)
			},
			want:    dto.AccountDTO{AccountID: 111, Balance: 1000},
			wantErr: false,
		},
		{
			name: "validation_error_invalid_account_id",
			args: args{
				ctx:     context.Background(),
				account: dto.AccountDTO{AccountID: 0, Balance: 1000}, // Invalid AccountID
			},
			want:    dto.AccountDTO{},
			wantErr: true,
		},
		{
			name: "validation_error_invalid_balance",
			args: args{
				ctx:     context.Background(),
				account: dto.AccountDTO{AccountID: 111, Balance: -100}, // Invalid Balance
			},
			want:    dto.AccountDTO{},
			wantErr: true,
		},
		{
			name: "find_one_error",
			args: args{
				ctx:     context.Background(),
				account: dto.AccountDTO{AccountID: 111, Balance: 1000},
			},
			setup: func(fields fields) {
				fields.accountRepo.EXPECT().FindOne(gomock.Any(), uint64(111)).
					Return(nil, errors.New("database error"))
			},
			want:    dto.AccountDTO{},
			wantErr: true,
		},
		{
			name: "account_already_exists",
			args: args{
				ctx:     context.Background(),
				account: dto.AccountDTO{AccountID: 111, Balance: 1000},
			},
			setup: func(fields fields) {
				existingAccount := &entity.Account{
					ID:      111,
					Balance: 500,
				}
				fields.accountRepo.EXPECT().FindOne(gomock.Any(), uint64(111)).Return(existingAccount, nil)
			},
			want:    dto.AccountDTO{},
			wantErr: true,
		},
		{
			name: "create_error",
			args: args{
				ctx:     context.Background(),
				account: dto.AccountDTO{AccountID: 111, Balance: 1000},
			},
			setup: func(fields fields) {
				fields.accountRepo.EXPECT().FindOne(gomock.Any(), uint64(111)).Return(nil, nil)
				fields.accountRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database constraint violation"))
			},
			want:    dto.AccountDTO{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAccountRepo := mock.NewMockAccountRepository(ctrl)
			mockTransactionRepo := mock.NewMockTransactionRepository(ctrl)

			uc := accountUsecase{
				accountRepo:     mockAccountRepo,
				transactionRepo: mockTransactionRepo,
				txManager:       &mock2.MockTxManager{},
			}

			testFields := fields{
				accountRepo:     mockAccountRepo,
				transactionRepo: mockTransactionRepo,
				txManager:       &mock2.MockTxManager{},
			}

			if tt.setup != nil {
				tt.setup(testFields)
			}

			got, err := uc.Create(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountUsecase_GetBalance(t *testing.T) {
	type args struct {
		ctx *gin.Context
		id  uint64
	}
	tests := []struct {
		name    string
		args    args
		setup   func(fields fields)
		want    dto.AccountDTO
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: &gin.Context{},
				id:  111,
			},
			setup: func(fields fields) {
				fields.accountRepo.EXPECT().FindOne(gomock.Any(), uint64(111)).Return(&entity.Account{
					ID:      111,
					Balance: 1500.75,
				}, nil)
			},
			want: dto.AccountDTO{
				AccountID: 111,
				Balance:   1500.75,
			},
			wantErr: false,
		},
		{
			name: "find_one_error",
			args: args{
				ctx: &gin.Context{},
				id:  111,
			},
			setup: func(fields fields) {
				fields.accountRepo.EXPECT().FindOne(gomock.Any(), uint64(111)).Return(nil, errors.New("database connection failed"))
			},
			want:    dto.AccountDTO{},
			wantErr: true,
		},
		{
			name: "account_not_found",
			args: args{
				ctx: &gin.Context{},
				id:  999,
			},
			setup: func(fields fields) {
				fields.accountRepo.EXPECT().FindOne(gomock.Any(), uint64(999)).Return(nil, nil)
			},
			want:    dto.AccountDTO{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAccountRepo := mock.NewMockAccountRepository(ctrl)
			mockTransactionRepo := mock.NewMockTransactionRepository(ctrl)

			uc := accountUsecase{
				accountRepo:     mockAccountRepo,
				transactionRepo: mockTransactionRepo,
				txManager:       mock2.NewMockTxManager(),
			}

			testFields := fields{
				accountRepo:     mockAccountRepo,
				transactionRepo: mockTransactionRepo,
				txManager:       mock2.NewMockTxManager(),
			}

			if tt.setup != nil {
				tt.setup(testFields)
			}

			got, err := uc.GetBalance(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountUsecase_MakeTransaction(t *testing.T) {
	type args struct {
		ctx *gin.Context
		req dto.TransactionDTO
	}
	tests := []struct {
		name    string
		args    args
		setup   func(fields fields)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               100.50,
				},
			},
			setup: func(fields fields) {
				// Mock FindForUpdate to return both accounts
				accounts := []*entity.Account{
					{ID: 111, Balance: 1000.00},
					{ID: 222, Balance: 500.00},
				}

				fields.txManager.ShouldFail = false
				fields.accountRepo.EXPECT().FindForUpdate(gomock.Any(), []uint64{111, 222}).Return(accounts, nil)
				fields.transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.Transaction{}, nil)
				fields.accountRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			wantErr: false,
		},
		{
			name: "validation_error_invalid_source_account_id",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      0, // Invalid
					DestinationAccountID: 222,
					Amount:               100.50,
				},
			},
			wantErr: true,
		},
		{
			name: "validation_error_invalid_destination_account_id",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 0, // Invalid
					Amount:               100.50,
				},
			},
			wantErr: true,
		},
		{
			name: "validation_error_invalid_amount",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               -100.50, // Invalid
				},
			},
			wantErr: true,
		},
		{
			name: "self_transfer_error",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 111, // Same as source
					Amount:               100.50,
				},
			},
			wantErr: true,
		},
		{
			name: "tx_manager_error",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               100.50,
				},
			},
			setup: func(fields fields) {
				fields.txManager.ShouldFail = true
			},
			wantErr: true,
		},
		{
			name: "find_for_update_error",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               100.50,
				},
			},
			setup: func(fields fields) {
				fields.accountRepo.EXPECT().FindForUpdate(gomock.Any(), []uint64{111, 222}).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "accounts_not_found",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               100.50,
				},
			},
			setup: func(fields fields) {
				// Return only one account (less than 2)
				accounts := []*entity.Account{
					{ID: 111, Balance: 1000.00},
				}
				fields.accountRepo.EXPECT().FindForUpdate(gomock.Any(), []uint64{111, 222}).Return(accounts, nil)
			},
			wantErr: true,
		},
		{
			name: "insufficient_balance",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               2000.00, // More than available balance
				},
			},
			setup: func(fields fields) {
				accounts := []*entity.Account{
					{ID: 111, Balance: 1000.00}, // Only 1000 available
					{ID: 222, Balance: 500.00},
				}
				fields.accountRepo.EXPECT().FindForUpdate(gomock.Any(), []uint64{111, 222}).Return(accounts, nil)
			},
			wantErr: true,
		},
		{
			name: "transaction_create_error",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               100.50,
				},
			},
			setup: func(fields fields) {
				accounts := []*entity.Account{
					{ID: 111, Balance: 1000.00},
					{ID: 222, Balance: 500.00},
				}
				fields.accountRepo.EXPECT().FindForUpdate(gomock.Any(), []uint64{111, 222}).Return(accounts, nil)
				fields.transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("transaction record creation failed"))
			},
			wantErr: true,
		},
		{
			name: "source_account_update_error",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               100.50,
				},
			},
			setup: func(fields fields) {
				accounts := []*entity.Account{
					{ID: 111, Balance: 1000.00},
					{ID: 222, Balance: 500.00},
				}
				fields.accountRepo.EXPECT().FindForUpdate(gomock.Any(), []uint64{111, 222}).Return(accounts, nil)
				fields.transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.Transaction{}, nil)
				// First update (source account) fails
				fields.accountRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("source account update failed"))
			},
			wantErr: true,
		},
		{
			name: "destination_account_update_error",
			args: args{
				ctx: &gin.Context{},
				req: dto.TransactionDTO{
					SourceAccountID:      111,
					DestinationAccountID: 222,
					Amount:               100.50,
				},
			},
			setup: func(fields fields) {
				accounts := []*entity.Account{
					{ID: 111, Balance: 1000.00},
					{ID: 222, Balance: 500.00},
				}
				fields.accountRepo.EXPECT().FindForUpdate(gomock.Any(), []uint64{111, 222}).Return(accounts, nil)
				fields.transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entity.Transaction{}, nil)
				// First update (source account) succeeds, second update (destination account) fails
				fields.accountRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				fields.accountRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("destination account update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAccountRepo := mock.NewMockAccountRepository(ctrl)
			mockTransactionRepo := mock.NewMockTransactionRepository(ctrl)
			mockTxManager := &mock2.MockTxManager{}

			testFields := fields{
				accountRepo:     mockAccountRepo,
				transactionRepo: mockTransactionRepo,
				txManager:       mockTxManager,
			}

			uc := accountUsecase{
				accountRepo:     mockAccountRepo,
				transactionRepo: mockTransactionRepo,
				txManager:       mockTxManager,
			}

			if tt.setup != nil {
				tt.setup(testFields)
			}

			err := uc.MakeTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
