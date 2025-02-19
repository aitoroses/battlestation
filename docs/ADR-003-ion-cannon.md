# ADR 003: Ion Cannon Management Strategy

## Status

Accepted

## Context

The battle station must manage 3 ion cannons with different characteristics:

- Each has different generation (1st, 2nd, 3rd) and fire times (3.5s, 1.5s, 2.5s)
- Must prioritize lower generation cannons
- Cannons become unavailable during fire time
- Must handle 1 request/second throughput
- Need to check cannon status via HTTP

## Decision

We will implement a priority-based cannon manager with the following design:

### Cannon Manager

```go
type IonCannon struct {
    Generation  int
    FireTime   float64
    BaseURL    string
    LastFired  time.Time
}

type CannonManager struct {
    cannons    []*IonCannon
    httpClient HTTPClient
    mu         sync.RWMutex
}
```

### Key Components

1. **Cannon Priority Queue**

   - Ordered by generation (lowest first)
   - Track availability based on fire time
   - Maintain thread-safe state

2. **Status Checking**

   - Parallel status checks for all cannons
   - Caching of status results
   - Circuit breaker for failed requests

3. **Selection Strategy**

   ```go
   func (cm *CannonManager) GetBestAvailable() (*IonCannon, error) {
       cm.mu.RLock()
       defer cm.mu.RUnlock()

       for _, cannon := range cm.cannons {
           if cannon.IsAvailable() && cannon.CheckStatus() {
               return cannon, nil
           }
       }
       return nil, ErrNoCannonAvailable
   }
   ```

4. **Availability Tracking**
   - Time-based availability calculation
   - Status endpoint health check
   - Automatic recovery after fire time

### Performance Optimizations

1. **Status Check Caching**

   - Cache status results for 100ms
   - Refresh cache in background
   - Use stale cache on errors

2. **Parallel Processing**

   - Check all cannon statuses concurrently
   - Use context for timeout control
   - Cancel unnecessary checks early

3. **Circuit Breaker**
   - Prevent cascade failures
   - Quick failure for known bad states
   - Exponential backoff for retries

## Consequences

### Positive

- Efficient cannon selection
- Reliable availability tracking
- Handles network issues gracefully
- Maintains optimal throughput
- Easy to monitor and debug

### Negative

- Complex state management
- Network dependency for status
- Potential race conditions
- Cache consistency challenges

## Implementation Notes

1. **Error Handling**

   - Timeout for status checks (100ms)
   - Fallback strategies for failures
   - Clear error messages

2. **Monitoring**

   - Track cannon usage metrics
   - Monitor status check latency
   - Alert on availability issues

3. **Testing**
   - Simulate network latency
   - Test concurrent requests
   - Verify priority ordering
   - Check edge cases

## Edge Cases

1. All cannons unavailable
2. Network timeouts
3. Invalid status responses
4. Clock skew issues
5. Concurrent fire requests

## Metrics

1. Cannon availability percentage
2. Status check latency
3. Selection time
4. Cache hit ratio
5. Error rates
