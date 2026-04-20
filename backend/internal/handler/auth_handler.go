package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trannghiach/support-dashboard/backend/internal/dto"
	"github.com/trannghiach/support-dashboard/backend/internal/service"
	"github.com/trannghiach/support-dashboard/backend/internal/response"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request payload")
		return
	}

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.LoginResponse{
			Token: token,
		},
	})
}