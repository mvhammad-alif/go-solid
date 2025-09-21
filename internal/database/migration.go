package database

import (
	"database/sql"
	"fmt"
)

type Migration struct {
	db *sql.DB
}

func NewMigration(db *sql.DB) *Migration {
	return &Migration{db: db}
}

func (m *Migration) CreateTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS posts (
			id INT PRIMARY KEY,
			user_id INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			body TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,
	}

	for i, query := range tables {
		if _, err := m.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table %d: %w", i+1, err)
		}
	}

	return nil
}

func (m *Migration) DropTables() error {
	tables := []string{"posts", "users"}

	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
		if _, err := m.db.Exec(query); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}
