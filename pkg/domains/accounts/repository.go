package accounts

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
	"github.com/uptrace/bun"
)

type AccountModel struct {
	bun.BaseModel  `bun:"table:accounts"`
	ID             int64     `bun:"id,pk,autoincrement"`
	DocumentNumber string    `bun:"document_number"`
	CreatedAt      time.Time `bun:"created_at"`
	UpdatedAt      time.Time `bun:"updated_at"`
}

func NewModelFromEntity(account *Account) *AccountModel {
	return &AccountModel{
		ID:             account.ID,
		DocumentNumber: account.DocumentNumber,
		CreatedAt:      account.CreatedAt,
		UpdatedAt:      account.UpdatedAt,
	}
}

func (model *AccountModel) ToEntity() *Account {
	return &Account{
		ID:             model.ID,
		DocumentNumber: model.DocumentNumber,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
}

type Repository interface {
	CreateAccount(context.Context, *Account) error
	GetAccountByID(context.Context, int64) (*Account, error)
}

type dbRepository struct {
	bunDB *bun.DB
}

func NewRepository(bunDB *bun.DB) Repository {
	return &dbRepository{bunDB}
}

func (repo *dbRepository) CreateAccount(ctx context.Context, account *Account) error {
	accountModel := NewModelFromEntity(account)

	_, err := repo.bunDB.NewInsert().
		Model(accountModel).
		Exec(ctx)

	if err != nil && strings.Contains(err.Error(), "unique constraint") {
		return errorlib.ErrDuplicated(err)
	} else if err != nil {
		return err
	}

	account.ID = accountModel.ID

	return nil
}

func (repo *dbRepository) GetAccountByID(ctx context.Context, id int64) (*Account, error) {
	accountModel := AccountModel{}

	err := repo.bunDB.NewSelect().
		Model(&accountModel).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errorlib.ErrNotFound(err)
	} else if err != nil {
		return nil, err
	}

	return accountModel.ToEntity(), nil
}
