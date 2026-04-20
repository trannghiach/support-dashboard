package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trannghiach/support-dashboard/backend/internal/dto"
	"github.com/trannghiach/support-dashboard/backend/internal/service"
	"github.com/trannghiach/support-dashboard/backend/internal/response"
)

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandler(service *service.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req dto.CreateTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request payload")
		return
	}

	userID, role, err := getAuthUser(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	ticket, err := h.service.CreateTicket(
		c.Request.Context(),
		userID,
		role,
		req,
	)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"data": ticket,
	})
}

func (h *TicketHandler) GetTickets(c *gin.Context) {
	status := c.Query("status")
	priority := c.Query("priority")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tickets, count, err := h.service.ListTickets(
		c.Request.Context(),
		status,
		priority,
		page,
		limit,
	)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch tickets")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tickets,
		"meta": gin.H{
			"page": page,
			"limit": limit,
			"count": count,
		},
	})
}

func (h *TicketHandler) UpdateTicketStatus(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid ticket id")
		return
	}

	var req dto.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	userID, role, err := getAuthUser(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	tickets, err := h.service.UpdateTicketStatus(
		c.Request.Context(),
		id,
		userID,
		role,
		req,
	)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tickets,
	})
}

func (h *TicketHandler) CreateReply(c *gin.Context) {
	idParam := c.Param("id")

	ticketID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid ticket id")
		return
	}
	
	var req dto.CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	userID, role, err := getAuthUser(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}
	
	reply, err := h.service.CreateReply(
		c.Request.Context(),
		ticketID,
		userID,
		role,
		req,
	)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": reply,
	})

}

func (h *TicketHandler) GetReplies(c *gin.Context) {
	idParam := c.Param("id")

	ticketID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid ticket id")
		return
	}

	replies, err := h.service.GetReplies(
		c.Request.Context(),
		ticketID,
	)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch replies")
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": replies,
	})
}

func (h *TicketHandler) AssignTicket(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid ticket id")
		return
	}
	
	var req dto.AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	userID, role, err := getAuthUser(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}
	
	ticket, err := h.service.AssignTicket(
		c.Request.Context(),
		id,
		userID,
		role,
		req,
	)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ticket,
	})
}

func (h *TicketHandler) GetTicketByID(c *gin.Context) {
	idParam := c.Param("id")
	
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid ticket id")
		return
	}
	
	userID, role, err := getAuthUser(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}
	
	ticket, err := h.service.GetTicketByID(
		c.Request.Context(),
		id,
		userID,
		role,
	)
	if err != nil {
		switch err.Error() {
		case "ticket not found":
			response.JSONError(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		case "customers can only view their own tickets", "agents can only view tickets assigned to them":
			response.JSONError(c, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			response.JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": ticket,
	})
}