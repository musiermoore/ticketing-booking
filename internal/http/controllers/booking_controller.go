package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/musiermoore/ticketing-booking/internal/http/middleware"
	"github.com/musiermoore/ticketing-booking/internal/service"
)

type BookingController struct {
	bookingService *service.BookingService
}

// Constructor
func NewBookingController(bs *service.BookingService) *BookingController {
	return &BookingController{bookingService: bs}
}

func (c *BookingController) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EventID json.Number `json:"event_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	eventID, err := input.EventID.Int64()
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == "" {
		http.Error(w, "missing user ID", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		http.Error(w, "invalid user ID type", http.StatusInternalServerError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	booking, err := c.bookingService.CreateBooking(r.Context(), userID, eventID, authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}
