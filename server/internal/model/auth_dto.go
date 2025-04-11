package model

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsRoot   bool   `json:"is_root"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest represents the request body for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword    string `json:"new_password" binding:"required"`
}

// TokenResponse represents the response body for successful login
type TokenResponse struct {
	Token string `json:"token"`
}

// AuthResponse represents the response body for successful registration
type AuthResponse struct {
	ID       uint           `json:"id"`
	Username string         `json:"username"`
	Token    *TokenResponse `json:"token"`
}
