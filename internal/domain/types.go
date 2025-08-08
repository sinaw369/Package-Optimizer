package domain

// OptimizationResult represents the result of a package optimization calculation.
// This structure is returned by the optimizer and contains all the information
// about the optimal package combination for a given quantity.
type OptimizationResult struct {
	// Requested is the original quantity that was requested
	Requested int `json:"requested"`

	// TotalDelivered is the total quantity that will be delivered
	// This may be greater than or equal to the requested quantity
	TotalDelivered int `json:"total_delivered"`

	// OverDelivery is the excess quantity delivered beyond what was requested
	// Calculated as: TotalDelivered - Requested
	OverDelivery int `json:"over_delivery"`

	// Packages is a map of package sizes to their counts
	// Key: package size as string (e.g., "250", "500", "1000")
	// Value: number of packages of that size to use
	Packages map[string]int `json:"packages"`
}

// PackageCount represents a package size and its count in a solution.
// This is an internal structure used by the optimizer to track package combinations.
type PackageCount struct {
	// Size is the package size (e.g., 250, 500, 1000)
	Size int
	// Count is the number of packages of this size to use
	Count int
}

// OptimizationRequest represents a request for package optimization.
// This structure can be used for future API extensions that accept JSON requests.
type OptimizationRequest struct {
	// Quantity is the requested quantity to be delivered
	Quantity int
}

// ErrorResponse represents an error response from the API.
// This structure is used to return consistent error messages in JSON format.
type ErrorResponse struct {
	// Error is the error message describing what went wrong
	Error string `json:"error"`
}
