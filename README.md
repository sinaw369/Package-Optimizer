# Package Optimizer

A Go backend service that optimizes package delivery by minimizing over-delivery while using configurable fixed-size packages.

## Problem Description

Given a requested quantity and a set of available package sizes, the service calculates the optimal combination of packages that:
1. **Primary goal**: Minimizes over-delivery (total_delivered - requested)
2. **Secondary goal**: Minimizes the total number of packages used (when over-delivery is tied)

## Features

- HTTP API endpoint for package optimization calculations
- Configurable package sizes via environment variables
- Comprehensive unit tests
- Docker containerization
- Simple web UI for testing
- Well-structured codebase with separate layers
- Built with Echo framework for high performance

## Quick Start

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd package-optimizer
```

2. Run with Docker Compose:
```bash
docker-compose up --build
```

3. Access the service:
- API: http://localhost:8080/calculate?qty=1201
- Web UI: http://localhost:8080

### Manual Setup

1. Install Go 1.21+ and ensure it's in your PATH

2. Set package sizes (optional, defaults to 250,500,1000,2000):
```bash
export PACKAGE_SIZES="250,500,1000,2000"
```

3. Run the service:
```bash
go run cmd/server/main.go
```

4. Run tests:
```bash
go test ./...
```

## API Usage

### Calculate Optimal Packages

**Endpoint**: `GET /calculate?qty={quantity}`

**Example**:
```bash
curl "http://localhost:8080/calculate?qty=1201"
```

**Response**:
```json
{
  "requested": 1201,
  "total_delivered": 1250,
  "over_delivery": 49,
  "packages": {
    "250": 5
  }
}
```

## Configuration

### Environment Variables

- `PACKAGE_SIZES`: Comma-separated list of available package sizes (default: "250,500,1000,2000")
- `PORT`: Server port (default: 8080)

### Example Configuration

```bash
export PACKAGE_SIZES="100,200,500,1000"
export PORT=3000
```

## Project Structure

```
package-optimizer/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── api/
│   │   ├── handler.go       # HTTP handlers (Echo framework)
│   │   └── middleware.go    # HTTP middleware (Echo framework)
│   ├── domain/
│   │   ├── optimizer.go     # Core optimization logic
│   │   └── types.go         # Domain types
│   └── config/
│       └── config.go        # Configuration management
├── web/
│   └── static/
│       ├── index.html       # Web UI
│       ├── style.css        # CSS styles
│       └── script.js        # JavaScript logic
├── tests/
│   └── optimizer_test.go    # Unit tests
├── Dockerfile               # Docker configuration
├── docker-compose.yml       # Docker Compose setup
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies checksum
└── README.md               # This file
```

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Specific Test File
```bash
go test ./tests/optimizer_test.go
```

## Docker Commands

### Build Image
```bash
docker build -t package-optimizer .
```

### Run Container
```bash
docker run -p 8080:8080 package-optimizer
```

### Run with Custom Package Sizes
```bash
docker run -p 8080:8080 -e PACKAGE_SIZES="100,200,500" package-optimizer
```

## Algorithm

The optimization algorithm uses a dynamic programming approach:

1. **State**: `dp[i]` represents the minimum over-delivery for quantity `i`
2. **Transition**: For each package size, try using it and update the minimum over-delivery
3. **Tie-breaking**: When over-delivery is equal, prefer fewer packages

### Time Complexity
- O(n × m) where n is the requested quantity and m is the number of package sizes
- Space complexity: O(n)

## Edge Cases Handled

- Zero quantity (returns empty result)
- Negative quantity (returns error)
- Very large quantities (handled efficiently)
- Invalid package sizes (validated)
- Empty package sizes list (returns error)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### MIT License Summary

- ✅ **Commercial Use**: Allowed
- ✅ **Modification**: Allowed  
- ✅ **Distribution**: Allowed
- ✅ **Private Use**: Allowed
- ❌ **Liability**: Limited
- ❌ **Warranty**: Limited

The MIT License is a permissive license that allows others to use, modify, and distribute your code with very few restrictions. 