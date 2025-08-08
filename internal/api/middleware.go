package api

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// LoggingMiddleware creates a middleware that logs HTTP requests.
// This middleware captures and logs information about each HTTP request including:
// - HTTP method (GET, POST, etc.)
// - Request URI (the endpoint being accessed)
// - Remote address (client IP address)
// - Request duration (how long the request took to process)
//
// The middleware logs requests in the format: METHOD URI REMOTE_ADDR DURATION
//
// Returns:
//   - echo.MiddlewareFunc: middleware function that can be used with Echo
//
// Example log output:
//
//	2025/08/07 12:13:11 GET /api/calculate?qty=1201 [::1]:33284 318.867µs
func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Record the start time of the request
			start := time.Now()

			// Call the next handler in the middleware chain
			err := next(c)

			// Calculate the duration of the request
			duration := time.Since(start)

			// Log the request details including method, URI, remote address, and duration
			log.Printf(
				"%s %s %s %v",
				c.Request().Method,     // HTTP method (GET, POST, etc.)
				c.Request().RequestURI, // Full request URI including query parameters
				c.Request().RemoteAddr, // Client's IP address
				duration,               // Request duration
			)

			// Return any error from the next handler
			return err
		}
	}
}

// CORSMiddleware creates a middleware that adds CORS (Cross-Origin Resource Sharing) headers.
// This middleware allows web applications from different origins to access the API.
// It's essential for web interfaces that need to make requests to the API from different domains.
//
// CORS Headers Added:
//   - Access-Control-Allow-Origin: "*" (allows all origins)
//   - Access-Control-Allow-Methods: "GET, POST, OPTIONS" (allowed HTTP methods)
//   - Access-Control-Allow-Headers: "Content-Type" (allowed headers)
//
// Special Handling:
//   - OPTIONS requests are handled immediately with a 200 status (preflight requests)
//
// Returns:
//   - echo.MiddlewareFunc: middleware function that can be used with Echo
//
// Note: In production, you might want to restrict Access-Control-Allow-Origin
// to specific domains for security reasons.
func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Add CORS headers to allow cross-origin requests
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")

			// Handle preflight OPTIONS requests
			// These are sent by browsers before making actual requests to check CORS permissions
			if c.Request().Method == "OPTIONS" {
				// Return immediately with 200 status for preflight requests
				return c.NoContent(http.StatusOK)
			}

			// Continue to the next handler for non-OPTIONS requests
			return next(c)
		}
	}
}
