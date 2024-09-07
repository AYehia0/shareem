package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"shareem/internal/database"
	"shareem/internal/middleware"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	logger *slog.Logger
	router *http.ServeMux
	db     *pgxpool.Pool
}

func New(logger *slog.Logger) *App {
	return &App{
		logger: logger,
		router: http.NewServeMux(),
	}
}

func (a *App) Run(ctx context.Context) error {

	db, err := database.Connect(ctx, a.logger)

	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	a.db = db

	a.reloadRoutes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(a.logger, a.router),
	}

	done := make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("server error", slog.Any("Error", err))
		}
		close(done)
	}()

	a.logger.Info("server started", slog.String("address", server.Addr))

	select {
	case <-done:
		break

	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		server.Shutdown(ctx)
		cancel()
	}

	return nil

}
