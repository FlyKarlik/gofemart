package usecase

import (
	"context"

	"github.com/FlyKarlik/gofemart/config"
	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/internal/repository"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	"github.com/google/uuid"
)

type IUserUsecase interface {
	RegisterUser(ctx context.Context, input model.UserInput) error
	LoginUser(ctx context.Context, input model.UserInput) (string, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)

	CreateUserOrder(ctx context.Context, input model.UserOrderInput) error
	GetUserOrders(ctx context.Context) ([]model.UserOrder, error)

	GetUserBalance(ctx context.Context) (*model.UserBalance[float64], error)
	WithdrawUserBalance(ctx context.Context, input model.UserWithdrawalInput[float64]) error
	GetUserWithdrawals(ctx context.Context) ([]model.UserWithdrawal[float64], error)
}

type Usecase struct {
	IUserUsecase
}

func New(cfg *config.Config, logger logger.Logger, repo *repository.Repository) *Usecase {
	return &Usecase{
		IUserUsecase: newUserUsecase(cfg, logger, repo.IUserRepository, repo.IUserCache),
	}
}
