package service

import (
	"errors"

	"github.com/minorcell/pfss/internal/model"
	"gorm.io/gorm"
)

// UserService handles user-related operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUsers returns a list of users with pagination
func (s *UserService) GetUsers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Get total count
	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}, currentUserID uint, isRoot bool) error {
	// Get the user to update
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// Check permissions
	if !isRoot && currentUserID != id {
		return errors.New("permission denied")
	}

	// Prevent non-root users from modifying root status
	if !isRoot {
		delete(updates, "is_root")
	}

	// Prevent modification of sensitive fields
	delete(updates, "password")
	delete(updates, "id")

	return s.db.Model(user).Updates(updates).Error
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uint, currentUserID uint, isRoot bool) error {
	// Get the user to delete
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// Check permissions
	if !isRoot {
		if currentUserID != id {
			return errors.New("permission denied")
		}
		if user.IsRoot {
			return errors.New("cannot delete root user")
		}
	} else {
		// Root users cannot delete themselves
		if currentUserID == id {
			return errors.New("root user cannot delete themselves")
		}
	}

	return s.db.Delete(user).Error
}

// UpdateUserStatus updates a user's status (active/inactive)
func (s *UserService) UpdateUserStatus(id uint, status string, currentUserID uint, isRoot bool) error {
	// Only root users can change user status
	if !isRoot {
		return errors.New("permission denied")
	}

	// Root users cannot deactivate themselves
	if currentUserID == id {
		return errors.New("root user cannot change their own status")
	}

	// Validate status
	if status != "active" && status != "inactive" {
		return errors.New("invalid status")
	}

	return s.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// GetUserPermissions returns a list of permissions for a user
func (s *UserService) GetUserPermissions(userID uint) ([]model.UserPermission, error) {
	var permissions []model.UserPermission
	if err := s.db.Where("user_id = ?", userID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// UpdateUserPermissions updates a user's permissions
func (s *UserService) UpdateUserPermissions(userID uint, permissions []model.UserPermission, currentUserID uint, isRoot bool) error {
	// Only root users can modify permissions
	if !isRoot {
		return errors.New("permission denied")
	}

	// Start transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing permissions
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserPermission{}).Error; err != nil {
			return err
		}

		// Add new permissions
		for _, perm := range permissions {
			perm.UserID = userID
			if err := tx.Create(&perm).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
