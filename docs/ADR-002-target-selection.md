# ADR 002: Target Selection Strategy

## Status

Accepted

## Context

The battle station must select targets based on multiple protocols that can be combined. The selection process must be:

- Fast and efficient
- Support multiple protocols simultaneously
- Handle edge cases gracefully
- Ignore targets beyond 100km

## Decision

We will implement a chain of responsibility pattern for protocol processing with the following design:

### Protocol Chain

1. Each protocol will be a separate handler in the chain
2. Protocols will be applied in order of specificity:
   - First: Validation protocols (avoid-mech, avoid-crossfire)
   - Second: Type protocols (prioritize-mech)
   - Third: Position protocols (closest-enemies, furthest-enemies)
   - Fourth: Tactical protocols (assist-allies)

### Target Filtering Process

1. **Initial Validation**

   - Filter out targets beyond 100km
   - Calculate distances for all targets

2. **Protocol Application**

   ```go
   type Target struct {
       Coordinates Position
       Enemies    EnemyGroup
       Allies     int
       Distance   float64
   }

   type TargetFilter interface {
       Apply(targets []Target, scan ScanData) []Target
   }
   ```

3. **Protocol Implementations**
   - avoid-mech: Filter out all mech targets
   - avoid-crossfire: Filter out targets with allies
   - prioritize-mech: Sort mech targets first
   - closest/furthest-enemies: Sort by distance
   - assist-allies: Prioritize targets with allies

### Distance Calculation

```go
func calculateDistance(pos Position) float64 {
    return math.Sqrt(float64(pos.X*pos.X + pos.Y*pos.Y))
}
```

## Consequences

### Positive

- Clear separation of protocol logic
- Easy to add new protocols
- Predictable target selection
- Efficient filtering and sorting
- Simple to test each protocol independently

### Negative

- Need to carefully order protocol application
- May need to process all targets for some protocols
- Memory usage scales with number of targets

## Performance Considerations

1. Early filtering of invalid targets
2. Pre-calculation of distances
3. Use of efficient sorting algorithms
4. Minimize memory allocations

## Edge Cases

1. No valid targets after filtering
2. Multiple targets at same distance
3. Conflicting protocol requirements
4. All targets beyond range
5. Empty scan data

## Testing Strategy

1. Unit tests for each protocol
2. Integration tests for protocol combinations
3. Performance tests for large scan data
4. Edge case validation
5. Protocol ordering tests
