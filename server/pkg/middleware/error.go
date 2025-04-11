package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/minorcell/pfss/pkg/util"
)

// ErrorHandler handles panics and returns appropriate error responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				debug.PrintStack()

				// Return 500 Internal Server Error
				util.SendError(c, &util.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}
		}()

		c.Next()

		// Handle 404 after no handler was matched
		if c.Writer.Status() == 404 {
			util.SendError(c, util.ErrNotFound)
		}
	}
}

// ValidationMiddleware validates request parameters
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBind(c.Request.Body); err != nil {
			util.SendError(c, &util.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request parameters: " + err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
