package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs detailed request and response information
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read the request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the body for later use
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a custom response writer to capture the response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get user information if available
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// Log the request details
		gin.DefaultWriter.Write([]byte(formatLog(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			userID,
			username,
			c.Writer.Status(),
			latency,
			string(requestBody),
			blw.body.String(),
		)))
	}
}

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// formatLog formats the log entry
func formatLog(method, path, clientIP interface{}, userID, username interface{}, status int, latency time.Duration, requestBody, responseBody string) string {
	return time.Now().Format("2006/01/02 - 15:04:05") + " | " +
		"Method: " + method.(string) + " | " +
		"Path: " + path.(string) + " | " +
		"IP: " + clientIP.(string) + " | " +
		"UserID: " + formatInterface(userID) + " | " +
		"Username: " + formatInterface(username) + " | " +
		"Status: " + string(rune(status)) + " | " +
		"Latency: " + latency.String() + "\n" +
		"Request: " + requestBody + "\n" +
		"Response: " + responseBody + "\n" +
		"--------------------------------------------------\n"
}

// formatInterface formats an interface for logging
func formatInterface(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", v)
}
