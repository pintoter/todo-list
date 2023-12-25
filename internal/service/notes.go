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

	_, err := s.repo.CreateNote(ctx, note)
	return err
}

func (s *Service) GetById(ctx context.Context, id int) (entity.Note, error) {
	note, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Note{}, entity.ErrNoteNotExists
		} else {
			return entity.Note{}, err
		}
	}

	return note, nil
}

func (s *Service) GetNotes(ctx context.Context) ([]entity.Note, error) {
	notes, err := s.repo.GetNotes(ctx)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *Service) GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time) ([]entity.Note, error) {
	notes, err := s.repo.GetNotesExtended(ctx, limit, offset, status, date)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *Service) UpdateNote(ctx context.Context, id int, title, description, status string) error {
	if !s.isNoteExists(ctx, id) {
		return entity.ErrNoteNotExists
	}

	if title != "" && s.isNoteExists(ctx, title) {
		return entity.ErrNoteExists
	}

	return s.repo.UpdateNote(ctx, id, title, description, status)
}

func (s *Service) DeleteById(ctx context.Context, id int) error {
	if s.isNoteExists(ctx, id) {
		return s.repo.DeleteById(ctx, id)
	}

	return entity.ErrNoteNotExists
}

func (s *Service) DeleteNotes(ctx context.Context) error {
	err := s.repo.DeleteNotes(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) isNoteExists(ctx context.Context, data interface{}) bool {
	var err error
	switch value := data.(type) {
	case int:
		_, err = s.repo.GetById(ctx, value)
	case string:
		_, err = s.repo.GetByTitle(ctx, value)
	}

	return err == nil
}
