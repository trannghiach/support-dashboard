package main

import (
	"log"
	"os"

	"github.com/trannghiach/support-dashboard/backend/internal/db"
	"github.com/trannghiach/support-dashboard/backend/internal/handler"
	"github.com/trannghiach/support-dashboard/backend/internal/service"
	"github.com/trannghiach/support-dashboard/backend/internal/repository"
	"github.com/trannghiach/support-dashboard/backend/internal/router"
	"github.com/trannghiach/support-dashboard/backend/internal/config"
)

func main() {
	config.LoadEnv()

	databaseURL := config.GetEnv("DATABASE_URL", "")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	
	dbPool, err := db.NewPostgres(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	userRepo := repository.NewUserRepository(dbPool)
	userHandler := handler.NewUserHandler(userRepo)

	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	ticketRepo := repository.NewTicketRepository(dbPool)
	ticketService := service.NewTicketService(ticketRepo, userRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	r := router.SetupRouter(authHandler, userHandler, ticketHandler, jwtSecret)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}