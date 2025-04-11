package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;unique;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	IsRoot    bool          `gorm:"default:false" json:"is_root"`
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserPermission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Resource  string         `gorm:"size:50;not null" json:"resource"`
	Action    string         `gorm:"size:20;not null" json:"action"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// TableName specifies the table name for UserPermission
func (UserPermission) TableName() string {
	return "user_permissions"
}
