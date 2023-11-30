package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

const (
	format = "2006-01-02"
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

	var notes []entity.Note
	var err error

	switch {
	case date != "" && status != "":
		var dateFormatted time.Time
		var err error
		if date != "" {
			dateFormatted, err = time.Parse(format, date)
			if err != nil {
				return nil, entity.ErrInvalidDate
			}
		}

		notes, err = s.repo.GetNotesByStatusAndDate(ctx, limit, offset, status, dateFormatted)
	case status != "":
		notes, err = s.repo.GetNotesByStatus(ctx, limit, offset, status)
	default:
		notes, err = s.repo.GetNotes(ctx, limit, offset)
	}

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
