package app

import (
	"log"

	"github.com/pintoter/todo-list/internal/config"
	migrations "github.com/pintoter/todo-list/internal/database"
	"github.com/pintoter/todo-list/internal/service"
	"github.com/pintoter/todo-list/pkg/database/postgres"
)

func Run() {
	cfg := config.Get()

	err := migrations.Do(&cfg.DB)

	db, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	service := service.New(db)

}
