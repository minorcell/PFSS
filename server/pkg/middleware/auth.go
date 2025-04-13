package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minorcell/pfss/pkg/util"
)

// AuthMiddleware JWT认证中间件，验证Token并设置用户信息到上下文
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.SendError(c, util.ErrUnauthorized)
			c.Abort()
			return
		}

		// 从Bearer格式中提取Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.SendError(c, util.NewError(401, "Invalid authorization header format"))
			c.Abort()
			return
		}

		// 解析并验证Token
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			util.SendError(c, util.NewError(401, "Invalid token: "+err.Error()))
			c.Abort()
			return
		}

		// 将用户信息设置到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_root", claims.IsRoot)

		c.Next()
	}
}

// RootRequired 验证用户是否为管理员
func RootRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取管理员标识
		isRoot, exists := c.Get("is_root")
		if !exists || !isRoot.(bool) {
			util.SendError(c, util.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
