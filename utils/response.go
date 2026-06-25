package utils

import (
	"github.com/labstack/echo/v4"
)

// APIResponse represents the standard success response format.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// APIErrorResponse represents the standard error response format.
type APIErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SendSuccess sends a JSON success response with appropriate status code.
func SendSuccess(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends a JSON error response with appropriate status code.
func SendError(c echo.Context, statusCode int, message string, errors interface{}) error {
	return c.JSON(statusCode, APIErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
