package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pintoter/todo-list/internal/entity"
)

func (s *Service) CreateNote(ctx context.Context, note entity.Note) (int, error) {
	if s.isNoteExists(ctx, note.Title) {
		return 0, entity.ErrNoteExists
	}

	return s.repo.Create(ctx, note)
}

func (s *Service) GetById(ctx context.Context, id int) (entity.Note, error) {
	note, err := s.repo.GetById(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Note{}, entity.ErrNoteNotExits
	}

	return note, nil
}

func (s *Service) GetNotes(ctx context.Context, limit, offset int, status, date string) ([]entity.Note, error) {
	if status != "" && status != entity.StatusDone && status != entity.StatusNotDone {
		return nil, entity.ErrInvalidStatus
	}

	notes, err := s.repo.GetNotes(ctx, limit, offset, status, date)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *Service) UpdateNote(ctx context.Context, id int, title, description, status string) error {
	if s.isNoteExists(ctx, id) {
		return entity.ErrNoteNotExits
	}

	return s.repo.UpdateNote(ctx, id, title, description, status)
}

func (s *Service) DeleteById(ctx context.Context, id int) error {
	if s.isNoteExists(ctx, id) {
		return entity.ErrNoteNotExits
	}

	return s.repo.DeleteById(ctx, id)
}

func (s *Service) DeleteNotes(ctx context.Context) error {
	return s.repo.DeleteNotes(ctx)
}

func (s *Service) isNoteExists(ctx context.Context, data interface{}) bool {
	var err error
	switch data.(type) {
	case int:
		_, err = s.repo.GetById(ctx, data.(int))
	case string:
		_, err = s.repo.GetByTitle(ctx, data.(string))
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}
