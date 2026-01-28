package service

import (
	"errors"

	"github.com/musiermoore/ticketing-booking/internal/domain"
)

type BookingService struct {
	repo domain.BookingRepository
}

func NewBookingService(repo domain.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

// CreateBooking just needs userID and eventID
func (s *BookingService) CreateBooking(userID, eventID int64) (*domain.Booking, error) {
	if userID == 0 || eventID == 0 {
		return nil, errors.New("userID and eventID are required")
	}

	b := domain.Booking{
		UserID:  userID,
		EventID: eventID,
	}

	return s.repo.Create(b)
}
