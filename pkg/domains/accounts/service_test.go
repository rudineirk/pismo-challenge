package accounts_test

import (
	"context"
	"testing"
	"time"

	"github.com/rudineirk/pismo-challenge/pkg/domains/accounts"
	mocks "github.com/rudineirk/pismo-challenge/pkg/domains/accounts/mocks"
	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
	assert "github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repo := mocks.NewMockRepository(mockCtrl)
	svc := accounts.NewService(repo)

	for _, data := range [][]string{
		{"23383829006", "233.838.290-06"},
		{"23383829006", "23383829006"},
		{"05677940000133", "05.677.940/0001-33"},
		{"05677940000133", "05677940000133"},
		{"05677940000133", "05.677.940.000133"},
	} {
		documentNumber := data[0]
		input := data[1]

		repo.EXPECT().
			CreateAccount(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, account *accounts.Account) {
				account.ID = 1
			}).
			Return(nil)

		t.Run("should create a new account", func(t *testing.T) {
			ctx := context.TODO()
			now := time.Now()
			req := &accounts.CreateAccountRequest{
				DocumentNumber: input,
			}

			account, err := svc.CreateAccount(ctx, req)
			assert.Nil(t, err)

			assert.Equal(t, int64(1), account.ID)
			assert.Equal(t, documentNumber, account.DocumentNumber)
			assert.WithinDuration(t, now, account.CreatedAt, 5*time.Millisecond)
			assert.WithinDuration(t, now, account.UpdatedAt, 5*time.Millisecond)
		})
	}
}

func TestCreateAccountInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repo := mocks.NewMockRepository(mockCtrl)
	svc := accounts.NewService(repo)

	for _, documentNumber := range []string{"23383829007", "05677940000134", "abc123", "2338382900"} {
		t.Run("should return error if document is invalid", func(t *testing.T) {
			ctx := context.TODO()
			req := &accounts.CreateAccountRequest{
				DocumentNumber: documentNumber,
			}

			_, err := svc.CreateAccount(ctx, req)
			assert.ErrorIs(t, err, accounts.ErrInvalidDocumentNumber(nil))

			repo.EXPECT().CreateAccount(nil, nil).Times(0)
		})
	}

	t.Run("should return error if document is not present", func(t *testing.T) {
		ctx := context.TODO()
		req := &accounts.CreateAccountRequest{
			DocumentNumber: "",
		}

		_, err := svc.CreateAccount(ctx, req)
		assert.ErrorIs(t, err, errorlib.ErrInvalidPayload(nil))

		repo.EXPECT().CreateAccount(nil, nil).Times(0)
	})
}

func TestGetAccountByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repo := mocks.NewMockRepository(mockCtrl)
	svc := accounts.NewService(repo)

	ctx := context.TODO()
	now := time.Now()
	documentNumber := "23383829006"
	account := &accounts.Account{
		ID:             1,
		DocumentNumber: documentNumber,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	repo.EXPECT().GetAccountByID(ctx, account.ID).
		Return(account, nil)

	result, err := svc.GetAccountByID(ctx, account.ID)
	assert.Nil(t, err)

	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, documentNumber, result.DocumentNumber)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, now, result.UpdatedAt)
}
