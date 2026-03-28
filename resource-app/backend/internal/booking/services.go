package booking

import (
	"time"

	perm "resource-app/internal/permission"
	usr "resource-app/internal/user"

	"github.com/google/uuid"
)

type Service struct {
	repo            Repository
	permissionSvc   *perm.Service
}

func NewService(repo Repository, permissionSvc *perm.Service) *Service {
	return &Service{
		repo:            repo,
		permissionSvc:   permissionSvc,
	}
}

func (s *Service) GetBookings() ([]Booking, error) {
	return s.repo.GetBookings()
}

func (s *Service) CreateBooking(booking *Booking, userID string, userRole usr.Role) error {
	// For non-admin users, enforce REQUEST permission check
	if userRole != usr.RoleAdmin {
		hasPermission, err := s.permissionSvc.HasRequestPermission(userID, booking.ResourceID)
		if err != nil {
			return err
		}
		if !hasPermission {
			return ErrBookingPermissionDenied
		}
	}

	booking.ID = uuid.New().String()
	booking.UserID = userID
	booking.CreatedAt = time.Now()

	if userRole == usr.RoleAdmin {
		booking.Status = StatusConfirmed
	} else {
		booking.Status = StatusPending
	}

	return s.repo.CreateBooking(booking)
}

func (s *Service) UpdateBookingStatus(id string, status BookingStatus, rejectionReason *string) (*Booking, error) {
	return s.repo.UpdateBookingStatus(id, status, rejectionReason)
}

func (s *Service) RescheduleBooking(id string, newStart, newEnd time.Time) (*Booking, error) {
	return s.repo.RescheduleBooking(id, newStart, newEnd)
}

func (s *Service) CancelBooking(id string) error {
	return s.repo.CancelBooking(id)
}

func (s *Service) GetUtilizationStats() ([]ResourceUsageStats, error) {
	return s.repo.GetUtilizationStats()
}