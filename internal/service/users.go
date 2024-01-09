package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

func (s *Service) SignUp(ctx context.Context, email, login, password string) (int, error) {
	if s.isLoginExists(ctx, login) || s.isEmailExists(ctx, email) {
		return 0, entity.ErrUserExists
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return 0, err
	}

	user := entity.User{
		Email:    email,
		Login:    login,
		Password: hashedPassword,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *Service) SignIn(ctx context.Context, login, password string) (Tokens, error) {
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetUserByCredentials(ctx, login, hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Tokens{}, entity.ErrUserNotExist
		}

		return Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetUserByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *Service) createSession(ctx context.Context, id int) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(id, s.accessTokenTTL)
	if err != nil {
		return Tokens{}, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	token := entity.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return res, s.repo.SetSession(ctx, id, token)
}

func (s *Service) isLoginExists(ctx context.Context, login string) bool {
	_, err := s.repo.GetUserByLogin(ctx, login)

	return err == nil
}

func (s *Service) isEmailExists(ctx context.Context, email string) bool {
	_, err := s.repo.GetUserByEmail(ctx, email)

	return err == nil
}
