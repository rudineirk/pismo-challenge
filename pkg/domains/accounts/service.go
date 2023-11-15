package accounts

import (
	"bytes"
	"context"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/paemuri/brdoc/v2"
	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
)

var ErrInvalidDocumentNumber = errorlib.NewError( //nolint:gochecknoglobals // error maker
	"invalid_document_number",
	"invalid document number",
)

type Service interface {
	CreateAccount(context.Context, *CreateAccountRequest) (*Account, error)
	GetAccountByID(context.Context, int64) (*Account, error)
}

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type accountsService struct {
	repo     Repository
	validate *validator.Validate
}

func NewService(repo Repository) Service {
	return &accountsService{
		repo:     repo,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (svc *accountsService) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*Account, error) {
	if err := svc.validate.Struct(req); err != nil {
		return nil, errorlib.ErrInvalidPayload(err)
	}

	documentNumber := svc.cleanDocumentNumber(req.DocumentNumber)
	if err := svc.validateDocument(documentNumber); err != nil {
		return nil, err
	}

	account := &Account{
		DocumentNumber: documentNumber,
		CreatedAt:      time.Now(),
	}

	account.UpdatedAt = account.CreatedAt

	if err := svc.repo.CreateAccount(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (svc *accountsService) GetAccountByID(ctx context.Context, id int64) (*Account, error) {
	return svc.repo.GetAccountByID(ctx, id)
}

func (svc *accountsService) cleanDocumentNumber(documentNumber string) string {
	buf := bytes.NewBufferString("")

	for _, r := range documentNumber {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

func (svc *accountsService) validateDocument(documentNumber string) error {
	if brdoc.IsCPF(documentNumber) || brdoc.IsCNPJ(documentNumber) {
		return nil
	}

	return ErrInvalidDocumentNumber(nil)
}
