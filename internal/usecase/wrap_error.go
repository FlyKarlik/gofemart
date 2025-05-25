package usecase

import (
	"errors"

	"github.com/FlyKarlik/gofemart/internal/errs"
	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/pkg/database/pghelpers"
)

func wrapUsecaseError(
	event model.EventTypeEnum,
	err error,
) error {
	if err == nil {
		return nil
	}

	var pgErr *pghelpers.PgError
	if errors.As(err, &pgErr) {
		if eventErrors, ok := errorMapping[pgErr.Code]; ok {
			if mappedErr, ok := eventErrors[event]; ok {
				return mappedErr
			}
		}
		return errs.New(errs.CodeUnknown, "unknown error")
	}

	return errs.New(errs.CodeUnknown, "unknown error")
}
