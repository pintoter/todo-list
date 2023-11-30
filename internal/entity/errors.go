package entity

import "errors"

var (
	ErrNoteNotFound = errors.New("note doesn't exists")
	ErrInvalidInput = errors.New("invalud input parameters")
	ErrNoteExists   = errors.New("note already exists")
)
