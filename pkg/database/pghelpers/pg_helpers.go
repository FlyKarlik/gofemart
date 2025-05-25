package pghelpers

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func ToNullString(s *string) sql.NullString {
	if s != nil && len(*s) != 0 {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{String: "", Valid: false}
}

func FromNullString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func FromNullUUID(uuid uuid.NullUUID) *uuid.UUID {
	if uuid.Valid {
		return &uuid.UUID
	}
	return nil
}

func ToNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u != nil {
		return uuid.NullUUID{UUID: *u, Valid: true}
	}
	return uuid.NullUUID{}
}

func FromNullTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func ToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{}
}

func ToNullInt64(i *int64) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: *i, Valid: true}
	}
	return sql.NullInt64{}
}

func FromNullInt64(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}
