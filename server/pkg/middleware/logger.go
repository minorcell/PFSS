package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 记录详细的请求和响应信息
// 返回值:
//   - gin.HandlerFunc: Gin中间件函数
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 读取请求体内容
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 恢复请求体以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建自定义响应写入器以捕获响应内容
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算请求处理耗时
		latency := time.Since(start)

		// 获取用户信息(如果存在)
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 记录请求详细信息
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

// bodyLogWriter 自定义响应写入器，用于捕获响应体
// 结构体字段:
//   - ResponseWriter: 原始响应写入器
//   - body: 用于存储响应内容的缓冲区
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 实现io.Writer接口，同时写入缓冲区和原始响应
func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// formatLog 格式化日志条目
// 参数:
//   - method: 请求方法
//   - path: 请求路径
//   - clientIP: 客户端IP
//   - userID: 用户ID
//   - username: 用户名
//   - status: 响应状态码
//   - latency: 请求处理耗时
//   - requestBody: 请求体内容
//   - responseBody: 响应体内容
//
// 返回值:
//   - string: 格式化后的日志字符串
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

// formatInterface 格式化接口值为字符串
// 参数:
//   - v: 需要格式化的值
//
// 返回值:
//   - string: 格式化后的字符串
func formatInterface(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", v)
}
