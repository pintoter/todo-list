package dbrepo

import (
	"context"

	"github.com/pintoter/todo-list/internal/entity"
)

func (r *DBRepo) SignIn(ctx context.Context, login, password string) (entity.Token, error) {

	return entity.Token{}, nil
}

func (r *DBRepo) RefreshToken(ctx context.Context, refreshToken string) (entity.Token, error) {

	return entity.Token{}, nil
}
