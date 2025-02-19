# Endor Battle Station

![AllianceStarbird](assets/alliance-logo.png)

A sophisticated battle station targeting system built in Go that manages ion cannon operations and target selection for the Rebel Alliance.

## Overview

The Endor Battle Station is a high-performance targeting system that:

- Processes radar scans from probe droids
- Applies complex targeting protocols
- Manages multiple ion cannons
- Executes coordinated attacks
- Provides real-time attack reporting

## Architecture

The system uses a hexagonal architecture (ports & adapters) with the following key components:

```
/
├── cmd/
│   └── battlestation/          # Main application entry point
├── internal/
│   ├── domain/                 # Core business logic
│   │   ├── target/            # Target selection logic
│   │   ├── protocol/          # Protocol implementations
│   │   ├── cannon/            # Ion cannon management
│   │   └── attack/            # Attack coordination
│   └── platform/              # External dependencies
│       ├── http/             # HTTP server & client
│       └── metrics/          # Performance monitoring
└── tests/                     # Integration tests
```

### Key Features

1. **Protocol Chain Processing**

   - Implements Chain of Responsibility pattern
   - Supports multiple simultaneous protocols
   - Efficient target filtering and prioritization
   - Handles edge cases gracefully

2. **Ion Cannon Management**

   - Priority-based cannon selection
   - Real-time availability tracking
   - Performance optimizations with caching
   - Concurrent operation support

3. **Performance Monitoring**
   - Prometheus metrics integration
   - Grafana dashboards
   - Response time tracking
   - System health monitoring

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.22 or later (for local development)
- Make (optional, for using Makefile commands)

### Running the System

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd battlestation-codetest
   ```

2. Start the system using Docker Compose:

   ```bash
   docker-compose up -d
   ```

3. The battle station API will be available at `http://localhost:8080`

### Development

1. Install dependencies:

   ```bash
   go mod download
   ```

2. Run tests:

   ```bash
   make test
   ```

3. Run integration tests:
   ```bash
   ./tests.sh
   ```

## API Documentation

### Attack Endpoint

POST `/attack`

Request body:

```json
{
  "protocols": ["avoid-mech"],
  "scan": [
    {
      "coordinates": { "x": 0, "y": 40 },
      "enemies": { "type": "soldier", "number": 10 }
    }
  ]
}
```

Response:

```json
{
  "target": { "x": 0, "y": 40 },
  "casualties": 1,
  "generation": 1
}
```

## Supported Protocols

- **closest-enemies**: Prioritize closest enemy point
- **furthest-enemies**: Prioritize furthest enemy point
- **assist-allies**: Prioritize enemy points with allies
- **avoid-crossfire**: Do not attack enemy points with allies
- **prioritize-mech**: Attack mech enemies if found
- **avoid-mech**: Do not attack any mech enemies

## Monitoring

The system includes Grafana dashboards for monitoring:

- Ion cannon availability
- Request latency
- Attack success rates
- System metrics

Access Grafana at `http://localhost:3000`

## Architecture Decisions

The project includes detailed Architecture Decision Records (ADRs):

- [ADR-001](docs/ADR-001-architecture.md): Overall Architecture Design
- [ADR-002](docs/ADR-002-target-selection.md): Target Selection Strategy
- [ADR-003](docs/ADR-003-ion-cannon.md): Ion Cannon Management

## Testing

The system includes:

- Unit tests for all components
- Integration tests using test cases
- Performance tests
- Mock implementations for external dependencies

Run the test suite:

```bash
make test
```

## Contributing

1. Follow Go best practices and idioms
2. Ensure tests pass
3. Update documentation as needed
4. Use conventional commits

## License

This project is licensed under the MIT License - see the LICENSE file for details.
