package model

import (
	"time"

	"gorm.io/gorm"
)

type Bucket struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:63;unique;not null" json:"name"`
	OwnerID     uint           `gorm:"not null" json:"owner_id"`
	Description string         `gorm:"size:255" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type BucketPermission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	BucketID  uint           `gorm:"not null" json:"bucket_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Access    string         `gorm:"size:20;not null" json:"access"` // read, write, admin
	ExpiresAt *time.Time     `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Bucket
func (Bucket) TableName() string {
	return "buckets"
}

// TableName specifies the table name for BucketPermission
func (BucketPermission) TableName() string {
	return "bucket_permissions"
}
