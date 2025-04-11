package model

// BucketCreateRequest represents the bucket creation request
type BucketCreateRequest struct {
	Name        string            `json:"name" binding:"required,min=3,max=63"`
	Description string            `json:"description" binding:"max=255"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// BucketUpdateRequest represents the bucket update request
type BucketUpdateRequest struct {
	Name        string            `json:"name,omitempty" binding:"omitempty,min=3,max=63"`
	Description string            `json:"description,omitempty" binding:"max=255"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// BucketPermissionRequest represents the bucket permission request
type BucketPermissionRequest struct {
	UserID    uint      `json:"user_id" binding:"required"`
	Access    string    `json:"access" binding:"required,oneof=read write admin"`
	ExpiresAt *JSONTime `json:"expires_at,omitempty"`
}

// BucketResponse represents the bucket response with permissions
type BucketResponse struct {
	Bucket      *Bucket            `json:"bucket"`
	Permissions []BucketPermission `json:"permissions,omitempty"`
}

// BucketListResponse represents the paginated bucket list response
type BucketListResponse struct {
	Buckets    []Bucket `json:"buckets"`
	TotalCount int64    `json:"total_count"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
}

// BucketStats represents bucket statistics
type BucketStats struct {
	FileCount int64 `json:"file_count"`
	TotalSize int64 `json:"total_size"` // in bytes
}
