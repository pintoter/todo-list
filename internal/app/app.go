package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pintoter/todo-list/internal/config"
	migrations "github.com/pintoter/todo-list/internal/database"
	"github.com/pintoter/todo-list/internal/server"
	"github.com/pintoter/todo-list/internal/service"
	"github.com/pintoter/todo-list/internal/transport"
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

	handler := transport.NewHandler(service)

	server := server.New(handler, &cfg.HTTP)

	server.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	select {
	case s := <-quit:
		log.Printf("Starting gracefully shutdown after signal %s", s.String())
	case err = <-server.Notify():
		log.Fatalf("Error when starting server: %s", err.Error())
	}

	if err := server.Shutdown(); err != nil {
		log.Fatal("Server", err)
	}

	log.Println("Server gracefully shutting down")
}
