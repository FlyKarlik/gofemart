package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        *uuid.UUID
	Login     *string
	Password  *string
	CreatedAt *time.Time
}

type UserInput struct {
	Login    *string `json:"login" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

type UserOrder struct {
	ID         *uuid.UUID       `json:"id,omitempty"`
	UserID     *uuid.UUID       `json:"user_id,omitempty"`
	Number     *string          `json:"number,omitempty"`
	Status     *OrderStatusEnum `json:"status,omitempty"`
	Accrual    *int64           `json:"accural,omitempty"`
	UploadedAt *time.Time       `json:"uploaded_at,omitempty"`
}

type UserOrderInput struct {
	UserID *uuid.UUID
	Number *string `json:"number" binding:"required"`
	Status *OrderStatusEnum
}

type UserBalance[T int64 | float64] struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Current   *T         `json:"current,omitempty"`
	Withdrawn *T         `json:"withdrawn,omitempty"`
}

type UserWithdrawal[T int64 | float64] struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	OrderNumber *string    `json:"order_number,omitempty"`
	Sum         *T         `json:"sum,omitempty"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}

type UserWithdrawalInput[T int64 | float64] struct {
	UserID          *uuid.UUID
	WitdrawnBalance *T
	CurrentBalance  *T
	OrderNumber     *string `json:"order_number" binding:"required"`
	Sum             *T      `json:"sum" binding:"required,gt=0"`
}
