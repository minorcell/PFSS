package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minorcell/pfss/internal/model"
	"github.com/minorcell/pfss/internal/service"
	"github.com/minorcell/pfss/pkg/util"
)

// FileHandler handles file-related requests
type FileHandler struct {
	fileService *service.FileService
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// CreateFile handles file creation
func (h *FileHandler) CreateFile(c *gin.Context) {
	var req model.FileCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	file, err := h.fileService.CreateFile(&req, userID, isRoot)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, file)
}

// ListFiles handles file listing
func (h *FileHandler) ListFiles(c *gin.Context) {
	bucketID, err := strconv.ParseUint(c.Param("bucket_id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid bucket ID",
		})
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	files, total, err := h.fileService.GetFiles(uint(bucketID), userID, isRoot, page, pageSize)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get files: " + err.Error(),
		})
		return
	}

	// Convert files to response format
	fileResponses := make([]model.FileResponse, len(files))
	for i, file := range files {
		fileResponses[i] = model.FileResponse{
			ID:          file.ID,
			BucketID:    file.BucketID,
			Name:        file.Name,
			Path:        file.Path,
			ContentType: file.ContentType,
			Size:        file.Size,
			Metadata:    file.Metadata,
			CreatedAt:   model.JSONTime(file.CreatedAt),
			UpdatedAt:   model.JSONTime(file.UpdatedAt),
		}
	}

	c.JSON(http.StatusOK, model.FileListResponse{
		Files:      fileResponses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	})
}

// GetFile handles getting a single file
func (h *FileHandler) GetFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid file ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	file, err := h.fileService.GetFileByID(uint(id), userID, isRoot)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.FileResponse{
		ID:          file.ID,
		BucketID:    file.BucketID,
		Name:        file.Name,
		Path:        file.Path,
		ContentType: file.ContentType,
		Size:        file.Size,
		Metadata:    file.Metadata,
		CreatedAt:   model.JSONTime(file.CreatedAt),
		UpdatedAt:   model.JSONTime(file.UpdatedAt),
	})
}

// UpdateFile handles file updates
func (h *FileHandler) UpdateFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid file ID",
		})
		return
	}

	var req model.FileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	if err := h.fileService.UpdateFile(uint(id), &req, userID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File updated successfully"})
}

// DeleteFile handles file deletion
func (h *FileHandler) DeleteFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid file ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	if err := h.fileService.DeleteFile(uint(id), userID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// GetUploadURL handles getting a pre-signed upload URL
func (h *FileHandler) GetUploadURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid file ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	uploadURL, err := h.fileService.GetUploadURL(uint(id), userID, isRoot)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.FileUploadResponse{
		File:      nil, // TODO: Add file details if needed
		UploadURL: uploadURL,
	})
}

// GetDownloadURL handles getting a pre-signed download URL
func (h *FileHandler) GetDownloadURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid file ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	downloadURL, expiresAt, err := h.fileService.GetDownloadURL(uint(id), userID, isRoot)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.FileDownloadResponse{
		File:        nil, // TODO: Add file details if needed
		DownloadURL: downloadURL,
		ExpiresAt:   model.JSONTime(expiresAt),
	})
}
