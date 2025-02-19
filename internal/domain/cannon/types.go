package cannon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aitoroses/battlestation-codetest/internal/domain/target"
)

// Generation represents an ion cannon generation
type Generation int

const (
	Generation1 Generation = 1
	Generation2 Generation = 2
	Generation3 Generation = 3
)

// FireTime returns the fire time in seconds for a given generation
func (g Generation) FireTime() float64 {
	switch g {
	case Generation1:
		return 3.5
	case Generation2:
		return 1.5
	case Generation3:
		return 2.5
	default:
		return 0
	}
}

// Status represents the current status of an ion cannon
type Status struct {
	Generation int  `json:"generation"`
	Available  bool `json:"available"`
}

// FireRequest represents a request to fire an ion cannon
type FireRequest struct {
	Target  target.Position `json:"target"`
	Enemies int             `json:"enemies"`
}

// FireResponse represents the response from firing an ion cannon
type FireResponse struct {
	Casualties int `json:"casualties"`
	Generation int `json:"generation"`
}

// IonCannon represents a single ion cannon
type IonCannon struct {
	generation  Generation
	baseURL     string
	lastFired   time.Time
	mu          sync.RWMutex
	httpClient  HTTPClient
	statusCache *StatusCache
}

// NewIonCannon creates a new ion cannon instance
func NewIonCannon(generation Generation, baseURL string, client HTTPClient) *IonCannon {
	return &IonCannon{
		generation:  generation,
		baseURL:     baseURL,
		httpClient:  client,
		statusCache: NewStatusCache(100 * time.Millisecond),
	}
}

// Generation returns the cannon's generation
func (c *IonCannon) Generation() Generation {
	return c.generation
}

// IsAvailable checks if the cannon is available based on its fire time
func (c *IonCannon) IsAvailable() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastFired.IsZero() {
		return true
	}

	return time.Since(c.lastFired) >= time.Duration(c.generation.FireTime()*float64(time.Second))
}

// CheckStatus checks the cannon's status via HTTP
func (c *IonCannon) CheckStatus(ctx context.Context) (*Status, error) {
	// Check cache first
	if status := c.statusCache.Get(); status != nil {
		return status, nil
	}

	// Make HTTP request
	status, err := c.httpClient.GetStatus(ctx, c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get cannon status: %w", err)
	}

	// Update cache
	c.statusCache.Set(status)
	return status, nil
}

// Fire sends a fire request to the cannon
func (c *IonCannon) Fire(ctx context.Context, req *FireRequest) (*FireResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double check availability
	if !c.IsAvailable() {
		return nil, fmt.Errorf("cannon generation %d is not available", c.generation)
	}

	// Send fire request
	resp, err := c.httpClient.Fire(ctx, c.baseURL, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fire cannon: %w", err)
	}

	// Update last fired time
	c.lastFired = time.Now()
	return resp, nil
}

// HTTPClient defines the interface for making HTTP requests to ion cannons
type HTTPClient interface {
	GetStatus(ctx context.Context, baseURL string) (*Status, error)
	Fire(ctx context.Context, baseURL string, req *FireRequest) (*FireResponse, error)
}

// StatusCache implements a simple time-based cache for cannon status
type StatusCache struct {
	status    *Status
	timestamp time.Time
	ttl       time.Duration
	mu        sync.RWMutex
}

// NewStatusCache creates a new status cache with the given TTL
func NewStatusCache(ttl time.Duration) *StatusCache {
	return &StatusCache{
		ttl: ttl,
	}
}

// Get returns the cached status if it's still valid, nil otherwise
func (c *StatusCache) Get() *Status {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.status == nil || time.Since(c.timestamp) > c.ttl {
		return nil
	}
	return c.status
}

// Set updates the cached status and timestamp
func (c *StatusCache) Set(status *Status) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.status = status
	c.timestamp = time.Now()
}
