# OBR 001: Attack Protocol Rules

## Business Rules

### BR-1: Protocol Compatibility

- Protocols must be compatible with each other
- Incompatible combinations (e.g., closest-enemies + furthest-enemies) will not be provided

### BR-2: Protocol Priority

1. **Validation Protocols** (First to apply)

   - avoid-mech: Skip all mech-type enemies
   - avoid-crossfire: Skip all positions with allies

2. **Type Protocols** (Second to apply)

   - prioritize-mech: Must select mech-type enemies if available

3. **Position Protocols** (Third to apply)

   - closest-enemies: Select nearest valid target
   - furthest-enemies: Select farthest valid target

4. **Tactical Protocols** (Last to apply)
   - assist-allies: Prioritize positions with allies present

### BR-3: Distance Rules

- Targets beyond 100km must be ignored
- Distance is calculated from origin (0,0)
- Distance formula: sqrt(x² + y²)

## Use Cases

### UC-1: Single Protocol Application

1. **avoid-mech**

   - Input: Mix of mech and soldier targets
   - Expected: Only soldier targets considered
   - Test: test_cases.txt line 1

2. **prioritize-mech**

   - Input: Mix of mech and soldier targets
   - Expected: Mech target selected if available
   - Test: test_cases.txt line 2

3. **closest-enemies**

   - Input: Multiple targets at different distances
   - Expected: Nearest valid target selected
   - Test: test_cases.txt line 3

4. **furthest-enemies**

   - Input: Multiple targets at different distances
   - Expected: Farthest valid target selected
   - Test: test_cases.txt line 4

5. **assist-allies**

   - Input: Mix of targets with/without allies
   - Expected: Target with allies selected
   - Test: test_cases.txt line 5

6. **avoid-crossfire**
   - Input: Mix of targets with/without allies
   - Expected: Only targets without allies considered
   - Test: test_cases.txt line 6

### UC-2: Multiple Protocol Combinations

1. **closest-enemies + avoid-mech**

   - Expected: Nearest non-mech target
   - Test: test_cases.txt line 8

2. **furthest-enemies + avoid-mech**

   - Expected: Farthest non-mech target
   - Test: test_cases.txt line 11

3. **closest-enemies + prioritize-mech**
   - Expected: Nearest mech target if available
   - Test: test_cases.txt line 12

## Edge Cases

### EC-1: Distance Edge Cases

1. Target exactly at 100km
   - Should be included in valid targets
2. Target slightly over 100km
   - Should be excluded from valid targets
3. Target at coordinates (0,0)
   - Should be considered closest target

### EC-2: Target Selection Edge Cases

1. Multiple targets at same distance
   - Should use stable sorting to maintain consistency
2. Single valid target after filtering
   - Should select regardless of other protocols
3. No valid targets after filtering
   - Should return appropriate error

### EC-3: Protocol Edge Cases

1. Empty protocols list
   - Should select based on default behavior
2. Single target matching all protocols
   - Should be selected regardless of other factors
3. No targets matching all protocols
   - Should return appropriate error

## Test Scenarios

### TS-1: Basic Protocol Tests

- Test each protocol individually
- Verify expected target selection
- Check casualty counts
- Validate response format

### TS-2: Protocol Combination Tests

- Test common protocol combinations
- Verify protocol priority order
- Check edge case handling
- Validate error responses

### TS-3: Performance Tests

- Test with maximum number of targets
- Verify 1 request/second handling
- Check response times
- Monitor resource usage

### TS-4: Error Handling Tests

- Test invalid input formats
- Check error response formats
- Verify error logging
- Validate error recovery

### TS-5: Integration Tests

- Test full attack sequence
- Verify cannon selection
- Check attack execution
- Validate final response
