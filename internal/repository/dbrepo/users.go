package dbrepo

import (
	"context"

	"github.com/pintoter/todo-list/internal/entity"
)

func (r *DBRepo) SignIn(ctx context.Context, login, password string) (entity.User, error) {

	return entity.User{}, nil
}

func (r *DBRepo) RefreshToken(ctx context.Context, refreshToken string) (entity.User, error) {

	return entity.User{}, nil
}
