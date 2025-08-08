package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration loaded from environment variables.
// This structure contains all the configurable parameters for the package optimizer service.
type Config struct {
	// Port is the HTTP server port (e.g., "8080")
	Port string
	// PackageSizes is a slice of available package sizes for optimization
	// These are the fixed-size packages that can be used to fulfill orders
	PackageSizes []int
}

// Load loads configuration from environment variables.
// This function reads the PORT and PACKAGE_SIZES environment variables
// and returns a configured Config struct.
//
// Environment Variables:
//   - PORT: HTTP server port (default: "8080")
//   - PACKAGE_SIZES: Comma-separated list of package sizes (default: "250,500,1000,2000")
//
// Returns:
//   - *Config: configured application settings
//   - error: if package sizes are invalid or cannot be parsed
//
// Example:
//
//	export PORT=3000
//	export PACKAGE_SIZES="100,200,500,1000"
func Load() (*Config, error) {
	// Get port from environment variable with default value
	port := getEnv("PORT", "8080")

	// Get package sizes from environment variable with default value
	packageSizesStr := getEnv("PACKAGE_SIZES", "250,500,1000,2000")

	// Parse the package sizes string into a slice of integers
	packageSizes, err := parsePackageSizes(packageSizesStr)
	if err != nil {
		return nil, fmt.Errorf("invalid package sizes: %w", err)
	}

	// Return the configured application settings
	return &Config{
		Port:         port,
		PackageSizes: packageSizes,
	}, nil
}

// getEnv gets an environment variable with a default value.
// If the environment variable is not set or is empty, it returns the default value.
//
// Args:
//   - key: the environment variable name
//   - defaultValue: the default value to return if the environment variable is not set
//
// Returns:
//   - string: the environment variable value or the default value
//
// Example:
//
//	port := getEnv("PORT", "8080") // Returns "8080" if PORT is not set
func getEnv(key, defaultValue string) string {
	// Check if the environment variable is set
	if value := os.Getenv(key); value != "" {
		return value
	}
	// Return the default value if the environment variable is not set
	return defaultValue
}

// parsePackageSizes parses a comma-separated string of package sizes into a slice of integers.
// This function validates that all package sizes are positive integers.
//
// Args:
//   - sizesStr: comma-separated string of package sizes (e.g., "250,500,1000,2000")
//
// Returns:
//   - []int: slice of validated package sizes
//   - error: if the string is empty, contains invalid numbers, or has non-positive values
//
// Example:
//
//	sizes, err := parsePackageSizes("250,500,1000") // Returns []int{250, 500, 1000}, nil
func parsePackageSizes(sizesStr string) ([]int, error) {
	// Validate that the input string is not empty
	if sizesStr == "" {
		return nil, fmt.Errorf("package sizes cannot be empty")
	}

	// Split the comma-separated string into individual size strings
	sizes := strings.Split(sizesStr, ",")
	result := make([]int, 0, len(sizes))

	// Process each package size string
	for _, sizeStr := range sizes {
		// Remove leading and trailing whitespace
		sizeStr = strings.TrimSpace(sizeStr)

		// Skip empty strings (e.g., "250,,500" -> skip the empty middle part)
		if sizeStr == "" {
			continue
		}

		// Convert string to integer
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid package size '%s': %w", sizeStr, err)
		}

		// Validate that the package size is positive
		if size <= 0 {
			return nil, fmt.Errorf("package size must be positive, got %d", size)
		}

		// Add the valid package size to the result slice
		result = append(result, size)
	}

	// Ensure that at least one valid package size was found
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid package sizes found")
	}

	// Return the parsed package sizes
	return result, nil
}
