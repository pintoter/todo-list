package entity

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input parameters")
	ErrInvalidId     = errors.New("invalid id")
	ErrNoteExists    = errors.New("note already exists")
	ErrNoteNotExists = errors.New("note doesn't exist")
	ErrInvalidPage   = errors.New("invalid page")
	ErrInvalidDate   = errors.New("invalid date")
	ErrInvalidStatus = errors.New("invalid status")
)
