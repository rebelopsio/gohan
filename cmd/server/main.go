package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rebelopsio/gohan/internal/container"
	httpinfra "github.com/rebelopsio/gohan/internal/infrastructure/http"
	"github.com/rebelopsio/gohan/internal/infrastructure/http/handlers"
)

func main() {
	// Initialize dependency container
	c, err := container.New()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer c.Close()

	log.Println("Starting Gohan Installation Server...")
	log.Printf("Configuration: Host=%s Port=%d", c.Config.API.Host, c.Config.API.Port)

	// Create handlers using pre-wired use cases from container
	installationHandler := handlers.NewInstallationHandler(
		c.StartInstallationUseCase,
		c.ExecuteInstallationUseCase,
		c.GetStatusUseCase,
		c.ListInstallationsUseCase,
		c.CancelInstallationUseCase,
	)

	// Create HTTP server
	serverConfig := httpinfra.Config{
		Host:         c.Config.API.Host,
		Port:         c.Config.API.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server := httpinfra.NewServer(serverConfig, installationHandler, false)

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("HTTP server listening on %s:%d", c.Config.API.Host, c.Config.API.Port)
		serverErrors <- server.Start()
	}()

	// Wait for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown...", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
			log.Fatal("Forcing shutdown")
		}

		log.Println("Server stopped gracefully")
	}
}

