package cannon

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aitoroses/battlestation-codetest/internal/domain/target"
)

// MockHTTPClient implements HTTPClient for testing
type MockHTTPClient struct {
	statusResponses map[string]*Status
	fireResponses   map[string]*FireResponse
	statusError     error
	fireError       error
}

func (m *MockHTTPClient) GetStatus(ctx context.Context, baseURL string) (*Status, error) {
	if m.statusError != nil {
		return nil, m.statusError
	}
	if status, ok := m.statusResponses[baseURL]; ok {
		return status, nil
	}
	return nil, errors.New("unexpected baseURL")
}

func (m *MockHTTPClient) Fire(ctx context.Context, baseURL string, req *FireRequest) (*FireResponse, error) {
	if m.fireError != nil {
		return nil, m.fireError
	}
	if resp, ok := m.fireResponses[baseURL]; ok {
		return resp, nil
	}
	return nil, errors.New("unexpected baseURL")
}

func TestManager_GetBestAvailable(t *testing.T) {
	tests := []struct {
		name            string
		cannons         []*IonCannon
		statusResponses map[string]*Status
		statusError     error
		wantGeneration  Generation
		wantErr         bool
	}{
		{
			name: "all cannons available",
			cannons: []*IonCannon{
				NewIonCannon(Generation1, "http://cannon1", nil),
				NewIonCannon(Generation2, "http://cannon2", nil),
				NewIonCannon(Generation3, "http://cannon3", nil),
			},
			statusResponses: map[string]*Status{
				"http://cannon1": {Generation: 1, Available: true},
				"http://cannon2": {Generation: 2, Available: true},
				"http://cannon3": {Generation: 3, Available: true},
			},
			wantGeneration: Generation1,
			wantErr:        false,
		},
		{
			name: "first generation unavailable",
			cannons: []*IonCannon{
				NewIonCannon(Generation1, "http://cannon1", nil),
				NewIonCannon(Generation2, "http://cannon2", nil),
				NewIonCannon(Generation3, "http://cannon3", nil),
			},
			statusResponses: map[string]*Status{
				"http://cannon1": {Generation: 1, Available: false},
				"http://cannon2": {Generation: 2, Available: true},
				"http://cannon3": {Generation: 3, Available: true},
			},
			wantGeneration: Generation2,
			wantErr:        false,
		},
		{
			name: "all cannons unavailable",
			cannons: []*IonCannon{
				NewIonCannon(Generation1, "http://cannon1", nil),
				NewIonCannon(Generation2, "http://cannon2", nil),
				NewIonCannon(Generation3, "http://cannon3", nil),
			},
			statusResponses: map[string]*Status{
				"http://cannon1": {Generation: 1, Available: false},
				"http://cannon2": {Generation: 2, Available: false},
				"http://cannon3": {Generation: 3, Available: false},
			},
			wantErr: true,
		},
		{
			name: "status check error",
			cannons: []*IonCannon{
				NewIonCannon(Generation1, "http://cannon1", nil),
			},
			statusError: errors.New("network error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockHTTPClient{
				statusResponses: tt.statusResponses,
				statusError:     tt.statusError,
			}

			// Set mock client for each cannon
			for _, c := range tt.cannons {
				c.httpClient = mockClient
			}

			manager := NewManager(tt.cannons)
			got, err := manager.GetBestAvailable(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetBestAvailable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.Generation() != tt.wantGeneration {
				t.Errorf("Manager.GetBestAvailable() = generation %v, want %v", got.Generation(), tt.wantGeneration)
			}
		})
	}
}

func TestManager_Fire(t *testing.T) {
	tests := []struct {
		name          string
		cannon        *IonCannon
		fireResponses map[string]*FireResponse
		fireError     error
		wantErr       bool
	}{
		{
			name:   "successful fire",
			cannon: NewIonCannon(Generation1, "http://cannon1", nil),
			fireResponses: map[string]*FireResponse{
				"http://cannon1": {Casualties: 10, Generation: 1},
			},
			wantErr: false,
		},
		{
			name:      "fire error",
			cannon:    NewIonCannon(Generation1, "http://cannon1", nil),
			fireError: errors.New("network error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockHTTPClient{
				fireResponses: tt.fireResponses,
				fireError:     tt.fireError,
			}

			// Set mock client for cannon
			tt.cannon.httpClient = mockClient

			manager := NewManager([]*IonCannon{tt.cannon})
			req := &FireRequest{
				Target:  target.Position{X: 0, Y: 10},
				Enemies: 10,
			}

			_, err := manager.Fire(context.Background(), tt.cannon, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Fire() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIonCannon_IsAvailable(t *testing.T) {
	tests := []struct {
		name       string
		generation Generation
		lastFired  time.Time
		want       bool
	}{
		{
			name:       "never fired",
			generation: Generation1,
			lastFired:  time.Time{},
			want:       true,
		},
		{
			name:       "recently fired",
			generation: Generation1,
			lastFired:  time.Now(),
			want:       false,
		},
		{
			name:       "fire time elapsed",
			generation: Generation1,
			lastFired:  time.Now().Add(-4 * time.Second), // Generation1 has 3.5s fire time
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cannon := NewIonCannon(tt.generation, "http://test", nil)
			cannon.lastFired = tt.lastFired

			if got := cannon.IsAvailable(); got != tt.want {
				t.Errorf("IonCannon.IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
