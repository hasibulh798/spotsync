package service

import (
	"errors"

	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"

	"gorm.io/gorm"
)

var (
	ErrForbidden           = errors.New("forbidden: insufficient permission to access this resource")
	ErrReservationNotFound = errors.New("reservation not found")
)

// ReservationService defines business logic for parking reservations.
type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.ReservationResponse, error)
	CancelReservation(reservationID uint, requestingUserID uint, requestingRole string) error
	GetAllReservations() ([]dto.ReservationResponse, error)
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ZoneRepository
}

// NewReservationService creates a new instance of ReservationService.
func NewReservationService(reservationRepo repository.ReservationRepository, zoneRepo repository.ZoneRepository) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

func (s *reservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	// 1. Verify that the zone exists first
	_, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	// 2. Call the locked transaction builder in repository
	res, err := s.reservationRepo.CreateReservationWithLock(req.ZoneID, userID, req.LicensePlate)
	if err != nil {
		return nil, err // Returns repository.ErrZoneFull or database errors
	}

	// 3. Map to Create Reservation Response (includes user_id, zone_id, status, license_plate)
	resp := mapToReservationResponse(res, false, false)
	return &resp, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var resp []dto.ReservationResponse
	for _, res := range reservations {
		resp = append(resp, mapToReservationResponse(&res, true, false))
	}
	return resp, nil
}

func (s *reservationService) CancelReservation(reservationID uint, requestingUserID uint, requestingRole string) error {
	// 1. Fetch reservation
	res, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReservationNotFound
		}
		return err
	}

	// 2. Authorization guard: drivers can only cancel their own reservations; admins can cancel any
	if requestingRole != "admin" && res.UserID != requestingUserID {
		return ErrForbidden
	}

	// 3. Perform status update to cancelled
	return s.reservationRepo.UpdateStatus(reservationID, "cancelled")
}

func (s *reservationService) GetAllReservations() ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var resp []dto.ReservationResponse
	for _, res := range reservations {
		resp = append(resp, mapToReservationResponse(&res, true, true))
	}
	return resp, nil
}

// mapToReservationResponse converts models.Reservation to DTO format with selective nested mapping.
func mapToReservationResponse(res *models.Reservation, includeZone bool, includeUser bool) dto.ReservationResponse {
	var zoneInfo *dto.ReservationZoneInfo
	if includeZone && res.Zone.ID != 0 {
		zoneInfo = &dto.ReservationZoneInfo{
			ID:   res.Zone.ID,
			Name: res.Zone.Name,
			Type: res.Zone.Type,
		}
	}

	var userInfo *dto.ReservationUserInfo
	if includeUser && res.User.ID != 0 {
		userInfo = &dto.ReservationUserInfo{
			ID:    res.User.ID,
			Name:  res.User.Name,
			Email: res.User.Email,
		}
	}

	var uID, zID uint
	if !includeUser {
		uID = res.UserID
	}
	if !includeZone {
		zID = res.ZoneID
	}

	return dto.ReservationResponse{
		ID:           res.ID,
		UserID:       uID,
		ZoneID:       zID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		Zone:         zoneInfo,
		User:         userInfo,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}
}
