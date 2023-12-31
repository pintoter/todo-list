package dbrepo

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func updateBuilder(id, userId int, title, description, status string) (string, []interface{}, error) {
	builder := sq.Update(notes).
		Where(sq.Eq{"id": id, "user_id": userId}).
		PlaceholderFormat(sq.Dollar)

	if title != "" {
		builder = builder.Set("title", title)
	}

	if description != "" {
		builder = builder.Set("description", description)
	}

	if status != "" {
		builder = builder.Set("status", status)
	}

	return builder.ToSql()
}

func (r *DBRepo) UpdateNote(ctx context.Context, id, userId int, title, description, status string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := updateBuilder(id, userId, title, description, status)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
