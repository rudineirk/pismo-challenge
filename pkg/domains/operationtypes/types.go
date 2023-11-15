package operationtypes

type Type int

const (
	CashPurchaseType Type = 1
	InstallmentType  Type = 2
	WithdrawType     Type = 3
	PaymentType      Type = 4
)

func IsValidOperationType(id Type) bool {
	return id >= CashPurchaseType && id <= PaymentType
}
