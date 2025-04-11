package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minorcell/pfss/internal/model"
	"github.com/minorcell/pfss/internal/service"
	"github.com/minorcell/pfss/pkg/util"
)

// BucketHandler handles bucket-related requests
type BucketHandler struct {
	bucketService *service.BucketService
}

// NewBucketHandler creates a new bucket handler
func NewBucketHandler(bucketService *service.BucketService) *BucketHandler {
	return &BucketHandler{
		bucketService: bucketService,
	}
}

// CreateBucket handles bucket creation
func (h *BucketHandler) CreateBucket(c *gin.Context) {
	var req model.BucketCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id")
	bucket, err := h.bucketService.CreateBucket(&req, userID)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, bucket)
}

// ListBuckets handles bucket listing
func (h *BucketHandler) ListBuckets(c *gin.Context) {
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

	buckets, total, err := h.bucketService.GetBuckets(userID, isRoot, page, pageSize)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get buckets: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.BucketListResponse{
		Buckets:    buckets,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	})
}

// GetBucket handles getting a single bucket
func (h *BucketHandler) GetBucket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid bucket ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	bucket, err := h.bucketService.GetBucketByID(uint(id), userID, isRoot)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	permissions, err := h.bucketService.GetBucketPermissions(uint(id), userID, isRoot)
	if err != nil {
		// Don't fail the request if permissions can't be retrieved
		c.JSON(http.StatusOK, model.BucketResponse{
			Bucket: bucket,
		})
		return
	}

	c.JSON(http.StatusOK, model.BucketResponse{
		Bucket:      bucket,
		Permissions: permissions,
	})
}

// UpdateBucket handles bucket updates
func (h *BucketHandler) UpdateBucket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid bucket ID",
		})
		return
	}

	var req model.BucketUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	if err := h.bucketService.UpdateBucket(uint(id), &req, userID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bucket updated successfully"})
}

// DeleteBucket handles bucket deletion
func (h *BucketHandler) DeleteBucket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid bucket ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	if err := h.bucketService.DeleteBucket(uint(id), userID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bucket deleted successfully"})
}

// UpdateBucketPermissions handles updating bucket permissions
func (h *BucketHandler) UpdateBucketPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid bucket ID",
		})
		return
	}

	var req []model.BucketPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	if err := h.bucketService.UpdateBucketPermissions(uint(id), req, userID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bucket permissions updated successfully"})
}

// GetBucketStats handles getting bucket statistics
func (h *BucketHandler) GetBucketStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid bucket ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	stats, err := h.bucketService.GetBucketStats(uint(id), userID, isRoot)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
