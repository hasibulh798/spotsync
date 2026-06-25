package service

import (
	"errors"

	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"
	"spotsync/utils"

	"gorm.io/gorm"
)

// ZoneService defines the business operations for Parking Zones.
type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

// NewZoneService creates a new instance of ZoneService.
func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     &zone.UpdatedAt,
	}, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	activeCounts, err := s.zoneRepo.GetActiveReservationCounts()
	if err != nil {
		return nil, err
	}

	var resp []dto.ZoneResponse
	for _, zone := range zones {
		activeCount := activeCounts[zone.ID]
		available := zone.TotalCapacity - activeCount
		if available < 0 {
			available = 0
		}

		availCopy := available
		resp = append(resp, dto.ZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: &availCopy,
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
		})
	}

	return resp, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrZoneNotFound
		}
		return nil, err
	}

	activeCount, err := s.zoneRepo.GetActiveReservationCount(id)
	if err != nil {
		return nil, err
	}

	available := zone.TotalCapacity - activeCount
	if available < 0 {
		available = 0
	}

	availCopy := available
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: &availCopy,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}, nil
}
