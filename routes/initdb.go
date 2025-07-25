package routes

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type ServiceDB struct {
	db *sqlx.DB
}

// SetupDBX initializes a connection to an SQLite database.
func SetupDBX() (*sqlx.DB, error) {
	// Define the SQLite file path (replace or get from environment if needed)
	dbFile := "database.db"

	// Open the SQLite connection
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// Set connection pool parameters for SQLite
	db.SetConnMaxIdleTime(5 * time.Minute) // Optional idle time
	db.SetConnMaxLifetime(1 * time.Hour)   // Optional connection lifetime

	return db, nil
}

func SetupTestDBX() (*sqlx.DB, error) {
	// Define the SQLite file path (replace or get from environment if needed)
	dbFile := "test_database.db"

	// Open the SQLite connection
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// Set connection pool parameters for SQLite
	db.SetConnMaxIdleTime(5 * time.Minute) // Optional idle time
	db.SetConnMaxLifetime(1 * time.Hour)   // Optional connection lifetime

	return db, nil
}

// SetupServiceDBX initializes a `ServiceDB` with an SQLite connection.
func SetupServiceDBX() (*ServiceDB, error) {
	// Define the SQLite file path (replace or get from environment if needed)
	dbFile := "database.db"

	// Open the SQLite connection
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// Set connection pool parameters for SQLite
	db.SetConnMaxIdleTime(5 * time.Minute) // Optional idle time
	db.SetConnMaxLifetime(1 * time.Hour)   // Optional connection lifetime

	return &ServiceDB{db: db}, nil
}
