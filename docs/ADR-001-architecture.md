# ADR 001: Overall Architecture Design

## Status

Accepted

## Context

We need to design a battle station targeting system that:

1. Receives radar scans via HTTP API
2. Processes attack protocols
3. Manages ion cannon availability
4. Executes attacks
5. Reports results

## Decision

We will use a hexagonal architecture (ports & adapters) with the following structure:

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
│   ├── platform/              # External dependencies
│   │   ├── http/             # HTTP server & client
│   │   └── metrics/          # Performance monitoring
│   └── ports/                 # Interface definitions
│       ├── incoming/         # HTTP handlers
│       └── outgoing/         # External services
└── tests/                     # Integration tests
```

### Key Components

1. **Domain Layer**

   - Target Selection: Implements target prioritization logic
   - Protocol Handler: Processes and combines multiple attack protocols
   - Cannon Manager: Tracks ion cannon availability and selection
   - Attack Coordinator: Orchestrates the attack sequence

2. **Platform Layer**

   - HTTP Server: Handles incoming attack requests
   - HTTP Client: Communicates with ion cannons
   - Metrics: Monitors system performance

3. **Ports Layer**
   - Incoming: HTTP API handlers
   - Outgoing: Ion cannon service interfaces

## Consequences

### Positive

- Clear separation of concerns
- Easy to test business logic in isolation
- Simple to add new protocols or modify existing ones
- Flexible ion cannon management
- Easy to extend with new features

### Negative

- More initial setup required
- More files and directories to manage
- Potential overhead in small operations

## Notes

- All components must be designed for concurrent operation
- Response time is critical - must handle 1 request/second
- Must maintain ion cannon availability state
