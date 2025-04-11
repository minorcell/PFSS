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

// CreateBucket godoc
// @Summary Create bucket
// @Description Create a new bucket
// @Tags buckets
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.BucketCreateRequest true "Bucket create request"
// @Success 201 {object} model.BucketResponse
// @Failure 400,401,403 {object} util.ErrorResponse
// @Router /buckets [post]
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

// ListBuckets godoc
// @Summary List buckets
// @Description Get a list of buckets with pagination
// @Tags buckets
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} model.BucketListResponse
// @Failure 400,401,403 {object} util.ErrorResponse
// @Router /buckets [get]
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

// GetBucket godoc
// @Summary Get bucket details
// @Description Get details of a specific bucket
// @Tags buckets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Bucket ID"
// @Success 200 {object} model.BucketResponse
// @Failure 400,401,403,404 {object} util.ErrorResponse
// @Router /buckets/{id} [get]
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

// UpdateBucket godoc
// @Summary Update bucket
// @Description Update bucket information
// @Tags buckets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Bucket ID"
// @Param request body model.BucketUpdateRequest true "Bucket update request"
// @Success 200 {object} model.BucketResponse
// @Failure 400,401,403,404 {object} util.ErrorResponse
// @Router /buckets/{id} [put]
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

// DeleteBucket godoc
// @Summary Delete bucket
// @Description Delete a bucket
// @Tags buckets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Bucket ID"
// @Success 204 "No Content"
// @Failure 400,401,403,404 {object} util.ErrorResponse
// @Router /buckets/{id} [delete]
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
