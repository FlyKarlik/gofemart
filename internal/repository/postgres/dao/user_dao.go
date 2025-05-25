package dao

import (
	"database/sql"

	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/pkg/database/pghelpers"
	"github.com/google/uuid"
)

type UserDAO struct {
	ID        uuid.NullUUID
	Login     sql.NullString
	Password  sql.NullString
	CreatedAt sql.NullTime
}

func (u *UserDAO) ToModel() *model.User {
	return &model.User{
		ID:        pghelpers.FromNullUUID(u.ID),
		Login:     pghelpers.FromNullString(u.Login),
		Password:  pghelpers.FromNullString(u.Password),
		CreatedAt: pghelpers.FromNullTime(u.CreatedAt),
	}
}

type UserInputDAO struct {
	Login    sql.NullString
	Password sql.NullString
}

func (u *UserInputDAO) FromModel(m model.UserInput) UserInputDAO {
	return UserInputDAO{
		Login:    pghelpers.ToNullString(m.Login),
		Password: pghelpers.ToNullString(m.Password),
	}
}

type UserOrderDAO struct {
	ID         uuid.NullUUID
	UserID     uuid.NullUUID
	Number     sql.NullString
	Status     sql.NullString
	Accrual    sql.NullInt64
	UploadedAt sql.NullTime
}

func (o *UserOrderDAO) ToModel() *model.UserOrder {
	return &model.UserOrder{
		ID:         pghelpers.FromNullUUID(o.ID),
		UserID:     pghelpers.FromNullUUID(o.UserID),
		Number:     pghelpers.FromNullString(o.Number),
		Status:     (*model.OrderStatusEnum)(&o.Status.String),
		Accrual:    pghelpers.FromNullInt64(o.Accrual),
		UploadedAt: pghelpers.FromNullTime(o.UploadedAt),
	}
}

type UserOrderInputDAO struct {
	UserID uuid.NullUUID
	Number sql.NullString
	Status sql.NullString
}

func (o *UserOrderInputDAO) FromModel(input model.UserOrderInput) UserOrderInputDAO {
	return UserOrderInputDAO{
		UserID: pghelpers.ToNullUUID(input.UserID),
		Number: pghelpers.ToNullString(input.Number),
		Status: pghelpers.ToNullString((*string)(input.Status)),
	}
}

type UserBalanceDAO struct {
	UserID    uuid.NullUUID
	Current   sql.NullInt64
	Withdrawn sql.NullInt64
}

func (b *UserBalanceDAO) ToModel() *model.UserBalance[int64] {
	return &model.UserBalance[int64]{
		UserID:    pghelpers.FromNullUUID(b.UserID),
		Current:   pghelpers.FromNullInt64(b.Current),
		Withdrawn: pghelpers.FromNullInt64(b.Withdrawn),
	}
}

type UserWithdrawalDAO struct {
	ID          uuid.NullUUID
	UserID      uuid.NullUUID
	OrderNumber sql.NullString
	Sum         sql.NullInt64
	ProcessedAt sql.NullTime
}

func (w *UserWithdrawalDAO) ToModel() *model.UserWithdrawal[int64] {
	return &model.UserWithdrawal[int64]{
		ID:          pghelpers.FromNullUUID(w.ID),
		UserID:      pghelpers.FromNullUUID(w.UserID),
		OrderNumber: pghelpers.FromNullString(w.OrderNumber),
		Sum:         pghelpers.FromNullInt64(w.Sum),
		ProcessedAt: pghelpers.FromNullTime(w.ProcessedAt),
	}
}

type UserWithdrawalInputDAO struct {
	UserID      uuid.NullUUID
	OrderNumber sql.NullString
	Sum         sql.NullInt64
}

func (w *UserWithdrawalInputDAO) FromModel(input model.UserWithdrawalInput[int64]) UserWithdrawalInputDAO {
	return UserWithdrawalInputDAO{
		UserID:      pghelpers.ToNullUUID(input.UserID),
		OrderNumber: pghelpers.ToNullString(input.OrderNumber),
		Sum:         pghelpers.ToNullInt64(input.Sum),
	}
}
