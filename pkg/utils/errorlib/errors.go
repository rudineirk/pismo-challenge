package errorlib

import (
	"runtime"
	"strings"

	"errors"
)

var (
	ErrNotFound       = NewError("not_found", "entity not found")      //nolint:gochecknoglobals // error maker
	ErrDuplicated     = NewError("duplicated", "duplicated entity")    //nolint:gochecknoglobals // error maker
	ErrInvalidPayload = NewError("invalid_payload", "invalid payload") //nolint:gochecknoglobals // error maker
)

type Error struct {
	Code       string          `json:"code"`
	Message    string          `json:"message"`
	err        error           `json:"-"`
	wrappedErr error           `json:"-"`
	stack      []runtime.Frame `json:"-"`
}

func NewError(code string, msg string) func(err error) *Error {
	return func(err error) *Error {
		return &Error{
			Code:       code,
			Message:    msg,
			wrappedErr: err,
			err:        errors.Join(errors.New(msg), err), //nolint:goerr113 // used to make constant/global errors
			stack:      getStack(),
		}
	}
}

func (err *Error) Error() string {
	return err.err.Error()
}

func (err *Error) Unwrap() error {
	return err.wrappedErr
}

func (err *Error) GetStack() []runtime.Frame {
	return err.stack
}

func (err *Error) Is(other error) bool {
	otherErr, ok := other.(*Error)
	if !ok {
		return false
	}

	return err.Code == otherErr.Code
}

func getStack() []runtime.Frame {
	pcs := make([]uintptr, 25)
	pcs = pcs[:runtime.Callers(3, pcs)]
	frames := runtime.CallersFrames(pcs)
	stacktrace := []runtime.Frame{}

	for {
		frame, more := frames.Next()

		if !more {
			break
		}

		if strings.Contains(frame.File, "vendor") && !strings.Contains(frame.File, "pismo-challenge") {
			continue
		}

		stacktrace = append(stacktrace, frame)
	}

	return stacktrace
}
