package model

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	IsRoot   bool   `json:"is_root"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"` // Expiration time in seconds
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  *User         `json:"user"`
	Token *TokenResponse `json:"token,omitempty"`
}
