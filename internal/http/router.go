package http

import (
	"database/sql"
	"net/http"

	"github.com/musiermoore/ticketing-booking/internal/clients"
	"github.com/musiermoore/ticketing-booking/internal/config"
	"github.com/musiermoore/ticketing-booking/internal/http/controllers"
	"github.com/musiermoore/ticketing-booking/internal/http/middleware"
	"github.com/musiermoore/ticketing-booking/internal/repository"
	"github.com/musiermoore/ticketing-booking/internal/service"
)

func NewRouter(cfg *config.Config, db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Protected routes
	protected := http.NewServeMux()

	protected.HandleFunc("/auth/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Authorized"))
	})

	bookingRepo := repository.NewPostgresBookingRepository(db)
	eventsClient := clients.NewEventsClient(cfg.APIBaseURL)
	bookingSvc := service.NewBookingService(bookingRepo, eventsClient)
	bookingCtrl := controllers.NewBookingController(bookingSvc)

	protected.HandleFunc("/tickets/book", postOnly(bookingCtrl.CreateBooking))
	protected.HandleFunc("/tickets/{id}/unbook", deleteOnly(bookingCtrl.RemoveBooking))

	// Apply JWT middleware ONCE
	mux.Handle("/", middleware.JWT(cfg)(protected))

	return mux
}

func checkMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func postOnly(h http.HandlerFunc) http.HandlerFunc {
	return checkMethod(h, http.MethodPost)
}

func deleteOnly(h http.HandlerFunc) http.HandlerFunc {
	return checkMethod(h, http.MethodDelete)
}
