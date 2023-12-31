package dbrepo

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

/*-----------------------------
					GET NOTE
 ----------------------------- */

func getNoteBuilder(data any, userId int) (string, []interface{}, error) {
	builder := sq.Select("id", "user_id", "title", "description", "date", "status").
		From(notes).
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar)

	switch data.(type) {
	case int:
		builder = builder.Where(sq.Eq{"id": data})
	case string:
		builder = builder.Where(sq.Eq{"title": data})
	}

	return builder.ToSql()
}

func (r *DBRepo) GetNoteById(ctx context.Context, id, userId int) (entity.Note, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getNoteBuilder(id, userId)
	if err != nil {
		return entity.Note{}, err
	}

	var note entity.Note
	err = tx.QueryRowContext(ctx, query, args...).Scan(&note.ID, &note.UserId, &note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, tx.Commit()
}

func (r *DBRepo) GetNoteByTitle(ctx context.Context, title string, userId int) (entity.Note, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getNoteBuilder(title, userId)
	if err != nil {
		return entity.Note{}, err
	}

	var note entity.Note
	err = tx.QueryRowContext(ctx, query, args...).Scan(&note.ID, &note.UserId, &note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, tx.Commit()
}

/*-----------------------------
					GET ALL NOTES
 ----------------------------- */

func getNotesBuilder(limit, offset int, status string, date time.Time, userId int) (string, []interface{}, error) {
	builder := sq.Select("id", "user_id", "title", "description", "date", "status").
		From(notes).
		OrderBy("id ASC").
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar)

	if status != "" {
		builder = builder.Where(sq.Eq{"status": status})
	}

	if !date.Equal(time.Time{}) {
		builder = builder.Where(sq.Eq{"date": date})
	}

	if limit != 0 || offset != 0 {
		builder = builder.Limit(uint64(limit)).Offset(uint64(offset))
	}

	return builder.ToSql()
}

func (r *DBRepo) GetNotes(ctx context.Context, userId int) ([]entity.Note, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getNotesBuilder(0, 0, "", time.Time{}, userId)
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
		if err := rows.Scan(&note.ID, &note.UserId, &note.Title, &note.Description, &note.Date, &note.Status); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, tx.Commit()
}

func (r *DBRepo) GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time, userId int) ([]entity.Note, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getNotesBuilder(limit, offset, status, date, userId)
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
		if err := rows.Scan(&note.ID, &note.UserId, &note.Title, &note.Description, &note.Date, &note.Status); err != nil {
			return notes, err
		}
		notes = append(notes, note)
	}

	return notes, tx.Commit()
}
