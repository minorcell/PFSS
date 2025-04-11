package model

// UserUpdateRequest represents the user update request
type UserUpdateRequest struct {
	Username string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Status   string `json:"status,omitempty" binding:"omitempty,oneof=active inactive"`
	IsRoot   *bool  `json:"is_root,omitempty"`
}

// UserPermissionRequest represents the permission update request
type UserPermissionRequest struct {
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required,oneof=read write admin"`
}

// UserListResponse represents the paginated user list response
type UserListResponse struct {
	Users      []User `json:"users"`
	TotalCount int64  `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

// UserResponse represents the user response with permissions
type UserResponse struct {
	User        *User             `json:"user"`
	Permissions []UserPermission `json:"permissions,omitempty"`
}
