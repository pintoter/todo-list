package dbrepo

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

func getGetByIdQuery(id int) (string, []interface{}, error) {
	builder := sq.Select("title", "description", "date", "status").
		From(notes).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) GetNoteById(ctx context.Context, id int) (entity.Note, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getGetByIdQuery(id)
	if err != nil {
		return entity.Note{}, err
	}

	var note entity.Note
	err = tx.QueryRowContext(ctx, query, args...).Scan(&note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, tx.Commit()
}

func getGetByTitleQuery(title string) (string, []interface{}, error) {
	builder := sq.Select("title", "description", "date", "status").
		From(notes).
		Where(sq.Eq{"title": title}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) GetNoteByTitle(ctx context.Context, title string) (entity.Note, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getGetByTitleQuery(title)
	if err != nil {
		return entity.Note{}, err
	}

	var note entity.Note
	err = tx.QueryRowContext(ctx, query, args...).Scan(&note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, tx.Commit()
}

func getGetNotesQuery() (string, []interface{}, error) {
	builder := sq.Select("id", "title", "description", "date", "status").
		From(notes).
		OrderBy("id ASC").
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) GetNotes(ctx context.Context) ([]entity.Note, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getGetNotesQuery()
	if err != nil {
		return nil, err
	}

	var notes []entity.Note
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note entity.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Description, &note.Date, &note.Status); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, tx.Commit()
}

func getNotesExtendedQuery(limit, offset int, status string, date time.Time) (string, []interface{}, error) {
	builder := sq.Select("id", "title", "description", "date", "status").
		From(notes).
		OrderBy("id ASC").
		PlaceholderFormat(sq.Dollar)

	if status != "" {
		builder = builder.Where(sq.Eq{"status": status})
	}

	if !date.Equal(time.Time{}) {
		builder = builder.Where(sq.Eq{"date": date})
	}

	builder = builder.Limit(uint64(limit)).Offset(uint64(offset))

	return builder.ToSql()
}

func (r *DBRepo) GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time) ([]entity.Note, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getNotesExtendedQuery(limit, offset, status, date)
	if err != nil {
		return nil, err
	}

	var notes []entity.Note
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note entity.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Description, &note.Date, &note.Status); err != nil {
			return notes, err
		}
		notes = append(notes, note)
	}

	return notes, tx.Commit()
}
