package repository

import (
	"context"
	"database/sql"

	"github.com/pintoter/todo-list/internal/entity"
)

type Notes interface {
	Create(ctx context.Context, note entity.Note) (int, error)
	GetByTitle(ctx context.Context, title string) (entity.Note, error)
	GetById(ctx context.Context, id int) (entity.Note, error)
	GetNotes(ctx context.Context, limit, offset int, status string, date string) ([]entity.Note, error)
	UpdateNote(ctx context.Context, id int, title, description, status string) error
	DeleteById(ctx context.Context, id int) error
	DeleteNotes(ctx context.Context) error
}

type Repository struct {
	Notes
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Notes: NewNotes(db),
	}
}
