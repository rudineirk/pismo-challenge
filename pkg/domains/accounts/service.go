package accounts

import "context"

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
	return svc.repo.CreateAccount(ctx, account)
}

func (svc *accountsService) GetAccountByID(ctx context.Context, id int64) (*Account, error) {
	return svc.repo.GetAccountByID(ctx, id)
}
