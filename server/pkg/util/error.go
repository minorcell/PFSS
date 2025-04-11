package util

import "github.com/gin-gonic/gin"

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Common error codes
const (
	ErrorCodeUnauthorized     = 401
	ErrorCodeForbidden        = 403
	ErrorCodeNotFound         = 404
	ErrorCodeInvalidInput     = 400
	ErrorCodeInternalError    = 500
	ErrorCodeServiceUnavailable = 503
)

// NewError creates a new error response
func NewError(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// SendError sends an error response to the client
func SendError(c *gin.Context, err *ErrorResponse) {
	c.JSON(err.Code, err)
}

// Common errors
var (
	ErrUnauthorized     = NewError(ErrorCodeUnauthorized, "Unauthorized access")
	ErrForbidden        = NewError(ErrorCodeForbidden, "Access forbidden")
	ErrNotFound         = NewError(ErrorCodeNotFound, "Resource not found")
	ErrInvalidInput     = NewError(ErrorCodeInvalidInput, "Invalid input")
	ErrInternalServer   = NewError(ErrorCodeInternalError, "Internal server error")
	ErrServiceUnavailable = NewError(ErrorCodeServiceUnavailable, "Service unavailable")
)
