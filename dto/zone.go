package dto

import "time"

// CreateZoneRequest holds the zone creation payload.
type CreateZoneRequest struct {
	Name          string  `json:"name" validate:"required"`
	Type          string  `json:"type" validate:"required,oneof=general ev_charging covered"`
	TotalCapacity int     `json:"total_capacity" validate:"required,gt=0"`
	PricePerHour  float64 `json:"price_per_hour" validate:"required,gt=0"`
}

// ZoneResponse represents zone info in API responses.
type ZoneResponse struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Type           string     `json:"type"`
	TotalCapacity  int        `json:"total_capacity"`
	AvailableSpots *int       `json:"available_spots,omitempty"`
	PricePerHour   float64    `json:"price_per_hour"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}
