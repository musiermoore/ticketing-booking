package http

import (
	"database/sql"
	"net/http"

	"github.com/musiermoore/ticketing-booking/internal/config"
)

type Server struct {
	cfg *config.Config
	db  *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{cfg: cfg, db: db}
}

func (s *Server) Start() {
	router := NewRouter(s.cfg, s.db)

	http.ListenAndServe(":"+s.cfg.AppPort, router)
}
