package main

import (
	"log"
	"net/http"

	"spotsync/config"
	"spotsync/models"
	"spotsync/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load configurations
	cfg := config.LoadConfig()

	// 2. Initialize Database and AutoMigrate models
	db := config.InitDB(cfg)
	log.Println("Running AutoMigration...")
	err := db.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{})
	if err != nil {
		log.Fatalf("AutoMigration failed: %v", err)
	}
	log.Println("AutoMigration completed successfully")

	// Avoid unused db error/warning by logging its pointer or using it when we configure repositories
	_ = db

	// 3. Initialize Echo Instance
	e := echo.New()

	// Register Custom Validator
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

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
