package dbrepo

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

func getCreateUserQuery(user entity.User) (string, []interface{}, error) {
	builder := sq.Insert(users).
		Columns("email", "login", "password", "register_at").
		Values(user.Email, user.Login, user.Password, time.Now()).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query, args, err := getCreateUserQuery(user)
	if err != nil {
		return 0, err
	}

	var userId int
	err = tx.QueryRowContext(ctx, query, args...).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, tx.Commit()
}
