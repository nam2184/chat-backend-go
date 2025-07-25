package migrations

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	qm "github.com/nam2184/mymy/models/db"
)

// Helper function to remove the existing SQLite file
func removeFileIfExists(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			log.Fatalf("Failed to remove existing database file: %v", err)
		}
	}
}

func TestAutoMigration(t *testing.T) {
	dbFile := "../database.db"

	// Remove the existing SQLite file
	removeFileIfExists(dbFile)

	// Open SQLite database connection (creates a new file if it doesn't exist)
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}

	err = autoMigrate(db)
	if err != nil {
		log.Fatalf("AutoMigration failed: %v", err)
	}

	log.Println("AutoMigration completed successfully!")
}

func autoMigrate(db *gorm.DB) error {
	// Auto migrate all tables based on the models
	err := db.AutoMigrate(
		&qm.Chat{},
		&qm.Message{},
		&qm.Auth{},
		&qm.User{},
		&qm.EncryptedMessage{},
		&qm.EvaluationKeys{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database tables: %w", err)
	}
	return nil
}
