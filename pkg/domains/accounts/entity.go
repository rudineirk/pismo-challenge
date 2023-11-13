package accounts

import "time"

type Account struct {
	ID             int64
	DocumentNumber string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
