package dbrepo

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

func createNoteBuilder(note entity.Note) (string, []interface{}, error) {
	builder := sq.Insert(notes).
		Columns("user_id", "title", "description", "date", "status").
		Values(note.UserId, note.Title, note.Description, note.Date, note.Status).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) CreateNote(ctx context.Context, note entity.Note) (int, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := createNoteBuilder(note)
	if err != nil {
		return 0, err
	}

	var noteId int
	err = tx.QueryRowContext(ctx, query, args...).Scan(&noteId)
	if err != nil {
		return 0, err
	}

	return noteId, tx.Commit()
}
