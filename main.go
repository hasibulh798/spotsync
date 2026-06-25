package main

import (
	"log"
	"net/http"

	"spotsync/config"
	"spotsync/handler"
	"spotsync/models"
	"spotsync/repository"
	"spotsync/service"
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

	// 3. Manual Dependency Injection
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	// 4. Initialize Echo Instance
	e := echo.New()

	// Register Custom Validator
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	// 5. Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 6. Health Check Route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "SpotSync API is healthy and running",
		})
	})

	// 7. API Routes
	api := e.Group("/api/v1")

	// Auth Group
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// 8. Start Server
	log.Printf("Starting SpotSync server on port %s...", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
