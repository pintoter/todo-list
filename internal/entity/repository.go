package entity

import (
	"context"
	"time"
)

type Notes interface {
	Create(ctx context.Context, note Note) (int, error)
	GetByTitle(ctx context.Context, title string) (Note, error)
	GetById(ctx context.Context, id int) (Note, error)
	GetNotes(ctx context.Context) ([]Note, error)
	GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time) ([]Note, error)
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
