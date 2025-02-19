package cannon

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// Manager handles the coordination of multiple ion cannons
type Manager struct {
	cannons []*IonCannon
	mu      sync.RWMutex
}

// NewManager creates a new cannon manager
func NewManager(cannons []*IonCannon) *Manager {
	// Sort cannons by generation to ensure consistent priority
	sort.Slice(cannons, func(i, j int) bool {
		return cannons[i].Generation() < cannons[j].Generation()
	})

	return &Manager{
		cannons: cannons,
	}
}

// GetBestAvailable finds the best available cannon based on generation priority
func (m *Manager) GetBestAvailable(ctx context.Context) (*IonCannon, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	type result struct {
		cannon *IonCannon
		status *Status
		err    error
	}

	// Check all cannons concurrently
	results := make(chan result, len(m.cannons))
	var wg sync.WaitGroup

	for _, c := range m.cannons {
		wg.Add(1)
		go func(cannon *IonCannon) {
			defer wg.Done()

			// Skip if not available based on fire time
			if !cannon.IsAvailable() {
				results <- result{err: fmt.Errorf("cannon generation %d not available", cannon.Generation())}
				return
			}

			// Check HTTP status
			status, err := cannon.CheckStatus(ctx)
			if err != nil {
				results <- result{err: fmt.Errorf("cannon generation %d status check failed: %w", cannon.Generation(), err)}
				return
			}

			results <- result{
				cannon: cannon,
				status: status,
			}
		}(c)
	}

	// Close results channel when all checks are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and find best available cannon
	var bestCannon *IonCannon
	var bestGeneration Generation = 4 // Higher than any valid generation

	for r := range results {
		if r.err != nil {
			continue
		}

		if !r.status.Available {
			continue
		}

		// Update best cannon if this one has lower generation
		if Generation(r.status.Generation) < bestGeneration {
			bestCannon = r.cannon
			bestGeneration = Generation(r.status.Generation)
		}
	}

	if bestCannon == nil {
		return nil, fmt.Errorf("no cannons available")
	}

	return bestCannon, nil
}

// Fire attempts to fire the specified cannon at the target
func (m *Manager) Fire(ctx context.Context, cannon *IonCannon, req *FireRequest) (*FireResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Verify cannon is still in our list
	found := false
	for _, c := range m.cannons {
		if c == cannon {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("invalid cannon")
	}

	// Attempt to fire
	resp, err := cannon.Fire(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fire failed: %w", err)
	}

	return resp, nil
}

// GetStatus returns the current status of all cannons
func (m *Manager) GetStatus(ctx context.Context) map[Generation]*Status {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[Generation]*Status)
	var wg sync.WaitGroup
	var statusMu sync.Mutex

	for _, c := range m.cannons {
		wg.Add(1)
		go func(cannon *IonCannon) {
			defer wg.Done()

			s, err := cannon.CheckStatus(ctx)
			if err != nil {
				return
			}

			statusMu.Lock()
			status[cannon.Generation()] = s
			statusMu.Unlock()
		}(c)
	}

	wg.Wait()
	return status
}
