package main

import (
	"fmt"
	"log"
	"go-solid/internal/config"
	"go-solid/internal/database"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Running database migrations...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("mysql", cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	migration := database.NewMigration(db)
	if err := migration.CreateTables(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Database migrations completed successfully!")
}
