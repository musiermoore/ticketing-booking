package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (c *BookingController) GetList(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1" // default page
	}
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		http.Error(w, "invalid page parameter", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by middleware)
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

	// Call service to get list of bookings
	bookings, err := c.bookingService.GetList(r.Context(), userID, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
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

func (c *BookingController) RemoveBooking(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	bookingID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "Invalid booking id", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == "" {
		http.Error(w, "missing user ID", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		http.Error(w, "invalid user", http.StatusInternalServerError)
		return
	}

	errRemoving := c.bookingService.RemoveBooking(r.Context(), userID, bookingID)
	if errRemoving != nil {
		http.Error(w, errRemoving.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("The booking was cancelled.")
}
