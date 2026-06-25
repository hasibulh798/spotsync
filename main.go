package main

import (
	"log"
	"net/http"

	"spotsync/config"
	"spotsync/handler"
	"spotsync/models"
	"spotsync/repository"
	"spotsync/router"
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
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	zoneService := service.NewZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo, zoneRepo)

	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)

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

	// 7. API Routes (Public & Protected Setup)
	router.SetupRoutes(router.RouterConfig{
		Echo:               e,
		AuthHandler:        authHandler,
		ZoneHandler:        zoneHandler,
		ReservationHandler: reservationHandler,
		JWTSecret:          cfg.JWTSecret,
	})

	// 8. Start Server
	log.Printf("Starting SpotSync server on port %s...", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
