package service

import (
	"context"
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

	id, err := s.repo.CreateUser(ctx, entity.User{
		Email:    email,
		Login:    login,
		Password: hashedPassword,
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) SignIn(ctx context.Context, login, password string) (entity.Token, error) {

	return entity.Token{}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (entity.Token, error) {

	return entity.Token{}, nil
}

func (s *Service) isLoginExists(ctx context.Context, login string) bool {
	_, err := s.repo.GetByLogin(ctx, login)

	return err == nil
}

func (s *Service) isEmailExists(ctx context.Context, email string) bool {
	_, err := s.repo.GetByEmail(ctx, email)

	return err == nil
}
