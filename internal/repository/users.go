package repository

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

func (r *DBRepo) SignUp(ctx context.Context, user entity.User) (int, error) {
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

	builder := sq.Insert(users).Columns("email", "login", "password", "register_at")

	return 0, nil
}

func (r *DBRepo) SignIn(ctx context.Context, login, password string) (entity.Token, error) {

	return entity.Token{}, nil
}

func (r *DBRepo) RefreshToken(ctx context.Context, refreshToken string) (entity.Token, error) {

	return entity.Token{}, nil
}
