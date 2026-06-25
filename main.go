package main

import (
	"log"
	"net/http"

	"spotsync/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load configurations
	cfg := config.LoadConfig()

	// 2. Database Connection Placeholder
	// TODO: Initialize database connection using config.DBURL (GORM + Postgres)
	log.Println("Database connection initialized (placeholder)")

	// 3. Initialize Echo Instance
	e := echo.New()

	// 4. Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 5. Health Check Route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "SpotSync API is healthy and running",
		})
	})

	// 6. Start Server
	log.Printf("Starting SpotSync server on port %s...", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
