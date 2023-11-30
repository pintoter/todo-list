package repository

import (
	"context"
	"database/sql"

	"github.com/pintoter/todo-list/internal/entity"
)

type Notes interface {
	Create(ctx context.Context, note entity.Note) (int, error)
	GetByTitle(ctx context.Context, title string) (entity.Note, error)
	// GetById(ctx context.Context, id int) (entity.Note, error)
	// GetAll(ctx context.Context, limit, offset int) ([]entity.Note, error)
	// GetByStatus(ctx context.Context, status string, limit, offset int) ([]entity.Note, error)
	// GetByDate(ctx context.Context, date time.Time) ([]entity.Note, error)
	// UpdateInfo(ctx context.Context, title, description string, time time.Time) error
	// UpdateStatus(ctx context.Context, status string, time time.Time) error
	// DeleteByTitle(ctx context.Context, title string) error
}

type Repository struct {
	Notes
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Notes: NewNotes(db),
	}
}
