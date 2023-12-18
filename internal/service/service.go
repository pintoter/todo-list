package service

import (
	"context"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Config interface {
	GetAccessTokenTTL() time.Duration
	GetRefreshTokenTTL() time.Duration
}

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

type IUsersRepository interface {
	SignUp(ctx context.Context, user entity.User) (int, error)
	SignIn(ctx context.Context, login, password string) (entity.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (entity.Token, error)
}

type IRepository interface {
	INotesRepository
	IUsersRepository
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type TokenManager interface {
	NewJWT(userId int, ttl time.Duration) (string, error)
	NewRefreshToken() (string, error)
	ParseToken(accessToken string) (int, error)
}

type Service struct {
	repo            IRepository
	hasher          PasswordHasher
	tokenManager    TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func New(cfg Config, repo IRepository, hasher PasswordHasher, tokenManager TokenManager) *Service {
	return &Service{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  cfg.GetAccessTokenTTL(),
		refreshTokenTTL: cfg.GetAccessTokenTTL(),
	}
}
