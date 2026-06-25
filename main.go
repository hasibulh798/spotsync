package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	// Register Global Custom Error Handler
	e.HTTPErrorHandler = utils.CustomHTTPErrorHandler

	// 5. Middleware Setup
	// Request logging middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[Echo] time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
	}))
	e.Use(middleware.Recover())

	// CORS middleware configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(cfg.AllowedOrigins, ","),
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

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

	// 8. Start Server with Graceful Shutdown
	go func() {
		log.Printf("Starting SpotSync server on port %s...", cfg.Port)
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down SpotSync server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server Graceful Shutdown Failed: %v", err)
	}
	log.Println("Server exited gracefully")
}
