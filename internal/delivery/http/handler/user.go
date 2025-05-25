package handler

import (
	"net/http"

	"github.com/FlyKarlik/gofemart/internal/delivery/http/response"
	"github.com/FlyKarlik/gofemart/internal/delivery/http/status"
	"github.com/FlyKarlik/gofemart/internal/errs"
	"github.com/FlyKarlik/gofemart/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/gin-gonic/gin"
)

// RegisterUser registers a new user
// @Summary User registration
// @Description Creates a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body model.UserInput true "Registration data"
// @Success 200 {object} response.BaseResponseAny "Successfully processed request"
// @Failure 400 {object} response.BaseResponseAny "Bad request - invalid input params"
// @Failure 409 {object} response.BaseResponseAny "Conflict - user login in use"
// @Failure 500 {object} response.BaseResponseAny "Internal system error"
// @Router /api/user/register [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	tracer := otel.Tracer("handler/register-user")
	ctx, span := tracer.Start(c.Request.Context(), "RegisterUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("handler", "RegisterUser"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.FullPath()),
	)

	var input model.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("handler[user]", "RegisterUser", "Failed to parse json object", err)
		response.New[any](c, http.StatusBadRequest, false, nil, errs.ErrInvalidRequest)
		return
	}

	if err := h.usecase.RegisterUser(ctx, input); err != nil {
		h.logger.Error("handler[user]", "RegisterUser", "Failed to register user", err)
		response.New[any](c, status.HTTPStatusFromError(err), false, nil, err)
		return
	}

	response.New[any](c, http.StatusOK, true, nil, nil)
}

type loginUserResponse struct {
	Token string `json:"token"`
}

// LoginUser authenticates a user and returns an access token
// @Summary Authenticate user
// @Description Verifies user credentials and returns JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body model.UserInput true "Login credentials"
// @Success 200 {object} response.BaseResponseLogin "Successful authentication"
// @Failure 400 {object} response.BaseResponseAny "Invalid request format"
// @Failure 401 {object} response.BaseResponseAny "Authentication failed"
// @Failure 500 {object} response.BaseResponseAny "Server error"
// @Router /api/user/login [post]
func (h *Handler) LoginUser(c *gin.Context) {
	tracer := otel.Tracer("handler/login-user")
	ctx, span := tracer.Start(c.Request.Context(), "LoginUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("handler", "LoginUser"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.FullPath()),
	)

	var input model.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("handler[user]", "LoginUser", "Failed to parse json object", err)
		response.New[any](c, http.StatusBadRequest, false, nil, errs.ErrInvalidRequest)
		return
	}

	token, err := h.usecase.LoginUser(ctx, input)
	if err != nil {
		h.logger.Error("handler[user]", "LoginUser", "Failed to login user", err)
		response.New[any](c, status.HTTPStatusFromError(err), false, nil, err)
		return
	}

	response.New(c, http.StatusOK, true, loginUserResponse{Token: token}, nil)
}

// CreateOrder uploads a new order number for processing
// @Summary Upload user order
// @Description Accepts a plain text order number and processes it
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order body model.UserOrderInput true "Order number"
// @Success 200 {object} response.BaseResponseAny "Order already uploaded by this user"
// @Success 202 {object} nil "New order accepted for processing"
// @Failure 400 {object} response.BaseResponseAny "Invalid request format"
// @Failure 401 {object} response.BaseResponseAny "Unauthorized"
// @Failure 409 {object} response.BaseResponseAny "Order already uploaded by another user"
// @Failure 422 {object} response.BaseResponseAny "Invalid order number format"
// @Failure 500 {object} response.BaseResponseAny "Internal server error"
// @Router /api/user/orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	tracer := otel.Tracer("handler/create-order")
	ctx, span := tracer.Start(c.Request.Context(), "CreateOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("handler", "CreateOrder"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.FullPath()),
	)

	var input model.UserOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("handler[user]", "CreateOrder", "Failed to parse JSON body", err)
		response.New[any](c, http.StatusBadRequest, false, nil, errs.ErrInvalidRequest)
		return
	}

	if !isValidOrderNumber(*input.Number) {
		h.logger.Error("handler[user]", "CreateOrder", "Failed to validate order number", errs.ErrInvalidOrderNumber)
		response.New[any](c, http.StatusUnprocessableEntity, false, nil, errs.ErrInvalidOrderNumber)
		return
	}

	err := h.usecase.CreateUserOrder(ctx, input)
	if err != nil {
		if status.CodeFromError(err) == errs.CodeOrderAlreadyUpload {
			response.New[any](c, http.StatusOK, true, nil, nil)
			return
		}
		h.logger.Error("handler[user]", "CreateOrder", "Failed to create order", err)
		response.New[any](c, status.HTTPStatusFromError(err), false, nil, err)
		return
	}

	response.New[any](c, http.StatusAccepted, true, nil, nil)
}

