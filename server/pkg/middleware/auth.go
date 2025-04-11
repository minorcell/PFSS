package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minorcell/pfss/pkg/util"
)

// AuthMiddleware validates JWT tokens and sets user information in the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.SendError(c, util.ErrUnauthorized)
			c.Abort()
			return
		}

		// Extract token from Bearer schema
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.SendError(c, util.NewError(401, "Invalid authorization header format"))
			c.Abort()
			return
		}

		// Parse and validate token
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			util.SendError(c, util.NewError(401, "Invalid token: "+err.Error()))
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_root", claims.IsRoot)

		c.Next()
	}
}

// RootRequired ensures the user is a root user
func RootRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isRoot, exists := c.Get("is_root")
		if !exists || !isRoot.(bool) {
			util.SendError(c, util.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
