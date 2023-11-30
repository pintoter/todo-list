package repository

import (
	"context"
	"database/sql"
	"log"
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

/*
db, err := sql.Open("mysql", *dsn)
if err != nil {
    log.Fatal(err)
}

db.SetMaxOpenConns(100)
db.SetMaxIdleConns(5)

*/

func (n *NotesRepo) Create(ctx context.Context, note entity.Note) (int, error) {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	builderInsert := sq.Insert(notes).
		Columns("title", "description", "date", "status").
		Values(note.Title, note.Description, note.Date, note.Status).Suffix("RETURNING id").PlaceholderFormat(sq.Dollar)

	query, args, err := builderInsert.ToSql()
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
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	builderSelect := sq.Select("*").From(notes).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return entity.Note{}, err
	}

	var note entity.Note

	err = tx.QueryRowContext(ctx, query, args...).Scan(&note.ID, &note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, tx.Commit()
}

func (n *NotesRepo) GetByTitle(ctx context.Context, title string) (entity.Note, error) {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	builderSelect := sq.Select("*").From(notes).Where(sq.Eq{"title": title}).PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return entity.Note{}, err
	}

	var note entity.Note

	err = tx.QueryRowContext(ctx, query, args...).Scan(&note.ID, &note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, tx.Commit()
}

func (n *NotesRepo) GetNotes(ctx context.Context, limit, offset int) ([]entity.Note, error) {
	log.Println("YA ZDES")

	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	builderSelect := sq.Select("*").
		From(notes).
		OrderBy("id ASC").
		PlaceholderFormat(sq.Dollar)

	if limit != 0 {
		builderSelect = builderSelect.Limit(uint64(limit)).Offset(uint64(offset))
	}

	query, args, err := builderSelect.ToSql()
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

func (n *NotesRepo) GetNotesByStatus(ctx context.Context, limit, offset int, status string) ([]entity.Note, error) {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	builderSelect := sq.Select("*").
		From(notes).
		Where(sq.Eq{"status": status}).
		PlaceholderFormat(sq.Dollar)

	if limit != 0 {
		builderSelect = builderSelect.Limit(uint64(limit)).Offset(uint64(offset))
	}

	query, args, err := builderSelect.ToSql()
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

func (n *NotesRepo) GetNotesByStatusAndDate(ctx context.Context, limit, offset int, status string, date time.Time) ([]entity.Note, error) {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	builderSelect := sq.Select("*").
		From(notes).
		Where(sq.Eq{"status": status}).
		Where(sq.Eq{"date": date}).
		OrderBy("id ASC").
		PlaceholderFormat(sq.Dollar)

	if limit != 0 {
		builderSelect = builderSelect.Limit(uint64(limit)).Offset(uint64(offset))
	}

	query, args, err := builderSelect.ToSql()
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

func (n *NotesRepo) UpdateNote(ctx context.Context, id int, title, description, status string) error {
	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	builderUpdate := sq.Update(notes).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	if title != "" {
		builderUpdate = builderUpdate.Set("title", title)
	}

	if description != "" {
		builderUpdate = builderUpdate.Set("description", description)
	}

	if status != "" {
		builderUpdate = builderUpdate.Set("status", status)
	}

	query, args, err := builderUpdate.ToSql()
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

	builderUpdate := sq.Delete(notes).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	query, args, err := builderUpdate.ToSql()
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

	builderUpdate := sq.Delete(notes).PlaceholderFormat(sq.Dollar)

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
