package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/pintoter/todo-list/docs"
	"github.com/pintoter/todo-list/internal/config"
	migrations "github.com/pintoter/todo-list/internal/database"
	"github.com/pintoter/todo-list/internal/server"
	"github.com/pintoter/todo-list/internal/service"
	"github.com/pintoter/todo-list/internal/transport"
	"github.com/pintoter/todo-list/pkg/database/postgres"
)

// @title           			todo-list
// @version         			1.0
// @description     			REST API for TODO app

// @contact.name   				Vlad Yurasov
// @contact.email  				meine23@yandex.ru

// @host      					localhost:8080
// @BasePath  					/api/v1

func Run() {
	cfg := config.Get()

	err := migrations.Do(&cfg.DB)
	if err != nil {
		log.Fatal("migrations failed:", err)
	}

	db, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected")

	service := service.New(db)

	handler := transport.NewHandler(service)

	server := server.New(&cfg.HTTP, handler)

	server.Run()
	log.Println("Starting server")

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
