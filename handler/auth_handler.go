package handler

import (
	"net/http"

	"spotsync/dto"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

// AuthHandler maps HTTP endpoints to AuthService actions.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration.
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		return err
	}

	return utils.SendSuccess(c, http.StatusCreated, "User registered successfully", resp)
}

// Login handles user login and credentials verification.
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		return err
	}

	return utils.SendSuccess(c, http.StatusOK, "Login successful", resp)
}
