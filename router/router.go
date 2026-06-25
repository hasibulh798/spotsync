package router

import (
	"spotsync/handler"
	"spotsync/middleware"

	"github.com/labstack/echo/v4"
)

// RouterConfig aggregates all handler dependencies to set up routes.
type RouterConfig struct {
	Echo        *echo.Echo
	AuthHandler *handler.AuthHandler
	JWTSecret   string
}

// SetupRoutes registers public and protected endpoints.
func SetupRoutes(cfg RouterConfig) {
	// Base API Group
	api := cfg.Echo.Group("/api/v1")

	// Public Auth Group
	authGroup := api.Group("/auth")
	authGroup.POST("/register", cfg.AuthHandler.Register)
	authGroup.POST("/login", cfg.AuthHandler.Login)

	// Protected Group (Requires login)
	// We initialize this now. In later steps, we will mount zone and reservation routes here.
	protected := api.Group("", middleware.JWTMiddleware(cfg.JWTSecret))
	
	// Just a dummy ping to test protected routing group
	protected.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"success": true,
			"message": "Authenticated ping successful",
			"user_id": c.Get("user_id"),
			"role":    c.Get("role"),
		})
	})
}
