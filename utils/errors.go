package utils

import (
	"errors"
	"net/http"
)

// Sentinel errors representing business logic violations.
var (
	ErrEmailExists         = errors.New("email already registered")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrZoneNotFound        = errors.New("parking zone not found")
	ErrReservationNotFound = errors.New("reservation not found")
	ErrZoneFull            = errors.New("parking zone is full")
	ErrForbidden           = errors.New("forbidden: insufficient permission to access this resource")
)

// MapErrorToHTTPStatus maps custom business errors to HTTP status codes and user-friendly messages.
func MapErrorToHTTPStatus(err error) (int, string) {
	if errors.Is(err, ErrEmailExists) {
		return http.StatusBadRequest, "User registration failed"
	}
	if errors.Is(err, ErrInvalidCredentials) {
		return http.StatusUnauthorized, "Login failed"
	}
	if errors.Is(err, ErrUserNotFound) {
		return http.StatusNotFound, "User not found"
	}
	if errors.Is(err, ErrZoneNotFound) {
		return http.StatusNotFound, "Parking zone not found"
	}
	if errors.Is(err, ErrReservationNotFound) {
		return http.StatusNotFound, "Reservation not found"
	}
	if errors.Is(err, ErrZoneFull) {
		return http.StatusConflict, "Reservation failed"
	}
	if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden, "Access denied"
	}
	return http.StatusInternalServerError, "Internal server error occurred"
}
