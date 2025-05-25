package usecase

import (
	"github.com/FlyKarlik/gofemart/internal/errs"
	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/pkg/database/pghelpers"
)

var errorMapping = map[pghelpers.PgErrorCode]map[model.EventTypeEnum]*errs.CustomError{
	pghelpers.ErrUniqueViolation: {
		model.EventTypeEnumRegisterUser: errs.New(errs.CodeLoginInUse, "login already exists"),
		model.EventTypeEnumCreateOrder:  errs.New(errs.CodeOrderByAnotherUserUpload, "current order already uploaded by another user"),
	},
	pghelpers.ErrNoRows: {
		model.EventTypeEnumLoginUser:   errs.New(errs.CodeUserNotFound, "user not found"),
		model.EventTypeEnumGetUserByID: errs.New(errs.CodeUnauthorized, "user creds not valid"),
	},
}
