package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
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

	return 0, nil
}

func (n *NotesRepo) GetByTitle(ctx context.Context, title string) (entity.Note, error) {

	return entity.Note{}, nil
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
