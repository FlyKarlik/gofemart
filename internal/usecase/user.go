package usecase

import (
	"context"
	"time"

	"github.com/FlyKarlik/gofemart/config"
	"github.com/FlyKarlik/gofemart/internal/errs"
	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/internal/repository"
	"github.com/FlyKarlik/gofemart/pkg/hash"
	"github.com/FlyKarlik/gofemart/pkg/jwt"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	"github.com/google/uuid"
)

type userUsecase struct {
	cfg       *config.Config
	logger    logger.Logger
	userCache repository.IUserCache
	userRepo  repository.IUserRepository
}

func newUserUsecase(cfg *config.Config, logger logger.Logger, userRepo repository.IUserRepository, userCache repository.IUserCache) *userUsecase {
	return &userUsecase{
		cfg:       cfg,
		logger:    logger,
		userCache: userCache,
		userRepo:  userRepo,
	}
}

func (u *userUsecase) RegisterUser(ctx context.Context, input model.UserInput) error {
	if input.Password != nil {
		hashedPass, err := hash.GenerateFromPassword(*input.Password)
		if err != nil {
			u.logger.Error("usecase[user]", "RegisterUser", "Failed to generate hash password", err)
			return wrapUsecaseError(model.EventTypeEnumRegisterUser, err)
		}
		input.Password = &hashedPass
	}

	_, err := u.userRepo.CreateUser(ctx, input)
	if err != nil {
		u.logger.Error("usecase[user]", "RegisterUser", "Failed to create user", err)
		return wrapUsecaseError(model.EventTypeEnumRegisterUser, err)
	}
	return nil
}

func (u *userUsecase) LoginUser(ctx context.Context, input model.UserInput) (string, error) {
	user, err := u.userRepo.GetUserByLogin(ctx, *input.Login)
	if err != nil {
		u.logger.Error("usecase[user]", "LoginUser", "Failed to get user", err)
		return "", wrapUsecaseError(model.EventTypeEnumLoginUser, err)
	}

	isVerified := func() bool {
		if err := hash.CompareHashAndPassword(*user.Password, *input.Password); err != nil {
			u.logger.Error("usecase[user]", "LoginUser", "Failed to compare hash and password", err)
			return false
		}

		if *user.Login != *input.Login {
			u.logger.Error("usecase[user]", "LoginUser", "Failed to compare db login with input login", errs.ErrInvalidLoginOrPassord)
			return false
		}
		return true
	}()

	if !isVerified {
		return "", errs.ErrInvalidLoginOrPassord
	}

	jwtPayload := jwt.JWTPayload{
		SecretKey:      u.cfg.AppGofemart.JWTSecret,
		Issuer:         u.cfg.AppGofemart.JWTIssuer,
		AccessTokenTTL: u.cfg.AppGofemart.JWTTokenTTL,
		Claims: jwt.Claims{
			UserID: user.ID.String(),
			Login:  *user.Login,
		},
	}

	accessToken, err := jwt.GenerateAccessToken(jwtPayload)
	if err != nil {
		u.logger.Error("usecase[user]", "LoginUser", "Failed to generate access token", err)
		return "", wrapUsecaseError(model.EventTypeEnumLoginUser, err)
	}

	return accessToken, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, found, err := u.userCache.Get(ctx, userID)
	if err != nil {
		u.logger.Error("usecase[user]", "GetUserByID", "Failed to get user from cache", err)
	}

	if found {
		u.logger.Debug("usecase[user]", "GetUserByID", "User found in cache", userID)
		return user, nil
	}

	user, err = u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		u.logger.Error("usecase[user]", "GetUserByID", "Failed to get user from database", err)
		return nil, wrapUsecaseError(model.EventTypeEnumGetUserByID, err)
	}

	if err := u.userCache.Set(ctx, userID, user, time.Minute*10); err != nil {
		u.logger.Warn("usecase[user]", "GetUserByID", "Failed to set user to cache", err)
	}

	return user, nil
}

func (u *userUsecase) CreateUserOrder(ctx context.Context, input model.UserOrderInput) error {
	userID := ctx.Value(model.ContextKeyEnumUserID).(uuid.UUID)
	orderExists, err := u.userRepo.CheckUserOrderExists(ctx, *input.Number, userID)
	if err != nil {
		u.logger.Error("usecase[user]", "CreateUserOrder", "Failed to check if order exists", err)
		return wrapUsecaseError(model.EventTypeEnumCreateOrder, err)
	}
	if orderExists {
		return errs.ErrOrderAlreadyUpload
	}

	input.UserID = &userID
	orderStatus := model.OrderStatusEnumNew
	input.Status = &orderStatus

	if _, err := u.userRepo.CreateUserOrder(ctx, input); err != nil {
		u.logger.Error("usecase[user]", "CreateUserOrder", "Failed to create user order", err)
		return wrapUsecaseError(model.EventTypeEnumCreateOrder, err)
	}

	return nil
}

