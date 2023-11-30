package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer      *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

type HTTPConfig interface {
	GetAddr() string
	GetShutdownTimeout() time.Duration
}

func New(handler http.Handler, cfg HTTPConfig) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.GetAddr(),
			Handler: handler,
		},
		notify:          make(chan error),
		shutdownTimeout: cfg.GetShutdownTimeout(),
	}
}

func (s *Server) Run() {
	go func() {
		s.notify <- s.httpServer.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
