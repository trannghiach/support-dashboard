package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"http://localhost:3000", // Next.js frontend
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))


	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API is healthy",
		})
	})

	r.POST("/auth/login", authHandler.Login)

	// public dev routes
	r.GET("/users", userHandler.GetAll)

	authorized := r.Group("/")
	authorized.Use(middleware.RequireAuth(jwtSecret)) 
	{
		authorized.POST("/tickets", ticketHandler.CreateTicket)
		authorized.PATCH("/tickets/:id/status", ticketHandler.UpdateTicketStatus)
		authorized.POST("/tickets/:id/replies", ticketHandler.CreateReply) 
		authorized.PATCH("/tickets/:id/assign", ticketHandler.AssignTicket)
		authorized.GET("/tickets/:id", ticketHandler.GetTicketByID)
		authorized.GET("/tickets", ticketHandler.GetTickets)
		authorized.GET("/tickets/:id/replies", ticketHandler.GetReplies) 
		authorized.POST("/tickets/:id/ai-assist", ticketHandler.GenerateAIAssist)
	}
	
	return r
}