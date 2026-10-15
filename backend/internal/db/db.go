package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`,
	`CREATE TABLE IF NOT EXISTS campaigns(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		UniqueStr TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL UNIQUE,
		tenant_id TEXT NOT NULL UNIQUE,
		resource TEXT NOT NULL,
		scopes TEXT NOT NULL CHECK (json_valid(scopes))
	);`,
	`CREATE TABLE IF NOT EXISTS tokens(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tenant_id TEXT NOT NULL,
		access_token TEXT NOT NULL,
		refresh_token TEXT NOT NULL
	);`, // just for testing, so can verify if works
}

func migrate(c *sql.DB) error {
	var v int
	if err := c.QueryRow(`PRAGMA user_version`).Scan(&v); err != nil {
		return err
	}

	for i := v; i < len(migrations); i++ {
		tx, err := c.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(migrations[i]); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration v%d: %w", i+1, err)
		}
		if _, err := tx.Exec(fmt.Sprintf(`PRAGMA user_version = %d`, i+1)); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func OpenAndMigrate(path string) (*sql.DB, error) {
	conn, err := open(path)
	if err != nil {
		return nil, err
	}

	if err := migrate(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

func open(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;`); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil

}

func New() (*DB, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	c, err := open(filepath.Join(h, ".local", "share", "token-jack", "database", "jacker.db"))
	if err != nil {
		return nil, err
	}

	return &DB{conn: c}, nil
}
