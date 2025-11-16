package main

import (
	"log"
	"os"

	"ms-venue-go/internal/handlers"
	"ms-venue-go/internal/middleware"
	"ms-venue-go/internal/repository"
	"ms-venue-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get environment variables
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://bar_user:bar_password@postgres-db:5432/bar_management_db?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")

	// Initialize repository
	venueRepo, err := repository.NewVenueRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize service
	venueService := service.NewVenueService(venueRepo)

	// Initialize handlers
	venueHandler := handlers.NewVenueHandler(venueService)

	// Initialize auth service
	authService := middleware.NewAuthService(jwtSecret)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", venueHandler.HealthCheck)

	// Public endpoints (no auth required)
	public := r.Group("/api/venue")
	{
		// Get all locations (public for frontend)
		public.GET("/locations", venueHandler.GetAllLocations)

		// Get tables by location (public for frontend)
		public.GET("/:locationId/tables", venueHandler.GetTablesByLocation)
	}

	// Protected endpoints (require authentication)
	protected := r.Group("/api/venue")
	protected.Use(middleware.RequireAuth(authService))
	{
		// Location management
		protected.POST("/locations", venueHandler.CreateLocation)
		protected.GET("/locations/:id", venueHandler.GetLocationByID)
		protected.PUT("/locations/:id", venueHandler.UpdateLocation)
		protected.DELETE("/locations/:id", venueHandler.DeleteLocation)

		// Table management
		protected.POST("/tables", venueHandler.CreateTable)
		protected.GET("/tables/:id", venueHandler.GetTableByID)
		protected.PUT("/tables/:id", venueHandler.UpdateTable)
		protected.DELETE("/tables/:id", venueHandler.DeleteTable)
	}

	// Admin-only endpoints
	admin := r.Group("/api/venue/admin")
	admin.Use(middleware.RequireRole(authService, "admin"))
	{
		// Admin can do everything that protected endpoints do
		// This is a placeholder for future admin-specific functionality
	}

	log.Printf("Starting MS-VENUE-GO server on port %s", port)
	log.Fatal(r.Run(":" + port))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
