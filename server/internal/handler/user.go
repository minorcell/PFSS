package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minorcell/pfss/internal/model"
	"github.com/minorcell/pfss/internal/service"
	"github.com/minorcell/pfss/pkg/util"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ListUsers handles the user list request
func (h *UserHandler) ListUsers(c *gin.Context) {
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

	// Get users
	users, total, err := h.userService.GetUsers(page, pageSize)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.UserListResponse{
		Users:      users,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	})
}

// GetUser handles the get user request
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	// Get user
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "User not found",
		})
		return
	}

	// Get user permissions
	permissions, err := h.userService.GetUserPermissions(uint(id))
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user permissions",
		})
		return
	}

	c.JSON(http.StatusOK, model.UserResponse{
		User:        user,
		Permissions: permissions,
	})
}

// UpdateUser handles the update user request
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Get current user info from context
	currentUserID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	// Convert request to updates map
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.IsRoot != nil {
		updates["is_root"] = *req.IsRoot
	}

	if err := h.userService.UpdateUser(uint(id), updates, currentUserID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser handles the delete user request
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	// Get current user info from context
	currentUserID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	if err := h.userService.DeleteUser(uint(id), currentUserID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateUserStatus handles the update user status request
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Get current user info from context
	currentUserID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	if err := h.userService.UpdateUserStatus(uint(id), req.Status, currentUserID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

// UpdateUserPermissions handles the update user permissions request
func (h *UserHandler) UpdateUserPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	var req []model.UserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Convert request to permissions
	permissions := make([]model.UserPermission, len(req))
	for i, p := range req {
		permissions[i] = model.UserPermission{
			UserID:   uint(id),
			Resource: p.Resource,
			Action:   p.Action,
		}
	}

	// Get current user info from context
	currentUserID := c.GetUint("user_id")
	isRoot := c.GetBool("is_root")

	if err := h.userService.UpdateUserPermissions(uint(id), permissions, currentUserID, isRoot); err != nil {
		util.SendError(c, &util.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User permissions updated successfully"})
}
