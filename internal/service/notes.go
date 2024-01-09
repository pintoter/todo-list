package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

func (s *Service) CreateNote(ctx context.Context, note entity.Note) error {
	if s.isNoteExists(ctx, note.Title, note.UserId) {
		return entity.ErrNoteExists
	}

	_, err := s.repo.CreateNote(ctx, note)
	return err
}

func (s *Service) GetNoteById(ctx context.Context, id, userId int) (entity.Note, error) {
	note, err := s.repo.GetNoteById(ctx, id, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Note{}, entity.ErrNoteNotExists
		} else {
			return entity.Note{}, err
		}
	}

	return note, nil
}

func (s *Service) GetNotes(ctx context.Context, userId int) ([]entity.Note, error) {
	notes, err := s.repo.GetNotes(ctx, userId)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *Service) GetNotesExtended(ctx context.Context, limit, offset int, status string, date time.Time, userId int) ([]entity.Note, error) {
	notes, err := s.repo.GetNotesExtended(ctx, limit, offset, status, date, userId)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *Service) UpdateNote(ctx context.Context, id int, title, description, status string, userId int) error {
	if !s.isNoteExists(ctx, id, userId) {
		return entity.ErrNoteNotExists
	}

	if title != "" && s.isNoteExists(ctx, title, userId) {
		return entity.ErrNoteExists
	}

	return s.repo.UpdateNote(ctx, id, userId, title, description, status)
}

func (s *Service) DeleteNoteById(ctx context.Context, id, userId int) error {
	if s.isNoteExists(ctx, id, userId) {
		return s.repo.DeleteNoteById(ctx, id, userId)
	}

	return entity.ErrNoteNotExists
}

func (s *Service) DeleteNotes(ctx context.Context, userId int) error {
	err := s.repo.DeleteNotes(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) isNoteExists(ctx context.Context, data any, userId int) bool {
	var err error
	switch value := data.(type) {
	case int:
		_, err = s.repo.GetNoteById(ctx, value, userId)
	case string:
		_, err = s.repo.GetNoteByTitle(ctx, value, userId)
	}

	return err == nil
}
