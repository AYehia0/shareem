package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"shareem/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

// define the database connection method

var (
	ErrorMissingMigrationPath = errors.New("MIGRATION_PATH env is missing")
	ErrorMissingDatabaseURL   = errors.New("DATABASE_URL env is missing")
)

// load the configs from the url
func loadConfigsFromEnvURL() (*pgxpool.Config, error) {
	url, ok := os.LookupEnv(config.DATABASE_URL)
	if !ok {
		return nil, ErrorMissingDatabaseURL
	}

	config, err := pgxpool.ParseConfig(url)

	if err != nil {
		return nil, err
	}

	return config, nil
}

// load the configs from the env path, if not found, load from the url
func loadConfigs() (*pgxpool.Config, error) {
	cfg, err := config.NewDatabaseConfig()

	if err != nil {
		return loadConfigsFromEnvURL()
	}

	return pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	))
}

func dbURL() (string, error) {
	cfg, err := config.NewDatabaseConfig()
	if err != nil {
		url, ok := os.LookupEnv(config.DATABASE_URL)
		if !ok {
			return "", ErrorMissingDatabaseURL
		}
		return url, nil
	}
	return cfg.URL(), nil
}

// New creates a new database connection pool using the provided configuration
func Connect(ctx context.Context, logger *slog.Logger) (*pgxpool.Pool, error) {
	cfg, err := loadConfigs()

	if err != nil {
		return nil, err
	}

	db, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	logger.Info("connected to database")

	// run the migrations
	migrationPath, ok := os.LookupEnv(config.MIGRATION_PATH)
	if !ok {
		migrationPath = "file://database/migrations"
	}

	url, err := dbURL()
	if err != nil {
		return nil, err
	}

	migrator, err := migrate.New(migrationPath, url)

	if err != nil {
		return nil, err
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("ran database migrations")

	return db, nil
}
