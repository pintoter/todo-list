package service

import (
	"time"

	"github.com/pintoter/todo-list/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Config interface {
	GetAccessTokenTTL() time.Duration
	GetRefreshTokenTTL() time.Duration
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
	repo            repository.Repository
	hasher          PasswordHasher
	tokenManager    TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type Deps struct {
	Cfg          Config
	Repo         repository.Repository
	Hasher       PasswordHasher
	TokenManager TokenManager
}

func New(deps Deps) *Service {
	return &Service{
		repo:            deps.Repo,
		hasher:          deps.Hasher,
		tokenManager:    deps.TokenManager,
		accessTokenTTL:  deps.Cfg.GetAccessTokenTTL(),
		refreshTokenTTL: deps.Cfg.GetRefreshTokenTTL(),
	}
}
