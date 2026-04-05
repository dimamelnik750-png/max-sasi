package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDB(cfg Config) (*sql.DB, error) {
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

func InitDB(db *sql.DB) error {
	_, err := db.Exec(`
	`)

	return err
}