func (u *userUsecase) GetUserOrders(ctx context.Context) ([]model.UserOrder, error) {
	userID := ctx.Value(model.ContextKeyEnumUserID).(uuid.UUID)

	orders, err := u.userRepo.GetUserOrders(ctx, userID)
	if err != nil {
		u.logger.Error("usecase[user]", "GetUserOrdersByUserID", "Failed to get user orders", err)
		return nil, wrapUsecaseError(model.EventTypeEnumGetUserOrders, err)
	}

	if len(orders) == 0 {
		return nil, errs.ErrNoOrders
	}

	return orders, nil
}

func (u *userUsecase) GetUserBalance(ctx context.Context) (*model.UserBalance[float64], error) {
	userID := ctx.Value(model.ContextKeyEnumUserID).(uuid.UUID)
	balance, err := u.userRepo.GetUserBalance(ctx, userID)
	if err != nil {
		u.logger.Error("usecase[user]", "GetUserBalance", "Failed to get user balance", err)
		return nil, wrapUsecaseError(model.EventTypeEnumGetUserBalance, err)
	}

	return &model.UserBalance[float64]{
		UserID:    balance.UserID,
		Current:   convertMoneyValueToFloat64(balance.Current),
		Withdrawn: convertMoneyValueToFloat64(balance.Withdrawn),
	}, nil
}

func (u *userUsecase) WithdrawUserBalance(ctx context.Context, input model.UserWithdrawalInput[float64]) error {
	userID := ctx.Value(model.ContextKeyEnumUserID).(uuid.UUID)
	isOrderExists, err := u.userRepo.CheckUserOrderExists(ctx, *input.OrderNumber, userID)
	if err != nil {
		u.logger.Error("usecase[user]", "WithdrawUserBalance", "Failed to check order exists", err)
		return wrapUsecaseError(model.EventTypeEnumWithdrawUserBalance, err)
	}

	if !isOrderExists {
		return errs.ErrOrderDoesNotExists
	}

	balance, err := u.userRepo.GetUserBalance(ctx, userID)
	if err != nil {
		u.logger.Error("usecase[user]", "WithdrawUserBalance", "Failed to get user balance", err)
		return wrapUsecaseError(model.EventTypeEnumWithdrawUserBalance, err)
	}

	if *balance.Current < *convertMoneyValueToInt64(input.Sum) {
		return errs.ErrNotEnoughBalance
	}

	prepareWitdrawal := model.UserWithdrawalInput[int64]{
		UserID: &userID,
		WitdrawnBalance: func() *int64 {
			withdrawnBalance := *balance.Withdrawn + *convertMoneyValueToInt64(input.Sum)
			return &withdrawnBalance
		}(),
		CurrentBalance: func() *int64 {
			currentBalance := *balance.Current - *convertMoneyValueToInt64(input.Sum)
			return &currentBalance
		}(),
		OrderNumber: input.OrderNumber,
		Sum:         convertMoneyValueToInt64(input.Sum),
	}

	if _, err := u.userRepo.CreateUserWithdrawal(ctx, prepareWitdrawal); err != nil {
		u.logger.Error("usecase", "WithdrawUserBalance", "Failed to create user withdrawal", err)
		return wrapUsecaseError(model.EventTypeEnumWithdrawUserBalance, err)
	}

	return nil
}

func (u *userUsecase) GetUserWithdrawals(ctx context.Context) ([]model.UserWithdrawal[float64], error) {
	userID := ctx.Value(model.ContextKeyEnumUserID).(uuid.UUID)
	withdrawals, err := u.userRepo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		u.logger.Error("usecase[user]", "WithdrawUserBalance", "Failed to get user balance", err)
		return nil, wrapUsecaseError(model.EventTypeEnumGetUserWithdrawals, err)
	}

	if len(withdrawals) == 0 {
		return nil, errs.ErrNooneWithdrawal
	}

	prepareWithdrawals := func() []model.UserWithdrawal[float64] {
		withdrawalsList := make([]model.UserWithdrawal[float64], len(withdrawals))
		for index, witdrawal := range withdrawals {
			userWithdrawal := model.UserWithdrawal[float64]{
				ID:          witdrawal.ID,
				UserID:      witdrawal.UserID,
				OrderNumber: witdrawal.OrderNumber,
				Sum:         convertMoneyValueToFloat64(witdrawal.Sum),
				ProcessedAt: witdrawal.ProcessedAt,
			}
			withdrawalsList[index] = userWithdrawal
		}
		return withdrawalsList
	}()

	return prepareWithdrawals, nil
}
