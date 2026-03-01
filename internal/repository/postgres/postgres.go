package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"filmDb/pkg/modules"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB *sqlx.DB
}

func NewStorage(ctx context.Context, cfg *modules.PostgreConfig) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.HOST, cfg.Port, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = AutoMigrate(dsn); err != nil {
		return nil, fmt.Errorf("migrations failed : %w", err)
	}

	return &Storage{
		DB: db,
	}, nil
}

func AutoMigrate(dsn string) error {
	sourceURL := "file://database/migrations"

	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("no new migrations")
		return nil
	}

	return nil
}