// GetUserOrders returns a list of user's uploaded orders
// @Summary Get user's orders
// @Description Retrieves all orders uploaded by the authenticated user
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponseOrders "Successful response with orders"
// @Success 204 {object} nil "No orders found for user"
// @Failure 401 {object} response.BaseResponseAny "Unauthorized"
// @Failure 500 {object} response.BaseResponseAny "Internal server error"
// @Router /api/user/orders [get]
func (h *Handler) GetUserOrders(c *gin.Context) {
	tracer := otel.Tracer("handler/get-user-orders")
	ctx, span := tracer.Start(c.Request.Context(), "GetUserOrders")
	defer span.End()

	span.SetAttributes(
		attribute.String("handler", "GetUserOrders"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.FullPath()),
	)

	orders, err := h.usecase.GetUserOrders(ctx)
	if err != nil {
		h.logger.Error("handler[user]", "GetUserOrders", "Failed to get user orders", err)
		response.New[any](c, status.HTTPStatusFromError(err), false, nil, err)
		return
	}

	response.New(c, http.StatusOK, true, orders, nil)
}

// GetUserBalance returns the current balance and total withdrawn amount
// @Summary Get user balance
// @Description Retrieves the current balance and total amount withdrawn by the user
// @Tags Balance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponseBalance "Successful response with balance data"
// @Failure 401 {object} response.BaseResponseAny "Unauthorized"
// @Failure 500 {object} response.BaseResponseAny "Internal server error"
// @Router /api/user/balance [get]
func (h *Handler) GetUserBalance(c *gin.Context) {
	tracer := otel.Tracer("handler/get-user-balance")
	ctx, span := tracer.Start(c.Request.Context(), "GetUserBalance")
	defer span.End()

	span.SetAttributes(
		attribute.String("handler", "GetUserBalance"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.FullPath()),
	)

	balance, err := h.usecase.GetUserBalance(ctx)
	if err != nil {
		h.logger.Error("handler[user]", "GetUserBalance", "Failed to get user balance", err)
		response.New[any](c, status.HTTPStatusFromError(err), false, nil, err)
		return
	}

	response.New(c, http.StatusOK, true, balance, nil)
}

// WithdrawUserBalance withdraws funds from the user's balance
// @Summary Withdraw user balance
// @Description Deducts the specified amount from the user's balance for the given order
// @Tags Balance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param withdrawal body response.WitdrawalUserBalanceInput true "Withdrawal request"
// @Success 200 {object} response.BaseResponseAny "Withdrawal successful"
// @Failure 400 {object} response.BaseResponseAny "Invalid request format"
// @Failure 401 {object} response.BaseResponseAny "Unauthorized"
// @Failure 402 {object} response.BaseResponseAny "Insufficient funds"
// @Failure 422 {object} response.BaseResponseAny "Invalid order number format"
// @Failure 500 {object} response.BaseResponseAny "Internal server error"
// @Router /api/user/balance/withdraw [post]
func (h *Handler) WithdrawUserBalance(c *gin.Context) {
	tracer := otel.Tracer("handler/withdraw-user-balance")
	ctx, span := tracer.Start(c.Request.Context(), "WithdrawUserBalance")
	defer span.End()

	span.SetAttributes(
		attribute.String("handler", "WithdrawUserBalance"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.FullPath()),
	)

	var input model.UserWithdrawalInput[float64]
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("handler[user]", "WithdrawUserBalance", "Failed to parse JSON body", err)
		response.New[any](c, http.StatusBadRequest, false, nil, errs.ErrInvalidRequest)
		return
	}

	if !isValidOrderNumber(*input.OrderNumber) {
		h.logger.Error("handler[user]", "WithdrawUserBalance", "Failed to validate order number", errs.ErrInvalidOrderNumber)
		response.New[any](c, http.StatusUnprocessableEntity, false, nil, errs.ErrInvalidOrderNumber)
		return
	}

	if err := h.usecase.WithdrawUserBalance(ctx, input); err != nil {
		h.logger.Error("handler[user]", "WithdrawUserBalance", "Failed to withdraw user balance", err)
		response.New[any](c, status.HTTPStatusFromError(err), false, nil, err)
		return
	}

	response.New[any](c, http.StatusOK, true, nil, nil)
}

// GetUserWithdrawals returns user's withdrawal history
// @Summary Get user withdrawals
// @Description Retrieves a list of all user's balance withdrawals
// @Tags Balance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponseWithdrawals "Successful response with withdrawal history"
// @Success 204 {object} nil "No withdrawals found"
// @Failure 401 {object} response.BaseResponseAny "Unauthorized"
// @Failure 500 {object} response.BaseResponseAny "Internal server error"
// @Router /api/user/withdrawals [get]
func (h *Handler) GetUserWithdrawals(c *gin.Context) {
	tracer := otel.Tracer("handler/get-user-withdrawals")
	ctx, span := tracer.Start(c.Request.Context(), "GetUserWithdrawals")
	defer span.End()

	span.SetAttributes(
		attribute.String("handler", "GetUserWithdrawals"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.FullPath()),
	)

	withdrawals, err := h.usecase.GetUserWithdrawals(ctx)
	if err != nil {
		if status.CodeFromError(err) == errs.CodeNooneWithdrawal {
			response.New[any](c, http.StatusNoContent, true, nil, nil)
			return
		}
		h.logger.Error("handler[user]", "WithdrawUserBalance", "Failed to withdraw user balance", err)
		response.New[any](c, status.HTTPStatusFromError(err), false, nil, err)
		return
	}

	response.New(c, http.StatusOK, true, withdrawals, nil)
}
