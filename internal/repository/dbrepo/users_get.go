package dbrepo

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

type getInput struct {
	id    *int
	login *string
	email *string
}

func getUserBuilder(data getInput) (string, []interface{}, error) {
	builder := sq.Select("id", "email", "login", "password", "register_at").
		From(users).
		PlaceholderFormat(sq.Dollar)

	if data.id != nil {
		builder = builder.Where(sq.Eq{"id": *(data.id)})
	}

	if data.login != nil {
		builder = builder.Where(sq.Eq{"login": *(data.login)})
	}

	if data.email != nil {
		builder = builder.Where(sq.Eq{"email": *(data.email)})
	}

	return builder.ToSql()
}

func (r *DBRepo) GetByID(ctx context.Context, id int) (entity.User, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.User{}, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getUserBuilder(getInput{id: &id})
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

	query, args, err := getUserBuilder(getInput{login: &login})
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

	query, args, err := getUserBuilder(getInput{email: &email})
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
