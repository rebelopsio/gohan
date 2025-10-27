package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rebelopsio/gohan/internal/infrastructure/http/handlers"
	"github.com/rebelopsio/gohan/internal/infrastructure/http/middleware"
)

// Server represents the HTTP server
type Server struct {
	router *chi.Mux
	server *http.Server
}

// Config holds server configuration
type Config struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewServer creates a new HTTP server with configured routes and middleware
func NewServer(
	config Config,
	installationHandler *handlers.InstallationHandler,
	enableTracing bool,
) *Server {
	r := chi.NewRouter()

	// Apply middleware
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)

	// Add tracing middleware if enabled
	if enableTracing {
		r.Use(middleware.Tracing("gohan-api"))
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Installation routes
		r.Route("/installation", func(r chi.Router) {
			r.Get("/", installationHandler.ListInstallations)
			r.Post("/start", installationHandler.StartInstallation)
			r.Post("/{sessionID}/execute", installationHandler.ExecuteInstallation)
			r.Get("/{sessionID}/status", installationHandler.GetStatus)
			r.Post("/{sessionID}/cancel", installationHandler.CancelInstallation)
		})
	})

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	return &Server{
		router: r,
		server: srv,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting HTTP server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}

// Router returns the chi router (useful for testing)
func (s *Server) Router() *chi.Mux {
	return s.router
}
