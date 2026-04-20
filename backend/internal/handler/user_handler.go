package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trannghiach/support-dashboard/backend/internal/repository"
	"github.com/trannghiach/support-dashboard/backend/internal/response"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	
	c.JSON(http.StatusOK, users)
}
