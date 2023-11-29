package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

func (s *Service) Create(ctx context.Context, note entity.Note) (int, error) {

	return 0, nil
}

func (s *Service) GetAll(ctx context.Context, limit, offset int) ([]entity.Note, error) {

	return nil, nil
}

func (s *Service) GetByStatus(ctx context.Context, status string, limit, offset int) ([]entity.Note, error) {

	return nil, nil
}

func (s *Service) GetByDate(ctx context.Context, date time.Time) ([]entity.Note, error) {

	return nil, nil
}

func (s *Service) UpdateInfo(ctx context.Context, title, description string, time time.Time) error {

	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, status string, time time.Time) error {

	return nil
}

func (s *Service) DeleteByTitle(ctx context.Context, title string) error {

	return nil
}

func (s *Service) isNoteExists(ctx context.Context, note entity.Note) bool {
	_, err := s.repo.GetByTitle(ctx, note.Title)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}
