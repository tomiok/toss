package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*http.Server
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (s *Server) Start() {
	slog.Info("server is running t", slog.String("port", s.Addr))
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("cannot start the server", slog.Any("err", err))
		}
	}()
	s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT)
	sig := <-quit
	slog.Info("server is shutting down", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		slog.Error("could not gracefully shutdown the server", slog.Any("err", err.Error()))
	}
	slog.Info("server stopped")
}
