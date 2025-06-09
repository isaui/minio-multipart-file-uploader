package config

import (
	"log"
	"os"

	"backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() (*gorm.DB, error) {
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "./database.db"
	}

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to the database
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.File{}); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return nil, err
	}

	DB = db
	return db, nil
}
