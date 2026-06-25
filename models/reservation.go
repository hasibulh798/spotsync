package models

import "time"

// Reservation represents the reservations table in PostgreSQL database.
type Reservation struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint        `gorm:"not null;index" json:"user_id"`
	User         User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user"`
	ZoneID       uint        `gorm:"not null;index" json:"zone_id"`
	Zone         ParkingZone `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"zone"`
	LicensePlate string      `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       string      `gorm:"type:varchar(20);default:active;not null" json:"status"` // active, completed, cancelled
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
