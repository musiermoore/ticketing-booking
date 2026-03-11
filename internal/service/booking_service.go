package service

import (
	"context"
	"errors"

	"github.com/musiermoore/ticketing-booking/internal/domain"
)

type EventValidator interface {
	ValidateEvent(ctx context.Context, eventID int64, authHeader string) error
}

type BookingService struct {
	repo   domain.BookingRepository
	events EventValidator
}

func NewBookingService(repo domain.BookingRepository, events EventValidator) *BookingService {
	return &BookingService{repo: repo, events: events}
}

// CreateBooking just needs userID and eventID
func (s *BookingService) CreateBooking(ctx context.Context, userID, eventID int64, authHeader string) (*domain.Booking, error) {
	if userID == 0 || eventID == 0 {
		return nil, errors.New("userID and eventID are required")
	}

	if s.events != nil {
		if err := s.events.ValidateEvent(ctx, eventID, authHeader); err != nil {
			return nil, err
		}
	}

	b := domain.Booking{
		UserID:  userID,
		EventID: eventID,
	}

	return s.repo.Create(b)
}

// CreateBooking just needs userID and eventID
func (s *BookingService) RemoveBooking(ctx context.Context, userID, bookingID int64) error {
	booking, err := s.repo.GetByID(bookingID)

	if err != nil || booking.UserID != userID || booking.ID != bookingID {
		return errors.New("Booking not found.")
	}

	err = s.repo.Remove(userID, bookingID)

	if err != nil {
		return errors.New("Booking is not deleted.")
	}

	return nil
}
