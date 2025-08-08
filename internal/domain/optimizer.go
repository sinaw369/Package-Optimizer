package domain

import (
	"fmt"
	"sort"
)

// Optimizer handles package optimization calculations using dynamic programming.
// It finds the optimal combination of packages that minimizes over-delivery
// while using the fewest number of packages when over-delivery is tied.
type Optimizer struct {
	// packageSizes stores available package sizes in descending order for efficiency
	packageSizes []int
}

// NewOptimizer creates a new optimizer with the given package sizes
func NewOptimizer(packageSizes []int) *Optimizer {
	// Validate that package sizes list is not empty
	if len(packageSizes) == 0 {
		panic("package sizes cannot be empty")
	}

	// Validate that all package sizes are positive integers
	for _, size := range packageSizes {
		if size <= 0 {
			panic("package sizes must be positive")
		}
	}

	// Sort package sizes in descending order for efficiency in dynamic programming
	// This allows us to try larger packages first, which often leads to better solutions
	sizes := make([]int, len(packageSizes))
	copy(sizes, packageSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	return &Optimizer{
		packageSizes: sizes,
	}
}

// Optimize calculates the optimal package combination for the given quantity.
// It uses dynamic programming to find the solution that:
// 1. Minimizes over-delivery (total_delivered - requested)
// 2. Minimizes the number of packages used (when over-delivery is tied)
func (o *Optimizer) Optimize(quantity int) (*OptimizationResult, error) {
	// Validate that quantity is non-negative
	if quantity < 0 {
		return nil, fmt.Errorf("quantity must be non-negative, got %d", quantity)
	}

	// Handle edge case: zero quantity requires no packages
	if quantity == 0 {
		return &OptimizationResult{
			Requested:      0,
			TotalDelivered: 0,
			OverDelivery:   0,
			Packages:       make(map[string]int),
		}, nil
	}

	// Use dynamic programming algorithm to find the optimal solution
	solution := o.findOptimalSolution(quantity)

	// Convert the internal solution format to the public result format
	result := &OptimizationResult{
		Requested:      quantity,
		TotalDelivered: solution.totalDelivered,
		OverDelivery:   solution.totalDelivered - quantity,
		Packages:       make(map[string]int),
	}

	// Convert package counts from internal format to string map for JSON response
	for _, pkg := range solution.packages {
		if pkg.Count > 0 {
			result.Packages[fmt.Sprintf("%d", pkg.Size)] = pkg.Count
		}
	}

	return result, nil
}

// solution represents a complete solution with package counts.
// This is an internal structure used by the dynamic programming algorithm.
type solution struct {
	totalDelivered int            // Total quantity delivered
	packages       []PackageCount // List of packages used with their counts
}

// findOptimalSolution uses dynamic programming to find the optimal package combination.
//
// Algorithm Overview:
// 1. Create a DP table where dp[i] represents the minimum over-delivery for quantity i
// 2. For each quantity i, try using each available package size
// 3. Update the solution if we find a better combination (less over-delivery or fewer packages)
// 4. Track package combinations for each quantity
//
// Time Complexity: O(n × m) where n is the requested quantity and m is the number of package sizes
// Space Complexity: O(n) for the DP arrays
//
// Args:
//   - quantity: the requested quantity
//
// Returns:
//   - *solution: the optimal solution found
func (o *Optimizer) findOptimalSolution(quantity int) *solution {
	// Calculate the maximum quantity we need to consider
	// We need to handle quantities up to quantity + maxPackageSize to find optimal solutions
	maxPackageSize := o.packageSizes[0] // Largest package size (first after sorting)
	maxQuantity := quantity + maxPackageSize

	// Initialize DP arrays
	// dp[i] represents the minimum over-delivery for quantity i
	dp := make([]int, maxQuantity+1)
	// packageCounts[i] stores the package combination for quantity i
	packageCounts := make([][]PackageCount, maxQuantity+1)

	// Initialize DP table with "infinity" (large number) to represent unreachable states
	for i := range dp {
		dp[i] = maxQuantity + 1
	}
	// Base case: quantity 0 requires 0 packages and has 0 over-delivery
	dp[0] = 0
	packageCounts[0] = []PackageCount{}

	// Fill the DP table using bottom-up approach
	for i := 1; i <= maxQuantity; i++ {
		// Try each available package size
		for _, packageSize := range o.packageSizes {
			// Only consider packages that can fit in the current quantity
			if packageSize <= i {
				// Calculate remaining quantity after using this package
				remaining := i - packageSize
				// Calculate new over-delivery for this quantity
				newOverDelivery := max(0, i-quantity)

				// Check if we can reach the remaining quantity
				if dp[remaining] != maxQuantity+1 {
					// Calculate total over-delivery for this combination
					totalOverDelivery := dp[remaining] + newOverDelivery

					// Update if this is a better solution:
					// 1. Less over-delivery, OR
					// 2. Same over-delivery but fewer packages
					if totalOverDelivery < dp[i] ||
						(totalOverDelivery == dp[i] && len(packageCounts[remaining])+1 < len(packageCounts[i])) {

						// Update the minimum over-delivery
						dp[i] = totalOverDelivery

						// Copy existing packages from the remaining quantity
						newPackages := make([]PackageCount, len(packageCounts[remaining]))
						copy(newPackages, packageCounts[remaining])

						// Add or update the current package count
						found := false
						for j := range newPackages {
							if newPackages[j].Size == packageSize {
								newPackages[j].Count++
								found = true
								break
							}
						}
						// If package size not found, add it as a new package
						if !found {
							newPackages = append(newPackages, PackageCount{Size: packageSize, Count: 1})
						}

						// Store the package combination for this quantity
						packageCounts[i] = newPackages
					}
				}
			}
		}
	}

	// Find the best solution among all quantities >= requested
	// Start with the requested quantity
	bestQuantity := quantity
	bestOverDelivery := dp[quantity]

	// Check all quantities from requested+1 to maxQuantity for better solutions
	for i := quantity + 1; i <= maxQuantity; i++ {
		// Update if we find a better solution (less over-delivery or fewer packages)
		if dp[i] < bestOverDelivery ||
			(dp[i] == bestOverDelivery && len(packageCounts[i]) < len(packageCounts[bestQuantity])) {
			bestQuantity = i
			bestOverDelivery = dp[i]
		}
	}

	// Return the optimal solution found
	return &solution{
		totalDelivered: bestQuantity,
		packages:       packageCounts[bestQuantity],
	}
}

// max returns the maximum of two integers.
// This is a utility function used in the dynamic programming algorithm.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
