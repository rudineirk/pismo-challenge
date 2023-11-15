package transactions_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rudineirk/pismo-challenge/pkg/domains/accounts"
	accountMocks "github.com/rudineirk/pismo-challenge/pkg/domains/accounts/mocks"
	"github.com/rudineirk/pismo-challenge/pkg/domains/operationtypes"
	"github.com/rudineirk/pismo-challenge/pkg/domains/transactions"
	mocks "github.com/rudineirk/pismo-challenge/pkg/domains/transactions/mocks"
	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
	assert "github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestCreateTransaction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repo := mocks.NewMockRepository(mockCtrl)
	accountsSvc := accountMocks.NewMockService(mockCtrl)

	svc := transactions.NewService(repo, accountsSvc)

	for _, req := range []*transactions.CreateTransactionRequest{
		{AccountID: 1, OperationTypeID: operationtypes.CashPurchaseType, Amount: -0.01},
		{AccountID: 1, OperationTypeID: operationtypes.CashPurchaseType, Amount: -1000},
		{AccountID: 1, OperationTypeID: operationtypes.InstallmentType, Amount: -0.01},
		{AccountID: 1, OperationTypeID: operationtypes.InstallmentType, Amount: -1000},
		{AccountID: 1, OperationTypeID: operationtypes.WithdrawType, Amount: -0.01},
		{AccountID: 1, OperationTypeID: operationtypes.WithdrawType, Amount: -1000},
		{AccountID: 1, OperationTypeID: operationtypes.PaymentType, Amount: 0.01},
		{AccountID: 1, OperationTypeID: operationtypes.PaymentType, Amount: 1000},
	} {
		t.Run("should create a new transaction", func(t *testing.T) {
			repo.EXPECT().
				CreateTransaction(gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, transaction *transactions.Transaction) {
					transaction.ID = 1
				}).
				Return(nil)

			accountsSvc.EXPECT().
				GetAccountByID(gomock.Any(), gomock.Any()).
				Return(&accounts.Account{ID: 1}, nil)

			ctx := context.TODO()
			now := time.Now()

			transaction, err := svc.CreateTransaction(ctx, req)
			assert.NoError(t, err)

			assert.Equal(t, int64(1), transaction.ID)
			assert.Equal(t, req.AccountID, transaction.AccountID)
			assert.Equal(t, req.OperationTypeID, transaction.OperationTypeID)
			assert.Equal(t, req.Amount, transaction.Amount)
			assert.WithinDuration(t, now, transaction.EventDate, 5*time.Millisecond)
		})
	}
}

func TestCreateAccountInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repo := mocks.NewMockRepository(mockCtrl)
	accountsSvc := accountMocks.NewMockService(mockCtrl)

	svc := transactions.NewService(repo, accountsSvc)

	for _, req := range []*transactions.CreateTransactionRequest{
		{},
		{OperationTypeID: operationtypes.PaymentType, Amount: 1.15},
		{AccountID: 1, Amount: 1.15},
		{Amount: 1.15},
		{AccountID: 1, OperationTypeID: operationtypes.PaymentType, Amount: 0},
		{AccountID: 1, OperationTypeID: operationtypes.PaymentType, Amount: -0},
		{AccountID: 1, OperationTypeID: operationtypes.PaymentType},
	} {
		t.Run("should validate required field", func(t *testing.T) {
			ctx := context.TODO()

			_, err := svc.CreateTransaction(ctx, req)
			assert.ErrorIs(t, err, errorlib.ErrInvalidPayload(nil))
		})
	}

	t.Run("should return error if transaction type is invalid", func(t *testing.T) {
		ctx := context.TODO()

		_, err := svc.CreateTransaction(ctx, &transactions.CreateTransactionRequest{
			AccountID:       1,
			OperationTypeID: operationtypes.Type(789),
			Amount:          1.15,
		})
		assert.ErrorIs(t, err, transactions.ErrInvalidOperationTypeID(nil))
	})

	for _, req := range []*transactions.CreateTransactionRequest{
		{AccountID: 1, OperationTypeID: operationtypes.CashPurchaseType, Amount: 0.01},
		{AccountID: 1, OperationTypeID: operationtypes.CashPurchaseType, Amount: -1.019},
		{AccountID: 1, OperationTypeID: operationtypes.InstallmentType, Amount: 0.01},
		{AccountID: 1, OperationTypeID: operationtypes.WithdrawType, Amount: 0.01},
		{AccountID: 1, OperationTypeID: operationtypes.PaymentType, Amount: -0.01},
		{AccountID: 1, OperationTypeID: operationtypes.PaymentType, Amount: 1.019},
	} {
		t.Run("should return error if amount is invalid", func(t *testing.T) {
			ctx := context.TODO()

			_, err := svc.CreateTransaction(ctx, req)
			assert.ErrorIs(t, err, transactions.ErrInvalidAmount(nil))
		})
	}

	t.Run("should return error if account is not found", func(t *testing.T) {
		accountsSvc.EXPECT().
			GetAccountByID(gomock.Any(), gomock.Any()).
			Return(nil, errorlib.ErrNotFound(nil))

		ctx := context.TODO()

		_, err := svc.CreateTransaction(ctx, &transactions.CreateTransactionRequest{
			AccountID:       1,
			OperationTypeID: operationtypes.PaymentType,
			Amount:          1.15,
		})
		assert.ErrorIs(t, err, transactions.ErrAccountIDNotFound(nil))
	})

	t.Run("should return error from repo create transaction", func(t *testing.T) {
		dbTimeoutErr := errors.New("db timeout error")

		repo.EXPECT().
			CreateTransaction(gomock.Any(), gomock.Any()).
			Return(dbTimeoutErr)

		accountsSvc.EXPECT().
			GetAccountByID(gomock.Any(), gomock.Any()).
			Return(&accounts.Account{ID: 1}, nil)

		ctx := context.TODO()

		transaction, err := svc.CreateTransaction(ctx, &transactions.CreateTransactionRequest{
			AccountID:       1,
			OperationTypeID: operationtypes.PaymentType,
			Amount:          1.15,
		})
		assert.ErrorIs(t, err, dbTimeoutErr)
		assert.Nil(t, transaction)
	})
}
