package repository

import (
	"context"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

type Notes interface {
	Create(ctx context.Context, note entity.Note) (int, error)
	GetByTitle(ctx context.Context, title string) (entity.Note, error)
	GetById(ctx context.Context, id int) (entity.Note, error)
	GetNotes(ctx context.Context) ([]entity.Note, error)
	GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time) ([]entity.Note, error)
	UpdateNote(ctx context.Context, id int, title, description, status string) error
	DeleteById(ctx context.Context, id int) error
	DeleteNotes(ctx context.Context) error
}

type Users interface {
	SignIn()
	SignUp()
	Logout()
}

type Repository struct {
	Notes
	Users
}

func New(repo Repository) *Repository {
	return &Repository{
		Notes: repo.Notes,
		Users: repo.Users,
	}
}
