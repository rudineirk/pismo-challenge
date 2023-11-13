package errorlib_test

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
	assert "github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	t.Run("should create a custom errors", func(t *testing.T) {
		customErr := errorlib.NewError("new_error_code", "new error code")
		otherErr := errorlib.NewError("other_error_code", "other error code")

		err := customErr(errors.New("wrapped error"))
		err2 := otherErr(nil)

		assert.True(t, errors.Is(err, customErr(nil)))
		assert.False(t, errors.Is(customErr(nil), errors.New("new error code")))
		assert.False(t, errors.Is(err2, customErr(nil)))
	})

	t.Run("should return full error string", func(t *testing.T) {
		dbErr := errors.New("unique constraint error")
		duplicatedErr := errorlib.ErrDuplicated(dbErr)

		assert.Equal(t, "duplicated entity\nunique constraint error", duplicatedErr.Error())
	})

	t.Run("should unwrap error", func(t *testing.T) {
		notFoundErr := errorlib.ErrNotFound(sql.ErrNoRows)

		assert.Equal(t, sql.ErrNoRows.Error(), notFoundErr.Unwrap().Error())
	})

	t.Run("should return stacktrace", func(t *testing.T) {
		err := errorlib.ErrDuplicated(nil)

		stacktrace := err.GetStack()
		assert.Greater(t, len(stacktrace), 0)

		includesTestFile := false
		for _, trace := range stacktrace {
			if strings.Contains(trace.File, "errors_test.go") {
				includesTestFile = true
				break
			}
		}

		assert.True(t, includesTestFile)
	})
}
