package transactions

import (
	"context"
	"time"

	"github.com/rudineirk/pismo-challenge/pkg/domains/operationtypes"
	"github.com/uptrace/bun"
)

type Repository interface {
	CreateTransaction(context.Context, *Transaction) error
}

type TransactionModel struct {
	bun.BaseModel   `bun:"table:transactions"`
	ID              int64               `bun:"id,pk,autoincrement"`
	AccountID       int64               `bun:"account_id"`
	OperationTypeID operationtypes.Type `bun:"operation_type_id"`
	Amount          float64             `bun:"amount"`
	EventDate       time.Time           `bun:"event_date"`
}

func NewModelFromEntity(transaction *Transaction) *TransactionModel {
	return &TransactionModel{
		ID:              transaction.ID,
		AccountID:       transaction.AccountID,
		OperationTypeID: transaction.OperationTypeID,
		Amount:          transaction.Amount,
		EventDate:       transaction.EventDate,
	}
}

type dbRepository struct {
	bunDB *bun.DB
}

func NewRepository(bunDB *bun.DB) Repository {
	return &dbRepository{bunDB}
}

func (repo *dbRepository) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	transactionModel := NewModelFromEntity(transaction)

	_, err := repo.bunDB.NewInsert().
		Model(transactionModel).
		Exec(ctx)

	if err != nil {
		return err
	}

	transaction.ID = transactionModel.ID

	return nil
}
