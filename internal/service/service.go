package service

import (
	"context"
	"database/sql"

	"github.com/pintoter/todo-list/internal/entity"
	"github.com/pintoter/todo-list/internal/repository"
)

type INotesRepository interface {
	Create(ctx context.Context, note entity.Note) (int, error)
	// GetById(ctx context.Context, id int) (entity.Note, error)
	GetByTitle(ctx context.Context, title string) (entity.Note, error)
	// GetAll(ctx context.Context, limit, offset int) ([]entity.Note, error)
	// GetByStatus(ctx context.Context, status string, limit, offset int) ([]entity.Note, error)
	// GetByDate(ctx context.Context, date time.Time) ([]entity.Note, error)
	// UpdateStatus(ctx context.Context, status string, time time.Time) error
	// DeleteByTitle(ctx context.Context, title string) error
	// UpdateInfo(ctx context.Context, title, description string, time time.Time) error
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
