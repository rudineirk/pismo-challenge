package accounts

import (
	"bytes"
	"context"
	"unicode"

	"github.com/paemuri/brdoc/v2"
	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
)

var ErrInvalidDocumentNumber = errorlib.NewError( //nolint:gochecknoglobals // error maker
	"invalid_document_number",
	"invalid document number",
)

type Service interface {
	CreateAccount(context.Context, *Account) error
	GetAccountByID(context.Context, int64) (*Account, error)
}

type accountsService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &accountsService{repo: repo}
}

func (svc *accountsService) CreateAccount(ctx context.Context, account *Account) error {
	documentNumber := svc.cleanDocumentNumber(account.DocumentNumber)
	if err := svc.validateDocument(documentNumber); err != nil {
		return err
	}

	account.DocumentNumber = documentNumber

	return svc.repo.CreateAccount(ctx, account)
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
