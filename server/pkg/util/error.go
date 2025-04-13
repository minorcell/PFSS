package util

import "github.com/gin-gonic/gin"

// 错误响应结构体
// 用于定义错误响应的格式
// 包含错误码和错误信息
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 错误码常量
// 用于定义错误码的常量值
// 方便在代码中使用
const (
	ErrorCodeUnauthorized     = 401
	ErrorCodeForbidden        = 403
	ErrorCodeNotFound         = 404
	ErrorCodeInvalidInput     = 400
	ErrorCodeInternalError    = 500
	ErrorCodeServiceUnavailable = 503
)

// 创建错误响应
func NewError(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// 发送错误响应
// 将错误响应以 JSON 格式发送给客户端
// 参数 c 是 Gin 上下文对象，用于发送响应
// 参数 err 是错误响应对象，包含错误码和错误信息
// 该函数会将错误响应以 JSON 格式发送给客户端
// 客户端可以根据错误码和错误信息进行相应的处理
func SendError(c *gin.Context, err *ErrorResponse) {
	c.JSON(err.Code, err)
}

// 定义一些常用的错误响应
// 用于在代码中直接使用这些错误响应
var (
	ErrUnauthorized     = NewError(ErrorCodeUnauthorized, "Unauthorized access")
	ErrForbidden        = NewError(ErrorCodeForbidden, "Access forbidden")
	ErrNotFound         = NewError(ErrorCodeNotFound, "Resource not found")
	ErrInvalidInput     = NewError(ErrorCodeInvalidInput, "Invalid input")
	ErrInternalServer   = NewError(ErrorCodeInternalError, "Internal server error")
	ErrServiceUnavailable = NewError(ErrorCodeServiceUnavailable, "Service unavailable")
)
