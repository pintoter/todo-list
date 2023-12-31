package entity

import "errors"

var (
	ErrNoteExists    = errors.New("note already exists")
	ErrNoteNotExists = errors.New("note doesn't exist")
	ErrInvalidAuth   = errors.New("missing authorization header")
	ErrInvalidDate   = errors.New("invalid date")
	ErrInvalidEmail  = errors.New("invalid email")
	ErrInvalidId     = errors.New("invalid id")
	ErrInvalidInput  = errors.New("invalid input parameters")
	ErrInvalidPage   = errors.New("invalid page")
	ErrInvalidStatus = errors.New("invalid status")
	ErrUserExists    = errors.New("user with input parameters already exists")
)
