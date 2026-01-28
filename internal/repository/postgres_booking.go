package repository

import (
	"database/sql"

	"github.com/musiermoore/ticketing-booking/internal/domain"
)

type PostgresBookingRepository struct {
	db *sql.DB
}

func NewPostgresBookingRepository(db *sql.DB) *PostgresBookingRepository {
	return &PostgresBookingRepository{db: db}
}

func (r *PostgresBookingRepository) Create(b domain.Booking) (*domain.Booking, error) {
	var id int64
	err := r.db.QueryRow(
		"INSERT INTO bookings (user_id, event_id) VALUES ($1, $2) RETURNING id",
		b.UserID, b.EventID,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	b.ID = id
	return &b, nil
}

func (r *PostgresBookingRepository) GetByID(id string) (*domain.Booking, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, event_id, created_at FROM bookings WHERE id=$1",
		id,
	)

	var b domain.Booking
	if err := row.Scan(&b.ID, &b.UserID, &b.EventID, &b.CreatedAt); err != nil {
		return nil, err
	}
	return &b, nil
}
