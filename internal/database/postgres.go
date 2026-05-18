package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"maxsasi/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func New(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migration driver error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/repository/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate instance error: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no new migrations")
			return nil
		}

		return fmt.Errorf("migrate up error: %w", err)
	}

	log.Println("migrations applied successfully")
	return nil
}
