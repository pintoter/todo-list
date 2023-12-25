package repository

import (
	"context"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

type NotesRepository interface {
	CreateNote(ctx context.Context, note entity.Note) (int, error)
	GetNoteByTitle(ctx context.Context, title string) (entity.Note, error)
	GetNoteById(ctx context.Context, id int) (entity.Note, error)
	GetNotes(ctx context.Context) ([]entity.Note, error)
	GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time) ([]entity.Note, error)
	UpdateNote(ctx context.Context, id int, title, description, status string) error
	DeleteNoteById(ctx context.Context, id int) error
	DeleteNotes(ctx context.Context) error
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByLogin(ctx context.Context, login string) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	SignIn(ctx context.Context, login, password string) (entity.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (entity.Token, error)
}

type Repository interface {
	NotesRepository
	UsersRepository
}
