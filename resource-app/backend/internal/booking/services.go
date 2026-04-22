package booking

import (
	"slices"
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

func (s *Service) GetBookings(filter BookingFilter) ([]Booking, error) {
	var userID *string
	var resourceIDs []string

	// Apply generic resource filter first so it works even when scope is omitted.
	if filter.ResourceID != "" {
		resourceIDs = []string{filter.ResourceID}
	}

	switch filter.Scope {
	// "me" scope: only return bookings made by the current user.
	case BookingScopeMe:
		userID = &filter.CurrentUserID

	// "approvable" scope: return bookings for resources the user has APPROVE permission on.
	case BookingScopeApprovable:
		approvableResourceIDs, err := s.permissionSvc.GetApprovableResourceIDs(filter.CurrentUserID)
		if err != nil {
			return nil, err
		}

		if len(approvableResourceIDs) == 0 {
			return []Booking{}, nil
		}

		if filter.ResourceID != "" {
			// If a specific resource is requested, it must be approvable by this user.
			if !slices.Contains(approvableResourceIDs, filter.ResourceID) {
				return nil, ErrBookingViewPermissionDenied
			}
		} else {
			// No specific resource requested: return bookings from all approvable resources.
			resourceIDs = approvableResourceIDs
		}
	}

	return s.repo.GetBookings(userID, filter.Statuses, resourceIDs)
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