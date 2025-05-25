package response

import (
	"time"

	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BaseResponse[T any] struct {
	Status bool  `json:"status"`
	Code   int   `json:"code"`
	Data   T     `json:"data,omitempty"`
	Error  error `json:"error,omitempty"`
}

func New[T any](c *gin.Context, code int, status bool, data T, err error) {
	obj := &BaseResponse[T]{
		Code:   code,
		Status: status,
		Data:   data,
		Error:  err,
	}
	c.JSON(code, obj)
}

// For Swagger
type BaseResponseAny struct {
	Status bool        `json:"status"`
	Code   int         `json:"code"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type BaseResponseLogin struct {
	Status bool `json:"status"`
	Code   int  `json:"code"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type BaseResponseOrders struct {
	Status bool              `json:"status"`
	Code   int               `json:"code"`
	Data   []model.UserOrder `json:"data,omitempty"`
	Error  string            `json:"error,omitempty"`
}

type BaseResponseBalance struct {
	Status bool               `json:"status"`
	Code   int                `json:"code"`
	Data   UserBalanceBalance `json:"data,omitempty"`
	Error  string             `json:"error,omitempty"`
}

type UserBalanceBalance struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Current   *float64   `json:"current,omitempty"`
	Withdrawn *float64   `json:"withdrawn,omitempty"`
}

type BaseResponseWithdrawals struct {
	Status bool                 `json:"status"`
	Code   int                  `json:"code"`
	Data   []UserWithdrawalData `json:"data,omitempty"`
	Error  string               `json:"error,omitempty"`
}

type UserWithdrawalData struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	OrderNumber *string    `json:"order_number,omitempty"`
	Sum         *float64   `json:"sum,omitempty"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}

type WitdrawalUserBalanceInput struct {
	OrderNumber *string  `json:"order_number" binding:"required"`
	Sum         *float64 `json:"sum" binding:"required,gt=0"`
}
