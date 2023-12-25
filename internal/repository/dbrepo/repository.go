package dbrepo

import "database/sql"

const (
	notes = "notes"
	users = "users"
)

type DBRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}
