package postgres

import (
	"context"
	"errors"

	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/internal/repository/postgres/dao"
	"github.com/FlyKarlik/gofemart/internal/repository/postgres/quries"
	"github.com/FlyKarlik/gofemart/pkg/database/pghelpers"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	logger logger.Logger
	c      *pgxpool.Pool
}

func NewUserRepo(logger logger.Logger, conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		logger: logger,
		c:      conn,
	}
}

func (u *UserRepo) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	userInputDAO := new(dao.UserInputDAO).FromModel(input)

	tx, err := u.c.Begin(ctx)
	if err != nil {
		u.logger.Error("postgres[user]", "CreateUser", "Failed to begin transaction", err)
		return nil, pghelpers.WrapError(err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				u.logger.Error("postgres[user]", "CreateUser", "Failed to rollback transaction", err)
			}
		}
	}()

	query, args, err := quries.BuildCreateUserQuery(userInputDAO)
	if err != nil {
		u.logger.Error("postgres[user]", "CreateUser", "Failed to build create user query", err)
		return nil, pghelpers.WrapError(err)
	}

	row := tx.QueryRow(ctx, query, args...)

	var userDAO dao.UserDAO
	if err := row.Scan(
		&userDAO.ID,
		&userDAO.Login,
		&userDAO.Password,
		&userDAO.CreatedAt,
	); err != nil {
		u.logger.Error("postgres[user]", "CreateUser", "Failed to scan user row", err)
		return nil, pghelpers.WrapError(err)
	}

	balanceQuery, balanceArgs, err := quries.BuildCreateUserBalanceQuery(userDAO.ID)
	if err != nil {
		u.logger.Error("postgres[user]", "CreateUser", "Failed to build create balance query", err)
		return nil, pghelpers.WrapError(err)
	}

	if _, err := tx.Exec(ctx, balanceQuery, balanceArgs...); err != nil {
		u.logger.Error("postgres[user]", "CreateUser", "Failed to exec balance insert", err)
		return nil, pghelpers.WrapError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		u.logger.Error("postgres[user]", "CreateUser", "Failed to commit transaction", err)
		return nil, pghelpers.WrapError(err)
	}

	return userDAO.ToModel(), nil
}

func (u *UserRepo) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	query, args, err := quries.BuildGetUserByLoginQuery(pghelpers.ToNullString(&login))
	if err != nil {
		u.logger.Error("postgres[user]", "GetUserByLogin", "Failed to build query get user by login", err)
		return nil, pghelpers.WrapError(err)
	}

	row := u.c.QueryRow(ctx, query, args...)

	var userDAO dao.UserDAO
	if err := row.Scan(
		&userDAO.ID,
		&userDAO.Login,
		&userDAO.Password,
		&userDAO.CreatedAt,
	); err != nil {
		u.logger.Error("postgres[user]", "GetUserByLogin", "Failed to scan row", err)
		return nil, pghelpers.WrapError(err)
	}

	return userDAO.ToModel(), nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	query, args, err := quries.BuildGetUserByIDQuery(pghelpers.ToNullUUID(&userID))
	if err != nil {
		u.logger.Error("postgres[user]", "GetUserByID", "Failed to build query get user by ID", err)
		return nil, pghelpers.WrapError(err)
	}

	row := u.c.QueryRow(ctx, query, args...)

	var userDAO dao.UserDAO
	if err := row.Scan(
		&userDAO.ID,
		&userDAO.Login,
		&userDAO.Password,
		&userDAO.CreatedAt,
	); err != nil {
		u.logger.Error("postgres[user]", "GetUserByID", "Failed to scan row", err)
		return nil, pghelpers.WrapError(err)
	}

	return userDAO.ToModel(), nil
}

func (u *UserRepo) CreateUserOrder(ctx context.Context, input model.UserOrderInput) (*model.UserOrder, error) {
	userOrderInputDAO := new(dao.UserOrderInputDAO).FromModel(input)
	query, args, err := quries.BuildCreateOrderQuery(userOrderInputDAO)
	if err != nil {
		u.logger.Error("postgres[user]", "CreateUserOrder", "Failed to build query create user order", err)
		return nil, pghelpers.WrapError(err)
	}

	row := u.c.QueryRow(ctx, query, args...)

	var userOrderDAO dao.UserOrderDAO
	if err := row.Scan(
		&userOrderDAO.ID,
		&userOrderDAO.UserID,
		&userOrderDAO.Number,
		&userOrderDAO.Status,
		&userOrderDAO.Accrual,
		&userOrderDAO.UploadedAt,
	); err != nil {
		u.logger.Error("postgres[user]", "GetUserByID", "Failed to scan row", err)
		return nil, pghelpers.WrapError(err)
	}

	return userOrderDAO.ToModel(), nil
}

func (u *UserRepo) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]model.UserOrder, error) {
	query, args, err := quries.BuildGetUserOrdersQuery(pghelpers.ToNullUUID(&userID))
	if err != nil {
		u.logger.Error("postgres[user]", "GetUserOrdersByUserID", "Failed to build query get user orders by user id", err)
		return nil, pghelpers.WrapError(err)
	}

	rows, err := u.c.Query(ctx, query, args...)
	if err != nil {
		u.logger.Error("postgres[user]", "GetUserOrdersByUserID", "Failed to query user orders rows", err)
		return nil, pghelpers.WrapError(err)
	}
	defer rows.Close()

	var orders []model.UserOrder
	for rows.Next() {
		var order dao.UserOrderDAO
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			u.logger.Error("postgres[user]", "GetUserOrdersByUserID", "Failed to scan row", err)
			return nil, pghelpers.WrapError(err)
		}
		orders = append(orders, *order.ToModel())
	}

	if err := rows.Err(); err != nil {
		u.logger.Error("postgres[user]", "GetUserOrdersByUserID", "Rows error", err)
		return nil, pghelpers.WrapError(err)
	}

	return orders, nil
}

