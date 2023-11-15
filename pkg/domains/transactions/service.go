package transactions

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rudineirk/pismo-challenge/pkg/domains/accounts"
	"github.com/rudineirk/pismo-challenge/pkg/domains/operationtypes"
	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
	"github.com/shopspring/decimal"
)

var ErrAccountIDNotFound = errorlib.NewError( //nolint:gochecknoglobals // error maker
	"account_id_not_found",
	"account_id not found",
)
var ErrInvalidOperationTypeID = errorlib.NewError( //nolint:gochecknoglobals // error maker
	"invalid_operation_type_id",
	"invalid operation_type_id",
)
var ErrInvalidAmount = errorlib.NewError( //nolint:gochecknoglobals // error maker
	"invalid_amount",
	"invalid amount",
)

type Service interface {
	CreateTransaction(context.Context, *CreateTransactionRequest) (*Transaction, error)
}

type CreateTransactionRequest struct {
	AccountID       int64               `json:"account_id"        validate:"required"`
	OperationTypeID operationtypes.Type `json:"operation_type_id" validate:"required"`
	Amount          float64             `json:"amount"            validate:"required"`
}

type transactionsService struct {
	repo        Repository
	accountsSvc accounts.Service
	validate    *validator.Validate
}

func NewService(repo Repository, accountsSvc accounts.Service) Service {
	return &transactionsService{
		repo:        repo,
		accountsSvc: accountsSvc,
		validate:    validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (svc *transactionsService) CreateTransaction(
	ctx context.Context,
	req *CreateTransactionRequest,
) (*Transaction, error) {
	if err := svc.validate.Struct(req); err != nil {
		return nil, errorlib.ErrInvalidPayload(err)
	} else if !operationtypes.IsValidOperationType(req.OperationTypeID) {
		return nil, ErrInvalidOperationTypeID(nil)
	} else if svc.isValidAmount(req.Amount, req.OperationTypeID) {
		return nil, ErrInvalidAmount(nil)
	} else if err := svc.validateAccountID(ctx, req.AccountID); err != nil {
		return nil, ErrAccountIDNotFound(err)
	}

	account := &Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          req.Amount,
		EventDate:       time.Now(),
	}

	if err := svc.repo.CreateTransaction(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (svc *transactionsService) isValidAmount(amount float64, opType operationtypes.Type) bool {
	sign := -1
	if opType == operationtypes.PaymentType {
		sign = 1
	}

	decimalAmount := decimal.NewFromFloat(amount)
	hasTwoDecimalsPrecision := decimalAmount.
		Mul(decimal.NewFromInt(100)).
		Mod(decimal.NewFromInt(1)).
		IsZero()

	return !hasTwoDecimalsPrecision || decimalAmount.Sign() != sign
}

func (svc *transactionsService) validateAccountID(ctx context.Context, id int64) error {
	_, err := svc.accountsSvc.GetAccountByID(ctx, id)

	return err
}
