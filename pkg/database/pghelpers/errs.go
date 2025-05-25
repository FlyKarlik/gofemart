package pghelpers

import "fmt"

type PgError struct {
	Code    PgErrorCode
	Message string
	RawErr  error
}

func (r *PgError) Error() string {
	return fmt.Sprintf("Code: %d, message: %s, raw error: %d", r.Code, r.Message, r.RawErr)
}

func newPgError(code PgErrorCode, msg string, rawErr error) *PgError {
	return &PgError{
		Code:    code,
		Message: msg,
		RawErr:  rawErr,
	}
}

type PgErrorCode int

const (
	ErrUniqueViolation = iota
	ErrForeignKey
	ErrNoRows
	ErrCheckViolation
	ErrNotNullViolation
	ErrSerialization
	ErrDeadlock
	ErrUnknown
	ErrUndefined
)
