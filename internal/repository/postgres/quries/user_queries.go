package quries

import (
	"database/sql"

	"github.com/FlyKarlik/gofemart/internal/repository/postgres/dao"
	"github.com/google/uuid"

	"github.com/Masterminds/squirrel"
)

func BuildCreateUserQuery(user dao.UserInputDAO) (string, []interface{}, error) {
	query, args, err := squirrel.
		Insert(`"user"`).
		Columns("login", "password_hash").
		Values(user.Login, user.Password).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

func BuildGetUserByLoginQuery(login sql.NullString) (string, []interface{}, error) {
	query, args, err := squirrel.
		Select("*").
		From(`"user"`).
		Where(squirrel.Eq{"login": login}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

func BuildGetUserByIDQuery(id uuid.NullUUID) (string, []interface{}, error) {
	query, args, err := squirrel.
		Select("*").
		From(`"user"`).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

func BuildCreateOrderQuery(order dao.UserOrderInputDAO) (string, []interface{}, error) {
	query, args, err := squirrel.
		Insert(`"user_order"`).
		Columns("number", "user_id", "status").
		Values(order.Number, order.UserID, order.Status).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	return query, args, err
}

func BuildCheckOrderExistsQuery(number sql.NullString, userID uuid.NullUUID) (string, []interface{}, error) {
	query, args, err := squirrel.
		Select("1").
		From(`"user_order"`).
		Where(squirrel.And{
			squirrel.Eq{"number": number},
			squirrel.Eq{"user_id": userID},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	return query, args, err
}

func BuildGetUserOrdersQuery(userID uuid.NullUUID) (string, []interface{}, error) {
	query := squirrel.
		Select("*").
		From(`"user_order"`).
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("uploaded_at ASC").
		PlaceholderFormat(squirrel.Dollar)

	return query.ToSql()
}

func BuildCreateUserBalanceQuery(userID uuid.NullUUID) (string, []interface{}, error) {
	return squirrel.
		Insert("user_balance").
		Columns("user_id", "current", "withdrawn").
		Values(userID, 0, 0).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
}

func BuildUpdateUserBalanceQuery(userID uuid.NullUUID, current sql.NullInt64, withdrawn sql.NullInt64) (string, []interface{}, error) {
	queryBuilder := squirrel.
		Update("user_balance").
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar)

	if current.Valid {
		queryBuilder = queryBuilder.Set("current", current.Int64)
	}

	if withdrawn.Valid {
		queryBuilder = queryBuilder.Set("withdrawn", withdrawn.Int64)
	}

	return queryBuilder.ToSql()
}

func BuildGetUserBalanceQuery(userID uuid.NullUUID) (string, []interface{}, error) {
	query := squirrel.
		Select("*").
		From(`"user_balance"`).
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar)

	return query.ToSql()
}

func BuildInsertWithdrawalQuery(withdrawal dao.UserWithdrawalInputDAO) (string, []interface{}, error) {
	query := squirrel.
		Insert(`"user_withdrawal"`).
		Columns("user_id", "order_number", "sum").
		Values(withdrawal.UserID, withdrawal.OrderNumber, withdrawal.Sum).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar)
	return query.ToSql()
}

func BuildGetUserWithdrawalsQuery(userID uuid.NullUUID) (string, []interface{}, error) {
	query := squirrel.
		Select("*").
		From(`"user_withdrawal"`).
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("processed_at ASC").
		PlaceholderFormat(squirrel.Dollar)
	return query.ToSql()
}
