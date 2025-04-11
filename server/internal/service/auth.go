package service

import (
	"errors"

	"github.com/minorcell/pfss/internal/model"
	"github.com/minorcell/pfss/pkg/util"
	"gorm.io/gorm"
)

// AuthService handles authentication related operations
type AuthService struct {
	db *gorm.DB
}

// NewAuthService creates a new authentication service
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	var user model.User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Validate password
	if !util.ValidatePassword(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Generate token
	token, err := util.GenerateToken(user.ID, user.Username, user.IsRoot)
	if err != nil {
		return nil, err
	}

	// Create response
	tokenResp := &model.TokenResponse{
		Token: token,
	}
	return &model.AuthResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    tokenResp,
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if username already exists
	var count int64
	if err := s.db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		IsRoot:   req.IsRoot,
		Status:   "active",
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// Generate token
	token, err := util.GenerateToken(user.ID, user.Username, user.IsRoot)
	if err != nil {
		return nil, err
	}

	// Create response
	tokenResp := &model.TokenResponse{
		Token: token,
	}
	return &model.AuthResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    tokenResp,
	}, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Validate current password
	if !util.ValidatePassword(currentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	return s.db.Model(&user).Update("password", hashedPassword).Error
}
