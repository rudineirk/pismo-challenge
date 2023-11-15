package transactions

import (
	"time"

	"github.com/rudineirk/pismo-challenge/pkg/domains/operationtypes"
)

type Transaction struct {
	ID              int64
	AccountID       int64
	OperationTypeID operationtypes.Type
	Amount          float64
	EventDate       time.Time
}
