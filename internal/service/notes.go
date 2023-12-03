package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

func (s *Service) CreateNote(ctx context.Context, note entity.Note) error {
	if s.isNoteExists(ctx, note.Title) {
		return entity.ErrNoteExists
	}

	_, err := s.repo.Create(ctx, note)
	return err
}

func (s *Service) GetById(ctx context.Context, id int) (entity.Note, error) {
	note, err := s.repo.GetById(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Note{}, entity.ErrNoteNotExists
	}

	return note, nil
}

func (s *Service) GetNotes(ctx context.Context) ([]entity.Note, error) {
	notes, err := s.repo.GetNotes(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotesNotExist
		} else {
			return nil, err
		}
	}

	return notes, nil
}

func (s *Service) GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time) ([]entity.Note, error) {
	notes, err := s.repo.GetNotesExtended(ctx, limit, offset, status, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotesNotExist
		} else {
			return nil, err
		}
	}

	return notes, nil
}

func (s *Service) UpdateNote(ctx context.Context, id int, title, description, status string) error {
	if s.isNoteExists(ctx, id) {
		return entity.ErrNoteNotExists
	}

	return s.repo.UpdateNote(ctx, id, title, description, status)
}

func (s *Service) DeleteById(ctx context.Context, id int) error {
	if s.isNoteExists(ctx, id) {
		return entity.ErrNoteNotExists
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
