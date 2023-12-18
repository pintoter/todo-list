package service

import (
	"context"
	"github.com/pintoter/todo-list/internal/entity"
)

func (s *Service) SignUp(ctx context.Context, user entity.User) (int, error) {
	if s.isLoginExists(ctx, user.Login) || s.isEmailExists(ctx, user.Email) {
		return 0, entity.ErrUserExists
	}

	id, err := s.repo.SignUp(ctx, user)
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

	return false
}

func (s *Service) isEmailExists(ctx context.Context, email string) bool {

	return false
}
