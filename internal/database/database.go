// Package database предоставляет управление подключением к базе данных PostgreSQL.
package database

import (
	"database/sql"
	"fmt"
	"time"

	// регистрация драйвера
	_ "github.com/lib/pq"
)

// DB представляет подключение к базе данных
type DB struct {
	*sql.DB
}

// New создает новое подключение к базе данных
func New(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Close закрывает подключение к базе данных
func (db *DB) Close() error {
	return db.DB.Close()
}
