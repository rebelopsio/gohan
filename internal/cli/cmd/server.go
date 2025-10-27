package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rebelopsio/gohan/internal/container"
	httpinfra "github.com/rebelopsio/gohan/internal/infrastructure/http"
	"github.com/rebelopsio/gohan/internal/infrastructure/http/handlers"
	"github.com/spf13/cobra"
)

var (
	serverHost string
	serverPort int
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the API server",
	Long: `Start the Gohan HTTP API server. This allows remote installation
management and provides a REST API for client applications.

The server supports:
  - Installation management
  - Progress monitoring
  - Session persistence
  - Multiple concurrent installations

Configuration can be provided via environment variables or command-line flags.`,
	RunE: runServer,
}

func init() {
	serverCmd.Flags().StringVar(&serverHost, "host", "0.0.0.0", "Server host")
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Server port")
}

func runServer(cmd *cobra.Command, args []string) error {
	// Initialize dependency container
	c, err := container.New()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer c.Close()

	// Override with command-line flags if provided
	if cmd.Flags().Changed("host") {
		c.Config.API.Host = serverHost
	}
	if cmd.Flags().Changed("port") {
		c.Config.API.Port = serverPort
	}

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
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown...", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
			return fmt.Errorf("forced shutdown: %w", err)
		}

		log.Println("Server stopped gracefully")
	}

	return nil
}
