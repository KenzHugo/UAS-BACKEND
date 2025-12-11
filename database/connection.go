package database

import (
	"database/sql"
	"fmt"
	"log"
	"UASBE/config"

	_ "github.com/lib/pq"
)

// ConnectPostgres establishes connection to PostgreSQL
func ConnectPostgres(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.GetDSN()
	
	log.Println("🔌 Connecting to PostgreSQL...")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("✅ Connected to PostgreSQL successfully!")
	return db, nil
}