func (u *UserRepo) CheckUserOrderExists(ctx context.Context, number string, userID uuid.UUID) (bool, error) {
	query, args, err := quries.BuildCheckOrderExistsQuery(
		pghelpers.ToNullString(&number),
		pghelpers.ToNullUUID(&userID),
	)
	if err != nil {
		u.logger.Error("postgres[user]", "CheckOrderExists", "Failed to build query", err)
		return false, pghelpers.WrapError(err)
	}

	var exists int
	err = u.c.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		u.logger.Error("postgres[user]", "CheckOrderExists", "Query failed", err)
		return false, pghelpers.WrapError(err)
	}

	return true, nil
}

func (u *UserRepo) GetUserBalance(ctx context.Context, userID uuid.UUID) (*model.UserBalance[int64], error) {
	query, args, err := quries.BuildGetUserBalanceQuery(pghelpers.ToNullUUID(&userID))
	if err != nil {
		u.logger.Error("postgres[user]", "GetUserBalance", "Failed to build query", err)
		return nil, pghelpers.WrapError(err)
	}

	var userBalanceDAO dao.UserBalanceDAO
	row := u.c.QueryRow(ctx, query, args...)
	if err := row.Scan(
		&userBalanceDAO.UserID,
		&userBalanceDAO.Current,
		&userBalanceDAO.Withdrawn,
	); err != nil {
		u.logger.Error("postgres[user]", "GetUserBalance", "Query failed", err)
		return nil, pghelpers.WrapError(err)
	}

	return userBalanceDAO.ToModel(), nil
}

func (u *UserRepo) CreateUserWithdrawal(ctx context.Context, input model.UserWithdrawalInput[int64]) (*model.UserWithdrawal[int64], error) {
	tx, err := u.c.Begin(ctx)
	if err != nil {
		u.logger.Error("postgres[user]", "CreateUserWithdrawal", "Failed to begin transaction", err)
		return nil, pghelpers.WrapError(err)
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				u.logger.Error("postgres[user]", "CreateUserWithdrawal", "Failed to rollback transaction", err)
			}
		}
	}()

	withdrawalDAO := new(dao.UserWithdrawalInputDAO).FromModel(input)

	query, args, err := quries.BuildInsertWithdrawalQuery(withdrawalDAO)
	if err != nil {
		u.logger.Error("postgres[user]", "CreateUserWithdrawal", "Failed to build insert query", err)
		return nil, pghelpers.WrapError(err)
	}

	row := tx.QueryRow(ctx, query, args...)

	var resultDAO dao.UserWithdrawalDAO
	if err := row.Scan(
		&resultDAO.ID,
		&resultDAO.UserID,
		&resultDAO.OrderNumber,
		&resultDAO.Sum,
		&resultDAO.ProcessedAt,
	); err != nil {
		u.logger.Error("postgres[user]", "CreateUserWithdrawal", "Failed to scan withdrawal result", err)
		return nil, pghelpers.WrapError(err)
	}

	updateQuery, updateArgs, err := quries.BuildUpdateUserBalanceQuery(
		pghelpers.ToNullUUID(input.UserID),
		pghelpers.ToNullInt64(input.CurrentBalance),
		pghelpers.ToNullInt64(input.WitdrawnBalance),
	)
	if err != nil {
		u.logger.Error("postgres[user]", "CreateUserWithdrawal", "Failed to build update balance query", err)
		return nil, pghelpers.WrapError(err)
	}

	_, err = tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		u.logger.Error("postgres[user]", "CreateUserWithdrawal", "Failed to update user balance", err)
		return nil, pghelpers.WrapError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		u.logger.Error("postgres[user]", "CreateUserWithdrawal", "Failed to commit transaction", err)
		return nil, pghelpers.WrapError(err)
	}

	return resultDAO.ToModel(), nil
}

func (u *UserRepo) GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.UserWithdrawal[int64], error) {
	query, args, err := quries.BuildGetUserWithdrawalsQuery(pghelpers.ToNullUUID(&userID))
	if err != nil {
		u.logger.Error("postgres[user]", "GetUserWithdrawals", "Failed to build query", err)
		return nil, pghelpers.WrapError(err)
	}

	rows, err := u.c.Query(ctx, query, args...)
	if err != nil {
		u.logger.Error("postgres[user]", "GetUserWithdrawals", "Failed to execute query", err)
		return nil, pghelpers.WrapError(err)
	}
	defer rows.Close()

	var withdrawals []model.UserWithdrawal[int64]
	for rows.Next() {
		var w dao.UserWithdrawalDAO
		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.OrderNumber,
			&w.Sum,
			&w.ProcessedAt,
		); err != nil {
			u.logger.Error("postgres[user]", "GetUserWithdrawals", "Failed to scan row", err)
			return nil, pghelpers.WrapError(err)
		}
		withdrawals = append(withdrawals, *w.ToModel())
	}

	if err := rows.Err(); err != nil {
		u.logger.Error("postgres[user]", "GetUserWithdrawals", "Rows error", err)
		return nil, pghelpers.WrapError(err)
	}

	return withdrawals, nil
}
