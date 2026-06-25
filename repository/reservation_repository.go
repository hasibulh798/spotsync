package repository

import (
	"spotsync/models"
	"spotsync/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ReservationRepository defines the contract for Reservation database operations.
type ReservationRepository interface {
	CreateReservationWithLock(zoneID uint, userID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	FindAll() ([]models.Reservation, error)
	UpdateStatus(id uint, status string) error
}

type reservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a GORM-based implementation of ReservationRepository.
func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// CreateReservationWithLock locks the ParkingZone row, checks its capacity, and creates a reservation atomically.
func (r *reservationRepository) CreateReservationWithLock(zoneID uint, userID uint, licensePlate string) (*models.Reservation, error) {
	var reservation *models.Reservation
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		// 1. Lock the ParkingZone row for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			return err
		}

		// 2. Count active reservations in this zone
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Capacity check
		if int(activeCount) >= zone.TotalCapacity {
			return utils.ErrZoneFull // Returns ErrZoneFull to rollback transaction
		}

		// 4. Create reservation
		newRes := &models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		if err := tx.Create(newRes).Error; err != nil {
			return err
		}
		reservation = newRes
		return nil // Commits transaction
	})

	if err != nil {
		return nil, err
	}
	return reservation, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error
}
