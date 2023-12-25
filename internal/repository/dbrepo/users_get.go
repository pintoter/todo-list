package dbrepo

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

func (r *DBRepo) GetByID(ctx context.Context, id int) (entity.User, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.User{}, err
	}
	defer func() { _ = tx.Rollback() }()

	builder := sq.Select("id", "email", "login", "password", "register_at").
		From(users).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	var user entity.User
	err = tx.QueryRowContext(ctx, query, args...).
		Scan(&user.ID, &user.Email, &user.Login, &user.Password, &user.RegisteredAt)
	if err != nil {
		return entity.User{}, err
	}

	return user, tx.Commit()
}

func (r *DBRepo) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.User{}, err
	}
	defer func() { _ = tx.Rollback() }()

	builder := sq.Select("id", "email", "login", "password", "register_at").
		From(users).
		Where(sq.Eq{"login": login}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	var user entity.User
	err = tx.QueryRowContext(ctx, query, args...).
		Scan(&user.ID, &user.Email, &user.Login, &user.Password, &user.RegisteredAt)
	if err != nil {
		return entity.User{}, err
	}

	return user, tx.Commit()
}

func (r *DBRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.User{}, err
	}
	defer func() { _ = tx.Rollback() }()

	builder := sq.Select("id", "email", "login", "password", "register_at").
		From(users).
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	var user entity.User
	err = tx.QueryRowContext(ctx, query, args...).
		Scan(&user.ID, &user.Email, &user.Login, &user.Password, &user.RegisteredAt)
	if err != nil {
		return entity.User{}, err
	}

	return user, tx.Commit()
}
