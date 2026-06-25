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
	ZoneHandler *handler.ZoneHandler
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

	// Public Zone Routes
	api.GET("/zones", cfg.ZoneHandler.GetAll)
	api.GET("/zones/:id", cfg.ZoneHandler.GetByID)

	// Protected Group (Requires login)
	protected := api.Group("", middleware.JWTMiddleware(cfg.JWTSecret))
	
	// Admin Only Group (Requires login & admin role)
	adminOnly := api.Group("", middleware.JWTMiddleware(cfg.JWTSecret), middleware.RoleMiddleware("admin"))
	adminOnly.POST("/zones", cfg.ZoneHandler.Create)

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
