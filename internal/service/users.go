package service

import (
	"context"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type TokenManager interface {
	NewJWT(userId int, ttl time.Duration) (string, error)
	NewRefreshToken() (string, error)
	ParseToken(accessToken string) (int, error)
}

type Config interface {
	GetAccessTokenTTL() time.Duration
	GetRefreshTokenTTL() time.Duration
}

type UserService struct {
	repo            IUsersRepository
	hasher          PasswordHasher
	tokenManager    TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(cfg Config, repo IUsersRepository, hasher PasswordHasher, tokenManager TokenManager) *UserService {
	return &UserService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  cfg.GetAccessTokenTTL(),
		refreshTokenTTL: cfg.GetRefreshTokenTTL(),
	}
}

func (u *UserService) SignUp(ctx context.Context, user entity.User) (int, error) {
	if u.isUserExists(ctx, user.)

	return 0, nil
}

func (u *UserService) SignIn(ctx context.Context, login, password string) (entity.Token, error) {

	return entity.Token{}, nil
}

func (u *UserService) RefreshToken(ctx context.Context, refreshToken string) (entity.Token, error) {

	return entity.Token{}, nil
}
