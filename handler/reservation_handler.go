package handler

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/dto"
	"spotsync/repository"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

// ReservationHandler maps HTTP endpoints to ReservationService actions.
type ReservationHandler struct {
	reservationService service.ReservationService
}

// NewReservationHandler creates a new instance of ReservationHandler.
func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

// Create handles booking a parking spot.
// POST /api/v1/reservations
func (h *ReservationHandler) Create(c echo.Context) error {
	// Retrieve authenticated user ID from context
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "User claims not found in context")
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request format", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Validation errors", err.Error())
	}

	resp, err := h.reservationService.CreateReservation(userID, req)
	if err != nil {
		if errors.Is(err, repository.ErrZoneFull) {
			return utils.SendError(c, http.StatusConflict, "Reservation failed", "Parking zone is full")
		}
		if errors.Is(err, service.ErrZoneNotFound) {
			return utils.SendError(c, http.StatusNotFound, "Reservation failed", err.Error())
		}
		return utils.SendError(c, http.StatusInternalServerError, "Internal server error occurred", err.Error())
	}

	return utils.SendSuccess(c, http.StatusCreated, "Reservation confirmed successfully", resp)
}

// GetMy handles retrieving reservations for the authenticated driver.
// GET /api/v1/reservations/my-reservations
func (h *ReservationHandler) GetMy(c echo.Context) error {
	// Retrieve authenticated user ID from context
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "User claims not found in context")
	}

	resp, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve reservations", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "My reservations retrieved successfully", resp)
}

// Cancel handles cancelling an active reservation.
// DELETE /api/v1/reservations/:id
func (h *ReservationHandler) Cancel(c echo.Context) error {
	// Retrieve user claims from context
	userIDVal := c.Get("user_id")
	userID, okUserID := userIDVal.(uint)
	roleVal := c.Get("role")
	role, okRole := roleVal.(string)

	if !okUserID || !okRole || userID == 0 || role == "" {
		return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "User claims not found in context")
	}

	idStr := c.Param("id")
	reservationID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid reservation ID parameter", err.Error())
	}

	err = h.reservationService.CancelReservation(uint(reservationID), userID, role)
	if err != nil {
		if errors.Is(err, service.ErrReservationNotFound) {
			return utils.SendError(c, http.StatusNotFound, "Cancellation failed", err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return utils.SendError(c, http.StatusForbidden, "Cancellation failed", err.Error())
		}
		return utils.SendError(c, http.StatusInternalServerError, "Internal server error occurred", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

// GetAll handles retrieving all reservations in the system.
// GET /api/v1/reservations
func (h *ReservationHandler) GetAll(c echo.Context) error {
	resp, err := h.reservationService.GetAllReservations()
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve reservations", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "All reservations retrieved successfully", resp)
}
