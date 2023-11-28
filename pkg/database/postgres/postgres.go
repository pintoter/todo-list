package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	driverName = "postgres"
)

type Config interface {
	GetDSN() string
}

func New(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(driverName, cfg.GetDSN())
	if err != nil {
		return nil, errors.Wrap(err, "opening db")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping DB")
	}

	return db, nil
}
