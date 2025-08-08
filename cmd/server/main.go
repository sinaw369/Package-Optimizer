package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"package-optimizer/internal/api"
	"package-optimizer/internal/config"
	"package-optimizer/internal/domain"
)

// main is the entry point of the package optimizer application.
// It sets up the server, configures routes, and starts the HTTP server with graceful shutdown.
func main() {
	// Load application configuration from environment variables
	// This includes port number and package sizes
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create the core optimizer with the configured package sizes
	// The optimizer will be used by the API handlers to calculate optimal package combinations
	optimizer := domain.NewOptimizer(cfg.PackageSizes)

	// Create the HTTP handler with the optimizer and package sizes
	// The handler provides the API endpoints for package optimization
	handler := api.NewHandler(optimizer, cfg.PackageSizes)

	// Create a new Echo instance for the HTTP server
	// Echo is a high-performance web framework for Go
	e := echo.New()

	// Add middleware to the Echo instance
	// Middleware functions are executed in order for each request
	e.Use(api.LoggingMiddleware()) // Log all HTTP requests
	e.Use(api.CORSMiddleware())    // Enable CORS for web interface

	// Configure API routes under the /api prefix
	// These routes handle the core functionality of the package optimizer
	apiGroup := e.Group("/api")
	apiGroup.GET("/calculate", handler.CalculateHandler)     // Main optimization endpoint
	apiGroup.GET("/package-sizes", handler.PackageSizesHandler) // Package sizes endpoint
	apiGroup.GET("/health", handler.HealthHandler)           // Health check endpoint

	// Configure web UI routes
	// These routes serve the static files for the web interface
	e.GET("/", handler.ServeWebUI)        // Main web interface
	e.GET("/style.css", handler.ServeCSS) // CSS styles
	e.GET("/script.js", handler.ServeJS)  // JavaScript functionality

	// Legacy route for backward compatibility
	// This allows the old /calculate endpoint to still work
	e.GET("/calculate", handler.CalculateHandler)

	// Start the server in a goroutine to allow for graceful shutdown
	go func() {
		// Log server startup information
		log.Printf("Starting server on port %s", cfg.Port)
		log.Printf("Available package sizes: %v", cfg.PackageSizes)
		log.Printf("API endpoint: http://localhost:%s/api/calculate?qty=<quantity>", cfg.Port)
		log.Printf("Package sizes endpoint: http://localhost:%s/api/package-sizes", cfg.Port)
		log.Printf("Web UI: http://localhost:%s", cfg.Port)

		// Start the HTTP server
		// This will block until the server is stopped
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	// This allows the server to shut down cleanly when receiving SIGINT or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Log that shutdown is beginning
	log.Println("Shutting down server...")

	// Perform graceful shutdown with a timeout
	// This gives the server time to finish processing current requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the Echo server gracefully
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Log successful shutdown
	log.Println("Server exited")
}
