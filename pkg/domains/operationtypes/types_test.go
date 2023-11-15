package operationtypes_test

import (
	"testing"

	"github.com/rudineirk/pismo-challenge/pkg/domains/operationtypes"
	assert "github.com/stretchr/testify/require"
)

func TestValidOperations(t *testing.T) {
	t.Run("should return true for valid operation types", func(t *testing.T) {
		operationTypes := []operationtypes.Type{
			operationtypes.CashPurchaseType,
			operationtypes.InstallmentType,
			operationtypes.WithdrawType,
			operationtypes.PaymentType,
		}

		for _, opType := range operationTypes {
			assert.True(t, operationtypes.IsValidOperationType(opType))
		}
	})

	t.Run("should return false for invalid operation type", func(t *testing.T) {
		opType := operationtypes.Type(5)
		assert.False(t, operationtypes.IsValidOperationType(opType))
	})
}
