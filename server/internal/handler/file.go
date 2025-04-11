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

// CreateFile godoc
// @Summary Create file
// @Description Create a new file
// @Tags files
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.FileCreateRequest true "File create request"
// @Success 201 {object} model.FileResponse
// @Failure 400,401,403 {object} util.ErrorResponse
// @Router /files [post]
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

// ListFiles godoc
// @Summary List files
// @Description Get a list of files in a bucket with pagination
// @Tags files
// @Accept json
// @Produce json
// @Security Bearer
// @Param bucket_id path int true "Bucket ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} model.FileListResponse
// @Failure 400,401,403 {object} util.ErrorResponse
// @Router /files/bucket/{bucket_id} [get]
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

// GetFile godoc
// @Summary Get file details
// @Description Get details of a specific file
// @Tags files
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "File ID"
// @Success 200 {object} model.FileResponse
// @Failure 400,401,403,404 {object} util.ErrorResponse
// @Router /files/{id} [get]
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

// UpdateFile godoc
// @Summary Update file
// @Description Update file information
// @Tags files
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "File ID"
// @Param request body model.FileUpdateRequest true "File update request"
// @Success 200 {object} model.FileResponse
// @Failure 400,401,403,404 {object} util.ErrorResponse
// @Router /files/{id} [put]
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

// DeleteFile godoc
// @Summary Delete file
// @Description Delete a file
// @Tags files
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "File ID"
// @Success 204 "No Content"
// @Failure 400,401,403,404 {object} util.ErrorResponse
// @Router /files/{id} [delete]
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

// GetUploadURL godoc
// @Summary Get upload URL
// @Description Get a pre-signed URL for file upload
// @Tags files
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "File ID"
// @Success 200 {object} model.FileUploadResponse
// @Failure 400,401,403,404 {object} util.ErrorResponse
// @Router /files/{id}/upload [get]
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

// GetDownloadURL godoc
// @Summary Get download URL
// @Description Get a pre-signed URL for file download
// @Tags files
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "File ID"
// @Success 200 {object} model.FileDownloadResponse
// @Failure 400,401,403,404 {object} util.ErrorResponse
// @Router /files/{id}/download [get]
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
