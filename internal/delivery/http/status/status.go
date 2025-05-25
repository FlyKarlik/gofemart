package status

import (
	"net/http"

	"github.com/FlyKarlik/gofemart/internal/errs"
)

func HTTPStatusFromError(err error) int {
	if customErr, ok := err.(*errs.CustomError); ok {
		switch customErr.Code {
		case errs.CodeLoginInUse:
			return http.StatusConflict
		case errs.CodeInvalidLoginOrPassword:
			return http.StatusUnauthorized
		case errs.CodeUserNotFound:
			return http.StatusBadRequest
		case errs.CodeInvalidRequest:
			return http.StatusBadRequest
		case errs.CodeOrderByAnotherUserUpload:
			return http.StatusConflict
		case errs.CodeNoOrders:
			return http.StatusNoContent
		case errs.CodeNotEnoughBalance:
			return http.StatusPaymentRequired
		case errs.CodeOrderDoesNotExists:
			return http.StatusUnprocessableEntity
		default:
			return http.StatusInternalServerError
		}
	}

	return http.StatusInternalServerError
}

func CodeFromError(err error) errs.CodeEnum {
	if customErr, ok := err.(*errs.CustomError); ok {
		return customErr.Code
	}
	return -1
}
