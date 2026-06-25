package models

import "time"

// ParkingZone represents the parking_zones table in PostgreSQL database.
type ParkingZone struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Type          string    `gorm:"type:varchar(50);not null" json:"type"` // general, ev_charging, covered
	TotalCapacity int       `gorm:"not null" json:"total_capacity"`
	PricePerHour  float64   `gorm:"type:numeric(10,2);not null" json:"price_per_hour"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
