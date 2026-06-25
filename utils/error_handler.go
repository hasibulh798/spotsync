package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomHTTPErrorHandler acts as the centralized error handler for the Echo framework.
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "Internal server error occurred"
	var details interface{} = "An unexpected error occurred on the server"

	// 1. Check if it is a custom business sentinel error
	if status, msg := MapErrorToHTTPStatus(err); status != http.StatusInternalServerError {
		code = status
		message = msg
		details = err.Error()
	} else if he, ok := err.(*echo.HTTPError); ok {
		// 2. Handle Echo's own HTTP errors (like 404 Not Found, 405 Method Not Allowed)
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
		details = he.Message
	} else if ve, ok := err.(validator.ValidationErrors); ok {
		// 3. Formulate user-friendly messaging for validation errors
		code = http.StatusBadRequest
		message = "Validation failed"

		valErrors := make(map[string]string)
		for _, fieldErr := range ve {
			valErrors[fieldErr.Field()] = fmt.Sprintf("failed validation on rule: %s", fieldErr.ActualTag())
		}
		details = valErrors
	} else {
		// 4. Handle internal system / GORM / raw DB errors safely without leaking details
		c.Logger().Errorf("Internal Server Error: %v", err)
		details = "An unexpected error occurred"
	}

	// Send standard error structure
	_ = c.JSON(code, APIErrorResponse{
		Success: false,
		Message: message,
		Errors:  details,
	})
}
