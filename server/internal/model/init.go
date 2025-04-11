package model

import (
	"log"

	"gorm.io/gorm"
)

// InitDB initializes the database schema and creates necessary tables
func InitDB(db *gorm.DB) error {
	// Auto-migrate all models
	if err := db.AutoMigrate(
		&User{},
		&UserPermission{},
		&File{},
		&FileMetadata{},
		&Bucket{},
		&BucketPermission{},
	); err != nil {
		log.Printf("Failed to auto-migrate database: %v", err)
		return err
	}

	// Check if root user exists
	var count int64
	db.Model(&User{}).Where("is_root = ?", true).Count(&count)
	if count == 0 {
		// Create root user if not exists
		rootUser := User{
			Username: "root",
			Password: "change_me_immediately", // TODO: Hash this password
			IsRoot:   true,
			Status:   "active",
		}
		if err := db.Create(&rootUser).Error; err != nil {
			log.Printf("Failed to create root user: %v", err)
			return err
		}
	}

	return nil
}
