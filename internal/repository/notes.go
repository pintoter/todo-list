package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	var (
		id    int
		query = fmt.Sprintf("INSERT INTO %s (title, description, date, status) VALUES ($1, $2, $3, $4) RETURNING id", notes)
	)

	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRowContext(ctx, query, note.Title, note.Description, note.Date, note.Status).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (n *NotesRepo) GetById(ctx context.Context, id int) (entity.Note, error) {

	return entity.Note{}, nil
}

func (n *NotesRepo) GetByTitle(ctx context.Context, title string) (entity.Note, error) {
	var (
		note  entity.Note
		query = fmt.Sprintf("SELECT title, description, date, status FROM %s", notes)
	)

	tx, err := n.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Note{}, err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRowContext(ctx, query).Scan(&note.Title, &note.Description, &note.Date, &note.Status)
	if err != nil {
		return entity.Note{}, err
	}

	return note, nil
}

func (n *NotesRepo) GetAll(ctx context.Context, limit, offset int) ([]entity.Note, error) {

	return nil, nil
}

func (n *NotesRepo) GetByStatus(ctx context.Context, status string, limit, offset int) ([]entity.Note, error) {

	return nil, nil
}

func (n *NotesRepo) GetByDate(ctx context.Context, date time.Time) ([]entity.Note, error) {

	return nil, nil
}

func (n *NotesRepo) UpdateInfo(ctx context.Context, title, description string, time time.Time) error {

	return nil
}

func (n *NotesRepo) UpdateStatus(ctx context.Context, status string, time time.Time) error {

	return nil
}

func (n *NotesRepo) DeleteByTitle(ctx context.Context, title string) error {

	return nil
}
