package dto

import "time"

// CreateReservationRequest holds reservation booking payload.
type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

// ReservationZoneInfo is nested zone info in reservation response.
type ReservationZoneInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ReservationUserInfo is nested user info in reservation response.
type ReservationUserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ReservationResponse represents booking details in API responses.
type ReservationResponse struct {
	ID           uint                 `json:"id"`
	UserID       uint                 `json:"user_id,omitempty"`
	ZoneID       uint                 `json:"zone_id,omitempty"`
	LicensePlate string               `json:"license_plate"`
	Status       string               `json:"status"`
	Zone         *ReservationZoneInfo `json:"zone,omitempty"`
	User         *ReservationUserInfo `json:"user,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty"`
}
