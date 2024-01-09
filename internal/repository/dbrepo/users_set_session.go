package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

func setUserSessionBuilder(userId int, refreshToken string, expiresAt time.Time) (string, []interface{}, error) {
	builder := sq.Update(users).
		Set("refresh_token", refreshToken).
		Set("expires_at", expiresAt).
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) SetSession(ctx context.Context, userId int, session entity.Session) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := setUserSessionBuilder(userId, session.RefreshToken, session.ExpiresAt)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return tx.Commit()
}
