package service

import (
	"time"

	"github.com/pintoter/todo-list/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

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

func New(cfg Config, repo repository.Repository, hasher PasswordHasher, tokenManager TokenManager) *Service {
	return &Service{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  cfg.GetAccessTokenTTL(),
		refreshTokenTTL: cfg.GetAccessTokenTTL(),
	}
}
