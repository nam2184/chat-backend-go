package migrations

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	qm "github.com/nam2184/mymy/models/db"
)

var khanhusername = os.Getenv("khanh_username")
var myusername = os.Getenv("my_username")
var khanhEmail = os.Getenv("khanh_email")
var myEmail = os.Getenv("my_email")
var plainPasswordMy = os.Getenv("my_password")
var plainPasswordKhanh = os.Getenv("khanh_password")

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

	hashedPasswordMy, err := hashPassword(plainPasswordMy)
	if err != nil {
		log.Fatal("failed to hash password for My: ", err)
	}

	hashedPasswordKhanh, err := hashPassword(plainPasswordKhanh)
	if err != nil {
		log.Fatal("failed to hash password for Khanh: ", err)
	}

	user1 := qm.User{
		FirstName: "My",
		Surname:   "Nguyen",
		Username:  myusername,
		Email:     myEmail,
		IsMy:      true,
		IsKhanh:   false,
	}
	auth1 := qm.Auth{
		Username: myusername,
		Password: hashedPasswordMy,
		IsKhanh:  false,
		IsMy:     true,
	}

	user2 := qm.User{
		FirstName: "Khanh",
		Surname:   "Cao",
		Username:  khanhusername,
		Email:     khanhEmail,
		IsMy:      false,
		IsKhanh:   true,
	}
	auth2 := qm.Auth{
		Username: khanhusername,
		Password: hashedPasswordKhanh,
		IsKhanh:  true,
		IsMy:     false,
	}

	chat := qm.Chat{
		User1ID:         1,
		User2ID:         2,
		CreatedAt:       time.Now(),
		Seen:            false,
		UpdatedAt:       time.Now(),
		LastMessageTime: time.Now(),
	}

	db.Create(&user1)
	db.Create(&auth1)
	db.Create(&user2)
	db.Create(&auth2)
	db.Create(&chat)
	log.Println("Database tables migrated successfully!")
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
