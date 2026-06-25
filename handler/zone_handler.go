package handler

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/dto"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

// ZoneHandler maps HTTP endpoints to ZoneService actions.
type ZoneHandler struct {
	zoneService service.ZoneService
}

// NewZoneHandler creates a new instance of ZoneHandler.
func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

// Create handles creating a new parking zone.
// POST /api/v1/zones
func (h *ZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request format", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Validation errors", err.Error())
	}

	resp, err := h.zoneService.CreateZone(req)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to create parking zone", err.Error())
	}

	return utils.SendSuccess(c, http.StatusCreated, "Parking zone created successfully", resp)
}

// GetAll handles retrieving all parking zones.
// GET /api/v1/zones
func (h *ZoneHandler) GetAll(c echo.Context) error {
	resp, err := h.zoneService.GetAllZones()
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve parking zones", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Parking zones retrieved successfully", resp)
}

// GetByID handles retrieving details of a single parking zone.
// GET /api/v1/zones/:id
func (h *ZoneHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid zone ID parameter", err.Error())
	}

	resp, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrZoneNotFound) {
			return utils.SendError(c, http.StatusNotFound, "Parking zone not found", err.Error())
		}
		return utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve parking zone", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Parking zone retrieved successfully", resp)
}
