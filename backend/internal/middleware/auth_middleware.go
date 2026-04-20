package middleware

import (
	"net/http"
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/trannghiach/support-dashboard/backend/internal/auth"
	"github.com/trannghiach/support-dashboard/backend/internal/response"
)

const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "role"
)

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(jwtSecret, parts[1])
		if err != nil {
			response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextRoleKey, claims.Role)
		c.Next()
	}
}