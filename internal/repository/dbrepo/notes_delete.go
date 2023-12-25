package dbrepo

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func getDeleteByIdQuery(id int) (string, []interface{}, error) {
	builder := sq.Delete(notes).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) DeleteNoteById(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getDeleteByIdQuery(id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func getDeleteNotesQuery() (string, []interface{}, error) {
	builder := sq.Delete(notes).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) DeleteNotes(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getDeleteNotesQuery()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
