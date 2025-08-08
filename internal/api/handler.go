package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"package-optimizer/internal/domain"

	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests for the package optimizer API.
// It provides endpoints for package optimization calculations and web UI serving.
type Handler struct {
	// optimizer is the core optimization engine that calculates optimal package combinations
	optimizer *domain.Optimizer
	// packageSizes stores the available package sizes for the API
	packageSizes []int
}

// NewHandler creates a new handler with the given optimizer.
// This function initializes the handler with the domain optimizer for package calculations.
//
// Args:
//   - optimizer: the domain optimizer instance for package calculations
//   - packageSizes: the available package sizes for the API
//
// Returns:
//   - *Handler: configured handler instance
func NewHandler(optimizer *domain.Optimizer, packageSizes []int) *Handler {
	return &Handler{
		optimizer:    optimizer,
		packageSizes: packageSizes,
	}
}

// CalculateHandler handles the /calculate endpoint for package optimization.
// This is the main API endpoint that accepts a quantity parameter and returns
// the optimal package combination that minimizes over-delivery.
//
// Query Parameters:
//   - qty: the requested quantity (required, must be a positive integer)
//
// Returns:
//   - JSON response with optimization result or error
//   - HTTP 400 if quantity is missing or invalid
//   - HTTP 200 with optimization result on success
//
// Example:
//
//	GET /api/calculate?qty=1201
//	Response: {"requested":1201,"total_delivered":1250,"over_delivery":49,"packages":{"1000":1,"250":1}}
func (h *Handler) CalculateHandler(c echo.Context) error {
	// Extract quantity parameter from query string
	qtyStr := c.QueryParam("qty")
	if qtyStr == "" {
		// Return error if quantity parameter is missing
		return echo.NewHTTPError(http.StatusBadRequest, "missing 'qty' parameter")
	}

	// Parse quantity string to integer
	quantity, err := strconv.Atoi(qtyStr)
	if err != nil {
		// Return error if quantity is not a valid integer
		return echo.NewHTTPError(http.StatusBadRequest, "invalid 'qty' parameter: must be an integer")
	}

	// Use the optimizer to calculate the optimal package combination
	result, err := h.optimizer.Optimize(quantity)
	if err != nil {
		// Log the optimization error for debugging
		log.Printf("Optimization error: %v", err)
		// Return error response to client
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("optimization error: %v", err))
	}

	// Return the optimization result as JSON response
	return c.JSON(http.StatusOK, result)
}

// PackageSizesHandler handles the /package-sizes endpoint.
// This endpoint returns the available package sizes that can be used for optimization.
//
// Returns:
//   - JSON response with available package sizes
//   - HTTP 200 with package sizes array
//
// Example:
//
//	GET /api/package-sizes
//	Response: {"package_sizes":[250,500,1000,2000]}
func (h *Handler) PackageSizesHandler(c echo.Context) error {
	// Return the available package sizes as JSON response
	return c.JSON(http.StatusOK, map[string][]int{
		"package_sizes": h.packageSizes,
	})
}

// HealthHandler handles health check requests.
// This endpoint is used by load balancers and monitoring systems to check if the service is running.
//
// Returns:
//   - JSON response with service status
//   - HTTP 200 with {"status":"healthy"}
//
// Example:
//
//	GET /api/health
//	Response: {"status":"healthy"}
func (h *Handler) HealthHandler(c echo.Context) error {
	// Return a simple health status response
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// ServeWebUI serves the main web interface.
// This endpoint serves the HTML page that provides a user-friendly interface
// for testing the package optimization API.
//
// Returns:
//   - HTML content of the web interface
//   - HTTP 200 with the index.html file
//
// Example:
//
//	GET /
//	Response: HTML content of the web interface
func (h *Handler) ServeWebUI(c echo.Context) error {
	// Serve the main HTML file for the web interface
	return c.File("web/static/index.html")
}

// ServeCSS serves CSS stylesheets for the web interface.
// This endpoint serves the CSS file that styles the web interface.
//
// Returns:
//   - CSS content for styling the web interface
//   - HTTP 200 with the style.css file
//
// Example:
//
//	GET /style.css
//	Response: CSS content for styling
func (h *Handler) ServeCSS(c echo.Context) error {
	// Serve the CSS file for styling the web interface
	return c.File("web/static/style.css")
}

// ServeJS serves JavaScript files for the web interface.
// This endpoint serves the JavaScript file that provides interactive functionality
// for the web interface.
//
// Returns:
//   - JavaScript content for web interface functionality
//   - HTTP 200 with the script.js file
//
// Example:
//
//	GET /script.js
//	Response: JavaScript content for interactivity
func (h *Handler) ServeJS(c echo.Context) error {
	// Serve the JavaScript file for web interface functionality
	return c.File("web/static/script.js")
}
