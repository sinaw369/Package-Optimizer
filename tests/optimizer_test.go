package tests

import (
	"testing"

	"package-optimizer/internal/domain"
)

func TestOptimizer_Optimize(t *testing.T) {
	tests := []struct {
		name           string
		packageSizes   []int
		quantity       int
		expectedResult *domain.OptimizationResult
		expectError    bool
	}{
		{
			name:         "Basic case - exact match",
			packageSizes: []int{250, 500, 1000, 2000},
			quantity:     1000,
			expectedResult: &domain.OptimizationResult{
				Requested:      1000,
				TotalDelivered: 1000,
				OverDelivery:   0,
				Packages: map[string]int{
					"1000": 1,
				},
			},
			expectError: false,
		},
		{
			name:         "Basic case - over delivery",
			packageSizes: []int{250, 500, 1000, 2000},
			quantity:     1201,
			expectedResult: &domain.OptimizationResult{
				Requested:      1201,
				TotalDelivered: 1250,
				OverDelivery:   49,
				Packages: map[string]int{
					"1000": 1,
					"250":  1,
				},
			},
			expectError: false,
		},
		{
			name:         "Zero quantity",
			packageSizes: []int{250, 500, 1000, 2000},
			quantity:     0,
			expectedResult: &domain.OptimizationResult{
				Requested:      0,
				TotalDelivered: 0,
				OverDelivery:   0,
				Packages:       map[string]int{},
			},
			expectError: false,
		},
		{
			name:         "Small quantity with large packages",
			packageSizes: []int{1000, 2000},
			quantity:     500,
			expectedResult: &domain.OptimizationResult{
				Requested:      500,
				TotalDelivered: 1000,
				OverDelivery:   500,
				Packages: map[string]int{
					"1000": 1,
				},
			},
			expectError: false,
		},
		{
			name:         "Large quantity",
			packageSizes: []int{250, 500, 1000, 2000},
			quantity:     5000,
			expectedResult: &domain.OptimizationResult{
				Requested:      5000,
				TotalDelivered: 5000,
				OverDelivery:   0,
				Packages: map[string]int{
					"2000": 2,
					"1000": 1,
				},
			},
			expectError: false,
		},
		{
			name:         "Negative quantity",
			packageSizes: []int{250, 500, 1000, 2000},
			quantity:     -100,
			expectError:  true,
		},
		{
			name:         "Single package size",
			packageSizes: []int{100},
			quantity:     250,
			expectedResult: &domain.OptimizationResult{
				Requested:      250,
				TotalDelivered: 300,
				OverDelivery:   50,
				Packages: map[string]int{
					"100": 3,
				},
			},
			expectError: false,
		},
		{
			name:         "Tie breaking - prefer fewer packages",
			packageSizes: []int{100, 200},
			quantity:     150,
			expectedResult: &domain.OptimizationResult{
				Requested:      150,
				TotalDelivered: 200,
				OverDelivery:   50,
				Packages: map[string]int{
					"200": 1,
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := domain.NewOptimizer(tt.packageSizes)
			result, err := optimizer.Optimize(tt.quantity)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Requested != tt.expectedResult.Requested {
				t.Errorf("Requested = %v, want %v", result.Requested, tt.expectedResult.Requested)
			}

			if result.TotalDelivered != tt.expectedResult.TotalDelivered {
				t.Errorf("TotalDelivered = %v, want %v", result.TotalDelivered, tt.expectedResult.TotalDelivered)
			}

			if result.OverDelivery != tt.expectedResult.OverDelivery {
				t.Errorf("OverDelivery = %v, want %v", result.OverDelivery, tt.expectedResult.OverDelivery)
			}

			// Check packages
			if len(result.Packages) != len(tt.expectedResult.Packages) {
				t.Errorf("Packages count = %v, want %v", len(result.Packages), len(tt.expectedResult.Packages))
			}

			for size, count := range tt.expectedResult.Packages {
				if result.Packages[size] != count {
					t.Errorf("Package %s count = %v, want %v", size, result.Packages[size], count)
				}
			}
		})
	}
}

func TestOptimizer_EdgeCases(t *testing.T) {
	t.Run("Very large quantity", func(t *testing.T) {
		optimizer := domain.NewOptimizer([]int{1, 2, 5, 10})
		result, err := optimizer.Optimize(10000)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if result.Requested != 10000 {
			t.Errorf("Requested = %v, want 10000", result.Requested)
		}

		if result.TotalDelivered < result.Requested {
			t.Errorf("TotalDelivered (%v) should be >= Requested (%v)", result.TotalDelivered, result.Requested)
		}
	})

	t.Run("Quantity equals smallest package", func(t *testing.T) {
		optimizer := domain.NewOptimizer([]int{100, 200, 500})
		result, err := optimizer.Optimize(100)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if result.TotalDelivered != 100 {
			t.Errorf("TotalDelivered = %v, want 100", result.TotalDelivered)
		}

		if result.OverDelivery != 0 {
			t.Errorf("OverDelivery = %v, want 0", result.OverDelivery)
		}
	})

	t.Run("Quantity between package sizes", func(t *testing.T) {
		optimizer := domain.NewOptimizer([]int{100, 300, 500})
		result, err := optimizer.Optimize(200)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if result.TotalDelivered < result.Requested {
			t.Errorf("TotalDelivered (%v) should be >= Requested (%v)", result.TotalDelivered, result.Requested)
		}
	})
}

func TestOptimizer_Validation(t *testing.T) {
	t.Run("Empty package sizes", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty package sizes")
			}
		}()

		domain.NewOptimizer([]int{})
	})

	t.Run("Negative package size", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative package size")
			}
		}()

		domain.NewOptimizer([]int{-100, 200})
	})
}

func BenchmarkOptimizer_Optimize(b *testing.B) {
	optimizer := domain.NewOptimizer([]int{250, 500, 1000, 2000})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := optimizer.Optimize(1000 + i%5000)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
