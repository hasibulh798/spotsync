package repository

import (
	"spotsync/models"

	"gorm.io/gorm"
)

// ZoneRepository defines the contract for ParkingZone database operations.
type ZoneRepository interface {
	Create(zone *models.ParkingZone) error
	FindAll() ([]models.ParkingZone, error)
	FindByID(id uint) (*models.ParkingZone, error)
	GetActiveReservationCounts() (map[uint]int, error)
	GetActiveReservationCount(zoneID uint) (int, error)
}

type zoneRepository struct {
	db *gorm.DB
}

// NewZoneRepository creates a GORM-based implementation of ZoneRepository.
func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepository) FindAll() ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	err := r.db.Find(&zones).Error
	return zones, err
}

func (r *zoneRepository) FindByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	err := r.db.First(&zone, id).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

// GetActiveReservationCounts returns a map of zone ID to active reservation count.
// It counts reservations where status is 'active' grouped by zone_id.
func (r *zoneRepository) GetActiveReservationCounts() (map[uint]int, error) {
	var results []struct {
		ZoneID uint
		Count  int
	}
	err := r.db.Model(&models.Reservation{}).
		Select("zone_id, count(*) as count").
		Where("status = ?", "active").
		Group("zone_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[uint]int)
	for _, res := range results {
		counts[res.ZoneID] = res.Count
	}
	return counts, nil
}

// GetActiveReservationCount returns the count of active reservations for a single zone.
func (r *zoneRepository) GetActiveReservationCount(zoneID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
