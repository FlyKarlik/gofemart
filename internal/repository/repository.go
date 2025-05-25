package repository

import (
	"context"
	"time"

	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/internal/repository/cache"
	"github.com/FlyKarlik/gofemart/internal/repository/postgres"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, input model.UserInput) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)

	CreateUserOrder(ctx context.Context, input model.UserOrderInput) (*model.UserOrder, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]model.UserOrder, error)
	CheckUserOrderExists(ctx context.Context, number string, userID uuid.UUID) (bool, error)

	GetUserBalance(ctx context.Context, userID uuid.UUID) (*model.UserBalance[int64], error)
	CreateUserWithdrawal(ctx context.Context, input model.UserWithdrawalInput[int64]) (*model.UserWithdrawal[int64], error)
	GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.UserWithdrawal[int64], error)
}

type IUserCache interface {
	Set(ctx context.Context, userID uuid.UUID, user *model.User, ttl time.Duration) error
	Get(ctx context.Context, userID uuid.UUID) (*model.User, bool, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type Repository struct {
	IUserRepository
	IUserCache
}

func New(logger logger.Logger, conn *pgxpool.Pool, redisClient *redis.Client) *Repository {
	return &Repository{
		IUserRepository: postgres.NewUserRepo(logger, conn),
		IUserCache:      cache.NewUserCache(logger, redisClient),
	}
}
