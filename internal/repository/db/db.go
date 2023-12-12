package repository

import "database/sql"

const (
	notes = "notes"
	users = "users"
)

type DBRepo struct {
	db *sql.DB
}

func NewDBRepo(db *sql.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}
