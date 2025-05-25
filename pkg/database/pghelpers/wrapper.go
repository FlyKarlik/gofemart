package pghelpers

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func WrapError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return newPgError(ErrUniqueViolation, "unique violation error", err)
		case "23503":
			return newPgError(ErrForeignKey, "foreign key violation error", err)
		case "23502":
			return newPgError(ErrNotNullViolation, "not null violation error", err)
		case "23514":
			return newPgError(ErrCheckViolation, "check violation error", err)
		case "40001":
			return newPgError(ErrSerialization, "serialization error", err)
		case "40P01":
			return newPgError(ErrDeadlock, "deadlock", err)
		default:
			return newPgError(ErrUndefined, "undefined error", err)
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return newPgError(ErrNoRows, "no rows in resul set", err)
	}

	return newPgError(ErrUnknown, "unknown db error", err)
}
