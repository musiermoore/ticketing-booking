package http

import (
	"database/sql"
	"net/http"

	"github.com/musiermoore/ticketing-booking/internal/config"
	"github.com/musiermoore/ticketing-booking/internal/http/middleware"
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

	protected.HandleFunc("/book", postOnly(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Booking"))
	}))

	// Apply JWT middleware ONCE
	mux.Handle("/", middleware.JWT(cfg)(protected))

	return mux
}

func postOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}
