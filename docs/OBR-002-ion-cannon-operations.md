# OBR 002: Ion Cannon Operational Rules

## Business Rules

### BR-1: Cannon Generation Priority

1. Always select lowest generation available cannon
2. Generation order:
   - 1st Generation (3.5s fire time)
   - 2nd Generation (1.5s fire time)
   - 3rd Generation (2.5s fire time)

### BR-2: Cannon Availability

1. Cannon becomes unavailable immediately after firing
2. Unavailability duration equals fire time
3. Must check HTTP status endpoint before selection
4. All cannons may be unavailable simultaneously

### BR-3: Request Handling

1. Must handle 1 request per second
2. Cannot delay requests waiting for cannons
3. Must return error if no cannons available
4. Must report actual casualties and generation used

## Use Cases

### UC-1: Basic Cannon Selection

1. **All Cannons Available**

   - Expected: Select 1st Generation
   - Test: test_cases.txt line 1

2. **1st Gen Unavailable**

   - Expected: Select 2nd Generation
   - Test: test_cases.txt line 2

3. **1st and 2nd Gen Unavailable**
   - Expected: Select 3rd Generation
   - Test: test_cases.txt line 3

### UC-2: Rapid Fire Scenarios

1. **Sequential Requests**

   - Input: Multiple requests within 5 seconds
   - Expected: Rotate through all cannons
   - Test: Consecutive lines in test_cases.txt

2. **Concurrent Requests**
   - Input: Multiple simultaneous requests
   - Expected: Different cannons for each request
   - Test: Load testing scenarios

### UC-3: Status Check Scenarios

1. **Normal Operation**

   - All cannons report available
   - Status checks complete quickly
   - Select lowest generation

2. **Degraded Operation**

   - Some cannons report unavailable
   - Skip unavailable cannons
   - Select lowest available generation

3. **Status Check Failures**
   - Timeout or error from status endpoint
   - Mark cannon as unavailable
   - Select next available cannon

## Edge Cases

### EC-1: Availability Edge Cases

1. **All Cannons Unavailable**

   - Due to fire time
   - Due to status check failures
   - Due to combination of both
   - Expected: Return error response

2. **Cannon Recovery**

   - Just finished fire time
   - Status check success after previous failure
   - Expected: Include in selection pool

3. **Status Check Timing**
   - Slow status response
   - Timeout scenarios
   - Network errors
   - Expected: Fail fast, try next cannon

### EC-2: Timing Edge Cases

1. **Fire Time Boundaries**

   - Request exactly at availability transition
   - Multiple cannons becoming available simultaneously
   - Expected: Consistent selection behavior

2. **Request Rate Spikes**
   - Burst of requests above 1/second
   - Requests during status check
   - Expected: Handle gracefully without errors

### EC-3: Error Cases

1. **Status Endpoint Errors**

   - Network timeout
   - Invalid response format
   - 5xx server errors
   - Expected: Handle gracefully, try next cannon

2. **Fire Endpoint Errors**
   - Network timeout
   - Invalid response format
   - 5xx server errors
   - Expected: Return error, keep cannon state consistent

## Test Scenarios

### TS-1: Basic Operation Tests

1. **Single Cannon Tests**

   - Fire and verify unavailability
   - Wait for recovery
   - Verify availability

2. **Multiple Cannon Tests**
   - Verify priority selection
   - Test rotation under load
   - Check status handling

### TS-2: Performance Tests

1. **Throughput Testing**

   - Verify 1 request/second handling
   - Test burst handling
   - Measure response times

2. **Status Check Performance**
   - Measure status check overhead
   - Test caching effectiveness
   - Verify timeout handling

### TS-3: Reliability Tests

1. **Network Issues**

   - Test with network delays
   - Handle connection failures
   - Verify recovery behavior

2. **Error Handling**
   - Test all error scenarios
   - Verify error responses
   - Check system recovery

### TS-4: Integration Tests

1. **Full Attack Sequence**

   - Target selection
   - Cannon selection
   - Fire execution
   - Response validation

2. **Load Testing**
   - Sustained request rate
   - Error rate monitoring
   - Performance metrics

### TS-5: Monitoring Metrics

1. **Availability Metrics**

   - Cannon availability percentage
   - Status check success rate
   - Error rates

2. **Performance Metrics**
   - Response times
   - Status check latency
   - Request queue length
