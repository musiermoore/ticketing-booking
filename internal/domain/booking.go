package domain

import "time"

type Booking struct {
	ID        int64  // UUID
	UserID    string // comes from auth microservice
	EventID   string // the booked event
	CreatedAt time.Time
}

// Repository interface (minimal)
type BookingRepository interface {
	Create(b Booking) (*Booking, error)
	GetByID(id int64) (*Booking, error)
}
