package http

import (
	"database/sql"
	stdhttp "net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/musiermoore/ticketing-booking/internal/clients"
	"github.com/musiermoore/ticketing-booking/internal/config"
	"github.com/musiermoore/ticketing-booking/internal/http/controllers"
	"github.com/musiermoore/ticketing-booking/internal/http/middleware"
	"github.com/musiermoore/ticketing-booking/internal/repository"
	"github.com/musiermoore/ticketing-booking/internal/service"
)

func NewRouter(cfg *config.Config, db *sql.DB) stdhttp.Handler {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), corsMiddleware(cfg))

	bookingRepo := repository.NewPostgresBookingRepository(db)
	eventsClient := clients.NewEventsClient(cfg.APIBaseURL)
	bookingSvc := service.NewBookingService(bookingRepo, eventsClient)
	bookingCtrl := controllers.NewBookingController(bookingSvc)

	router.GET("/health", func(c *gin.Context) {
		c.String(stdhttp.StatusOK, "OK")
	})

	protected := router.Group("/")
	protected.Use(wrapStdMiddleware(middleware.JWT(cfg)))

	protected.GET("/auth/check", func(c *gin.Context) {
		c.String(stdhttp.StatusOK, "Authorized")
	})
	protected.GET("/tickets", wrapHTTPHandler(bookingCtrl.GetList))
	protected.POST("/tickets/book", wrapHTTPHandler(bookingCtrl.CreateBooking))
	protected.DELETE("/tickets/:id/unbook", wrapHTTPHandler(bookingCtrl.RemoveBooking))

	return router
}

func wrapHTTPHandler(handler stdhttp.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request.Clone(c.Request.Context())
		for _, param := range c.Params {
			req.SetPathValue(param.Key, param.Value)
		}

		handler(c.Writer, req)
	}
}

func wrapStdMiddleware(mw func(stdhttp.Handler) stdhttp.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nextCalled bool

		handler := mw(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			nextCalled = true
			c.Request = r
			c.Next()
		}))

		handler.ServeHTTP(c.Writer, c.Request)

		if !nextCalled {
			c.Abort()
		}
	}
}

func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if isAllowedOrigin(origin, cfg.UIBaseURL) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		}

		if c.Request.Method == stdhttp.MethodOptions {
			c.Status(stdhttp.StatusNoContent)
			c.Abort()
			return
		}

		c.Next()
	}
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
