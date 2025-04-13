package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/minorcell/pfss/pkg/util"
)

// ErrorHandler 处理panic并返回适当的错误响应
// 返回值:
//   - gin.HandlerFunc: Gin中间件函数
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印调用栈信息
				debug.PrintStack()

				// 返回500内部服务器错误
				util.SendError(c, &util.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}
		}()

		c.Next()

		// 处理404未找到路由的情况
		if c.Writer.Status() == http.StatusNotFound {
			util.SendError(c, util.ErrNotFound)
		}
	}
}

// ValidationMiddleware 验证请求参数
// 返回值:
//   - gin.HandlerFunc: Gin中间件函数
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定并验证请求体
		if err := c.ShouldBind(c.Request.Body); err != nil {
			util.SendError(c, &util.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request parameters: " + err.Error(),
			})
			c.Abort() // 终止后续处理
			return
		}
		c.Next()
	}
}
