package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
	"github.com/pintoter/todo-list/internal/repository"
)

type INotesRepository interface {
	Create(ctx context.Context, note entity.Note) (int, error)
	GetByTitle(ctx context.Context, title string) (entity.Note, error)
	GetById(ctx context.Context, id int) (entity.Note, error)
	GetNotes(ctx context.Context) ([]entity.Note, error)
	GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time) ([]entity.Note, error)
	UpdateNote(ctx context.Context, id int, title, description, status string) error
	DeleteById(ctx context.Context, id int) error
	DeleteNotes(ctx context.Context) error
}

type IRepository interface {
	INotesRepository
}

type Service struct {
	repo IRepository
}

func New(db *sql.DB) *Service {
	return &Service{
		repo: repository.New(db),
	}
}
