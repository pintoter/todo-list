package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pintoter/todo-list/internal/repository/dbrepo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/pintoter/todo-list/docs"
	"github.com/pintoter/todo-list/internal/config"
	migrations "github.com/pintoter/todo-list/internal/database"
	"github.com/pintoter/todo-list/internal/server"
	"github.com/pintoter/todo-list/internal/service"
	"github.com/pintoter/todo-list/internal/transport"
	"github.com/pintoter/todo-list/pkg/auth"
	"github.com/pintoter/todo-list/pkg/database/postgres"
	"github.com/pintoter/todo-list/pkg/hash"
	"github.com/pintoter/todo-list/pkg/logger"
)

// @title           			todo-list
// @version         			1.0
// @description     			REST API for TODO app

// @contact.name   				Vlad Yurasov
// @contact.email  				meine23@yandex.ru

// @host      					localhost:8080
// @BasePath  					/api/v1

func Run() {
	ctx := context.Background()

	cfg := config.Get()

	syncLogger := initLogger(ctx, cfg)
	defer syncLogger()

	err := migrations.Do(&cfg.DB)
	if err != nil {
		logger.FatalKV(ctx, "Failed init migrations", "err", err)
	}

	db, err := postgres.New(&cfg.DB)
	if err != nil {
		logger.FatalKV(ctx, "Failed connect database", "err", err)
	}

	deps := service.Deps{
		Cfg:          &cfg.Auth,
		Repo:         dbrepo.New(db),
		Hasher:       hash.New(&cfg.Auth),
		TokenManager: auth.NewManager(&cfg.Auth),
	}

	service := service.New(deps)
	handler := transport.NewHandler(service)
	server := server.New(&cfg.HTTP, handler)

	server.Run()
	log.Println("Starting server")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	select {
	case <-quit:
		logger.InfoKV(ctx, "Starting gracefully shutdown")
	case err = <-server.Notify():
		logger.FatalKV(ctx, "Failed starting server", "err", err.Error())
	}

	if err := server.Shutdown(); err != nil {
		logger.FatalKV(ctx, "Failed shutdown server", "err", err.Error())
	}
}

func initLogger(ctx context.Context, cfg *config.Config) (syncFn func()) {
	loggingLevel := zap.InfoLevel
	if cfg.Project.Level == logger.DebugLevel {
		loggingLevel = zap.DebugLevel
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		zap.NewAtomicLevelAt(loggingLevel),
	)

	notSuggaredLogger := zap.New(consoleCore)

	sugarLogger := notSuggaredLogger.Sugar()
	logger.SetLogger(sugarLogger.With(
		"service", cfg.Project.Name,
	))

	return func() {
		notSuggaredLogger.Sync()
	}
}
