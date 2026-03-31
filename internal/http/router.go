package http

import (
	"database/sql"
	"net/http"
	"net/url"
	"strings"

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

	protected.HandleFunc("/tickets", getOnly(bookingCtrl.GetList))
	protected.HandleFunc("/tickets/book", postOnly(bookingCtrl.CreateBooking))
	protected.HandleFunc("/tickets/{id}/unbook", deleteOnly(bookingCtrl.RemoveBooking))

	// Apply JWT middleware ONCE
	mux.Handle("/", middleware.JWT(cfg)(protected))

	return withCORS(cfg, mux)
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

func getOnly(h http.HandlerFunc) http.HandlerFunc {
	return checkMethod(h, http.MethodGet)
}

func postOnly(h http.HandlerFunc) http.HandlerFunc {
	return checkMethod(h, http.MethodPost)
}

func deleteOnly(h http.HandlerFunc) http.HandlerFunc {
	return checkMethod(h, http.MethodDelete)
}

func withCORS(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isAllowedOrigin(origin, cfg.UIBaseURL) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin, configured string) bool {
	if origin == "" {
		return false
	}

	normalizedOrigin := strings.TrimRight(origin, "/")
	normalizedConfigured := strings.TrimRight(configured, "/")

	if normalizedOrigin == normalizedConfigured {
		return true
	}

	originURL, err := url.Parse(normalizedOrigin)
	if err != nil {
		return false
	}

	configuredURL, err := url.Parse(normalizedConfigured)
	if err != nil {
		return false
	}

	if originURL.Scheme != configuredURL.Scheme || originURL.Port() != configuredURL.Port() {
		return false
	}

	if configuredURL.Hostname() == "ui" {
		return originURL.Hostname() == "localhost" || originURL.Hostname() == "127.0.0.1"
	}

	return false
}
