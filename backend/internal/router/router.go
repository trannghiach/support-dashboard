package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trannghiach/support-dashboard/backend/internal/handler"
	"github.com/trannghiach/support-dashboard/backend/internal/middleware"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	ticketHandler *handler.TicketHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API is healthy",
		})
	})

	r.POST("/auth/login", authHandler.Login)

	// public dev routes
	r.GET("/users", userHandler.GetAll)
	r.GET("/tickets", ticketHandler.GetTickets)
	r.GET("/tickets/:id/replies", ticketHandler.GetReplies) 

	authorized := r.Group("/")
	authorized.Use(middleware.RequireAuth(jwtSecret)) 
	{
		authorized.POST("/tickets", ticketHandler.CreateTicket)
		authorized.PATCH("/tickets/:id/status", ticketHandler.UpdateTicketStatus)
		authorized.POST("/tickets/:id/replies", ticketHandler.CreateReply) 
		authorized.PATCH("/tickets/:id/assign", ticketHandler.AssignTicket)
	}
	
	return r
}