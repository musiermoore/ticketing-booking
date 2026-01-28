package domain

import "time"

type Booking struct {
	ID        int64 // UUID
	UserID    int64 // comes from auth microservice
	EventID   int64 // the booked event
	CreatedAt time.Time
}

// Repository interface (minimal)
type BookingRepository interface {
	Create(b Booking) (*Booking, error)
	GetByID(id int64) (*Booking, error)
}
