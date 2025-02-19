# OBR 003: API Operations and Error Handling

## Business Rules

### BR-1: API Endpoints

1. **Attack Endpoint**

   - Path: `/attack`
   - Method: POST
   - Content-Type: application/json
   - Required fields: protocols, scan

2. **Response Format**
   - Content-Type: application/json
   - Required fields: target, casualties, generation
   - Status codes: 200, 400, 500

### BR-2: Input Validation

1. **Protocols Array**

   - Must be non-empty
   - Must contain valid protocol names
   - Must not contain conflicting protocols

2. **Scan Data**
   - Must contain at least one target
   - Must have valid coordinates
   - Must have valid enemy data
   - Optional allies field

### BR-3: Performance Requirements

1. Must handle 1 request per second
2. Status check timeout: 100ms
3. Fire request timeout: 500ms
4. Total request timeout: 1000ms

## Use Cases

### UC-1: Valid Request Processing

1. **Simple Attack**

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

   - Expected: Successful attack
   - Test: test_cases.txt line 1

2. **Complex Attack**
   ```json
   {
     "protocols": ["closest-enemies", "prioritize-mech"],
     "scan": [
       {
         "coordinates": { "x": 0, "y": 1 },
         "enemies": { "type": "mech", "number": 1 }
       },
       {
         "coordinates": { "x": 0, "y": 10 },
         "enemies": { "type": "soldier", "number": 10 }
       }
     ]
   }
   ```
   - Expected: Select closest mech
   - Test: test_cases.txt line 12

### UC-2: Error Scenarios

1. **Invalid Protocol**

   ```json
   {
     "protocols": ["invalid-protocol"],
     "scan": [...]
   }
   ```

   - Expected: 400 Bad Request
   - Error: "Invalid protocol specified"

2. **Missing Required Fields**

   ```json
   {
     "protocols": []
   }
   ```

   - Expected: 400 Bad Request
   - Error: "Missing required scan data"

3. **Invalid Coordinates**
   ```json
   {
     "protocols": ["avoid-mech"],
     "scan": [
       {
         "coordinates": { "x": "invalid" },
         "enemies": { "type": "soldier", "number": 10 }
       }
     ]
   }
   ```
   - Expected: 400 Bad Request
   - Error: "Invalid coordinate format"

## Edge Cases

### EC-1: Input Validation

1. **Empty Arrays**

   - Empty protocols array
   - Empty scan array
   - Expected: 400 Bad Request

2. **Invalid Types**

   - Non-numeric coordinates
   - Non-string protocols
   - Non-numeric enemy numbers
   - Expected: 400 Bad Request

3. **Boundary Values**
   - Zero enemies
   - Negative coordinates
   - Extremely large numbers
   - Expected: Appropriate error handling

### EC-2: Performance Edge Cases

1. **Large Payload**

   - Maximum number of scan points
   - Multiple protocols
   - Expected: Handle within timeout

2. **Concurrent Requests**
   - Multiple simultaneous requests
   - Rapid sequential requests
   - Expected: Maintain throughput

### EC-3: System Errors

1. **Cannon Service Unavailable**

   - All cannons down
   - Network errors
   - Expected: 503 Service Unavailable

2. **Timeout Scenarios**
   - Status check timeout
   - Fire request timeout
   - Expected: 504 Gateway Timeout

## Test Scenarios

### TS-1: Input Validation Tests

1. **Protocol Validation**

   - Test all valid protocols
   - Test invalid protocols
   - Test protocol combinations

2. **Scan Data Validation**
   - Test coordinate ranges
   - Test enemy types
   - Test ally numbers

### TS-2: Performance Tests

1. **Load Testing**

   - Sustained 1 req/sec
   - Burst testing
   - Response time monitoring

2. **Concurrency Testing**
   - Parallel requests
   - Resource utilization
   - Error rate monitoring

### TS-3: Error Handling Tests

1. **Client Errors**

   - Invalid input formats
   - Missing fields
   - Invalid values

2. **Server Errors**
   - Service unavailable
   - Timeout handling
   - Recovery testing

### TS-4: Integration Tests

1. **End-to-End Flow**

   - Request validation
   - Target selection
   - Cannon operation
   - Response formatting

2. **System Interaction**
   - Cannon service communication
   - Error propagation
   - State management

## Monitoring Metrics

### M-1: Request Metrics

1. Request rate
2. Response times
3. Error rates
4. Status code distribution

### M-2: Business Metrics

1. Protocol usage
2. Target selection patterns
3. Cannon utilization
4. Attack success rate

### M-3: System Metrics

1. CPU usage
2. Memory usage
3. Network latency
4. Error logs
