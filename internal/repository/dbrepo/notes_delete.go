package dbrepo

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func getDeleteByIdQuery(id, userId int) (string, []interface{}, error) {
	builder := sq.Delete(notes).
		Where(sq.Eq{"id": id, "user_id": userId}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) DeleteNoteById(ctx context.Context, id, userId int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getDeleteByIdQuery(id, userId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func getDeleteNotesQuery(userId int) (string, []interface{}, error) {
	builder := sq.Delete(notes).
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) DeleteNotes(ctx context.Context, userId int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getDeleteNotesQuery(userId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
