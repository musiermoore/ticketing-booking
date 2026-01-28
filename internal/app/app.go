package app

import (
	"log"

	"github.com/musiermoore/ticketing-booking/internal/config"
	"github.com/musiermoore/ticketing-booking/internal/db"
	"github.com/musiermoore/ticketing-booking/internal/http"
)

func Run() {
	cfg := config.Load()

	database := db.Connect(cfg)
	db.Migrate(database)

	server := http.NewServer(cfg, database)

	log.Printf("Booking service running on :%s\n", cfg.AppPort)
	server.Start()
}
