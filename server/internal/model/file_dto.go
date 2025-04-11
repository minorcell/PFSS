package model

// FileCreateRequest represents the file creation request
type FileCreateRequest struct {
	BucketID    uint              `json:"bucket_id" binding:"required"`
	Name        string            `json:"name" binding:"required,min=1,max=255"`
	Path        string            `json:"path" binding:"required"`
	ContentType string            `json:"content_type" binding:"required"`
	Size        int64             `json:"size" binding:"required"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// FileUpdateRequest represents the file update request
type FileUpdateRequest struct {
	Name        string            `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Path        string            `json:"path,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// FileResponse represents the file response
type FileResponse struct {
	ID          uint              `json:"id"`
	BucketID    uint              `json:"bucket_id"`
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	ContentType string            `json:"content_type"`
	Size        int64             `json:"size"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   JSONTime         `json:"created_at"`
	UpdatedAt   JSONTime         `json:"updated_at"`
}

// FileListResponse represents the paginated file list response
type FileListResponse struct {
	Files      []FileResponse `json:"files"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}

// FileUploadResponse represents the file upload response
type FileUploadResponse struct {
	File      *FileResponse `json:"file"`
	UploadURL string       `json:"upload_url"`
}

// FileDownloadResponse represents the file download response
type FileDownloadResponse struct {
	File        *FileResponse `json:"file"`
	DownloadURL string       `json:"download_url"`
	ExpiresAt   JSONTime     `json:"expires_at"`
}
