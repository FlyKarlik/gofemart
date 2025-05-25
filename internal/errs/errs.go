package errs

import "fmt"

type CustomError struct {
	Code    CodeEnum `json:"code"`
	Message string   `json:"message"`
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("Code: %d, message: %s", c.Code, c.Message)
}

func New(code CodeEnum, msg string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: msg,
	}
}

type CodeEnum int

const (
	CodeUnknown = iota
	CodeLoginInUse
	CodeInvalidRequest
	CodeUserNotFound
	CodeInvalidLoginOrPassword
	CodeUnauthorized
	CodeEmptyAuthHeader
	CodeInvalidToken
	CodeOrderByAnotherUserUpload
	CodeOrderAlreadyUpload
	CodeNoOrders
	CodeInvalidOrderNumber
	CodeOrderDoesNotExists
	CodeNotEnoughBalance
	CodeNooneWithdrawal
)

var (
	ErrInvalidToken          = New(CodeInvalidToken, "invalid token")
	ErrEmptyAuthHeader       = New(CodeEmptyAuthHeader, "empty auth header")
	ErrInvalidRequest        = New(CodeInvalidRequest, "invalid request")
	ErrInvalidLoginOrPassord = New(CodeInvalidLoginOrPassword, "invalid login or password")
	ErrUnauthorized          = New(CodeUnauthorized, "unauthorized")
	ErrOrderAlreadyUpload    = New(CodeOrderAlreadyUpload, "your ordder already upload")
	ErrNoOrders              = New(CodeNoOrders, "you have not any uploaded orders")
	ErrInvalidOrderNumber    = New(CodeInvalidOrderNumber, "invalid order number")
	ErrOrderDoesNotExists    = New(CodeOrderDoesNotExists, "order does not exists")
	ErrNotEnoughBalance      = New(CodeNotEnoughBalance, "not enough balance")
	ErrNooneWithdrawal       = New(CodeNooneWithdrawal, "no one withdrawal")
)
