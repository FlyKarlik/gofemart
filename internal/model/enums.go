package model

type EventTypeEnum string

const (
	EventTypeEnumRegisterUser        EventTypeEnum = "REGISTER_USER"
	EventTypeEnumLoginUser           EventTypeEnum = "LOGIN_USER"
	EventTypeEnumGetUserByID         EventTypeEnum = "GET_USER_BY_ID"
	EventTypeEnumCreateOrder         EventTypeEnum = "CREATE_ORDER"
	EventTypeEnumGetUserOrders       EventTypeEnum = "GET_USER_ORDERS"
	EventTypeEnumGetUserBalance      EventTypeEnum = "GET_USER_BALANCE"
	EventTypeEnumWithdrawUserBalance EventTypeEnum = "WITHDRAW_USER_BALANCE"
	EventTypeEnumGetUserWithdrawals  EventTypeEnum = "GET_USER_WITHDRAWALS"
)

type ContextKeyEnum string

const (
	ContextKeyEnumUserID ContextKeyEnum = "USER"
)

func (c ContextKeyEnum) String() string {
	return string(c)
}

type OrderStatusEnum string

const (
	OrderStatusEnumNew        OrderStatusEnum = "NEW"
	OrderStatusEnumProcessing OrderStatusEnum = "PROCESSING"
	OrderStatusEnumInvalid    OrderStatusEnum = "INVALID"
	OrderStatusEnumProcessed  OrderStatusEnum = "PROCESSED"
)

func (c OrderStatusEnum) String() string {
	return string(c)
}
