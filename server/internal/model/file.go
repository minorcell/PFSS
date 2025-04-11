package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Path         string         `gorm:"size:1024;not null" json:"path"`
	BucketID     uint           `gorm:"not null" json:"bucket_id"`
	Size         int64          `gorm:"not null" json:"size"`
	ContentType  string         `gorm:"size:100;not null" json:"content_type"`
	Hash         string         `gorm:"size:64" json:"hash"`
	Metadata     map[string]string `gorm:"-" json:"metadata,omitempty"`
	CreatedBy    uint           `gorm:"not null" json:"created_by"`
	UpdatedBy    uint           `gorm:"not null" json:"updated_by"`
	LastModified time.Time      `json:"last_modified"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type FileMetadata struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	FileID    uint           `gorm:"not null" json:"file_id"`
	Key       string         `gorm:"size:50;not null" json:"key"`
	Value     string         `gorm:"size:255" json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for File
func (File) TableName() string {
	return "files"
}

// TableName specifies the table name for FileMetadata
func (FileMetadata) TableName() string {
	return "file_metadata"
}
