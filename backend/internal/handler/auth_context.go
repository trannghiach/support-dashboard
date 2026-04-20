package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/trannghiach/support-dashboard/backend/internal/middleware"
)

func getAuthUser(c *gin.Context) (int64, string, error) {
	userIDValue, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		return 0, "", errors.New("missing user id in context")
	}

	roleValue, ok := c.Get(middleware.ContextRoleKey)
	if !ok {
		return 0, "", errors.New("missing role in context")
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		return 0, "", errors.New("invalid user id type in context")
	}

	role, ok := roleValue.(string)
	if !ok {
		return 0, "", errors.New("invalid role type in context")
	}

	return userID, role, nil
}