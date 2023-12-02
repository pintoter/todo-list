package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/todo-list/internal/entity"
)

const (
	notes = "notes"
)

type NotesRepo struct {
	db *sql.DB
}

func NewNotes(db *sql.DB) *NotesRepo {
	return &NotesRepo{
		db: db,
	}
}

func (n *NotesRepo) Create(ctx context.Context, note entity.Note) (int, error) {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	builder := sq.Insert(notes).
		Columns("title", "description", "date", "status").
		Values(note.Title, note.Description, note.Date, note.Status).Suffix("RETURNING id").PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
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

func (n *NotesRepo) GetById(ctx context.Context, id int) (entity.Note, error) {
	builderSelect := sq.Select("title", "description", "date", "status").From(notes).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return entity.Note{}, err
	}

	var note entity.Note

	err = n.db.QueryRowContext(ctx, query, args...).Scan(&note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, nil
}

func (n *NotesRepo) GetByTitle(ctx context.Context, title string) (entity.Note, error) {
	builderSelect := sq.Select("title", "description", "date", "status").From(notes).Where(sq.Eq{"title": title}).PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return entity.Note{}, err
	}

	var note entity.Note
	err = n.db.QueryRowContext(ctx, query, args...).Scan(&note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, err
}

func (n *NotesRepo) GetNotes(ctx context.Context) ([]entity.Note, error) {
	builder := sq.Select("id", "title", "description", "date", "status").
		From(notes).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var notes []entity.Note
	rows, err := n.db.QueryContext(ctx, query, args...)
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

	return notes, err
}

func getSelectQuery(limit, offset int, status string, date time.Time) (string, []interface{}, error) {
	builder := sq.Select("id", "title", "description", "date", "status").
		From(notes).
		PlaceholderFormat(sq.Dollar)

	if status != "" {
		builder = builder.Where(sq.Eq{"status": status})
	}

	if !date.Equal(time.Time{}) {
		builder = builder.Where(sq.Eq{"date": date})
	}

	if limit != 0 {
		builder = builder.Limit(uint64(limit)).Offset(uint64(offset))
	}

	return builder.ToSql()
}

func (n *NotesRepo) GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time) ([]entity.Note, error) {
	query, args, err := getSelectQuery(limit, offset, status, date)
	if err != nil {
		return nil, err
	}

	var notes []entity.Note
	rows, err := n.db.QueryContext(ctx, query, args...)
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

	return notes, err
}

func getUpdateBuilder(id int, title, description, status string) (string, []interface{}, error) {
	builder := sq.Update(notes).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

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

func (n *NotesRepo) UpdateNote(ctx context.Context, id int, title, description, status string) error {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getUpdateBuilder(id, title, description, status)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (n *NotesRepo) DeleteById(ctx context.Context, id int) error {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	builder := sq.Delete(notes).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (n *NotesRepo) DeleteNotes(ctx context.Context) error {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	builder := sq.Delete(notes).PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
