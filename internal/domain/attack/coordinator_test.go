package attack

import (
	"context"
	"testing"

	"github.com/aitoroses/battlestation-codetest/internal/domain/cannon"
	"github.com/aitoroses/battlestation-codetest/internal/domain/target"
)

// MockCannonManager implements a test double for the cannon manager
type MockCannonManager struct {
	bestCannon *cannon.IonCannon
	bestErr    error
	fireResp   *cannon.FireResponse
	fireErr    error
}

func (m *MockCannonManager) GetBestAvailable(ctx context.Context) (*cannon.IonCannon, error) {
	return m.bestCannon, m.bestErr
}

func (m *MockCannonManager) Fire(ctx context.Context, c *cannon.IonCannon, req *cannon.FireRequest) (*cannon.FireResponse, error) {
	return m.fireResp, m.fireErr
}

func TestCoordinator_ProcessAttack(t *testing.T) {
	tests := []struct {
		name         string
		request      *Request
		mockCannon   *cannon.IonCannon
		mockFireResp *cannon.FireResponse
		wantErr      bool
	}{
		{
			name: "successful attack - avoid mech",
			request: &Request{
				Protocols: []string{"avoid-mech"},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 40},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 10},
					},
					{
						Coordinates: target.Position{X: 0, Y: 80},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeMech, Number: 1},
					},
				},
			},
			mockCannon: &cannon.IonCannon{},
			mockFireResp: &cannon.FireResponse{
				Casualties: 10,
				Generation: 1,
			},
			wantErr: false,
		},
		{
			name: "successful attack - prioritize mech",
			request: &Request{
				Protocols: []string{"prioritize-mech"},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 40},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 10},
					},
					{
						Coordinates: target.Position{X: 0, Y: 80},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeMech, Number: 1},
					},
				},
			},
			mockCannon: &cannon.IonCannon{},
			mockFireResp: &cannon.FireResponse{
				Casualties: 1,
				Generation: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid protocols",
			request: &Request{
				Protocols: []string{"invalid-protocol"},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 40},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 10},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no valid targets",
			request: &Request{
				Protocols: []string{"avoid-mech"},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 150}, // Beyond range
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 10},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockManager := &MockCannonManager{
				bestCannon: tt.mockCannon,
				fireResp:   tt.mockFireResp,
			}

			coordinator := NewCoordinator(mockManager)
			resp, err := coordinator.ProcessAttack(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Coordinator.ProcessAttack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Error("Expected response, got nil")
					return
				}

				if resp.Casualties != tt.mockFireResp.Casualties {
					t.Errorf("Expected casualties %d, got %d", tt.mockFireResp.Casualties, resp.Casualties)
				}

				if resp.Generation != tt.mockFireResp.Generation {
					t.Errorf("Expected generation %d, got %d", tt.mockFireResp.Generation, resp.Generation)
				}
			}
		})
	}
}

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *Request
		wantErr bool
	}{
		{
			name: "valid request",
			request: &Request{
				Protocols: []string{"avoid-mech"},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 40},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 10},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty protocols",
			request: &Request{
				Protocols: []string{},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 40},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 10},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty scan",
			request: &Request{
				Protocols: []string{"avoid-mech"},
				Scan:      []ScanPoint{},
			},
			wantErr: true,
		},
		{
			name: "invalid enemy type",
			request: &Request{
				Protocols: []string{"avoid-mech"},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 40},
						Enemies:     target.EnemyGroup{Type: "invalid", Number: 10},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid enemy number",
			request: &Request{
				Protocols: []string{"avoid-mech"},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 40},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 0},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid allies number",
			request: &Request{
				Protocols: []string{"avoid-mech"},
				Scan: []ScanPoint{
					{
						Coordinates: target.Position{X: 0, Y: 40},
						Enemies:     target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 10},
						Allies:      func() *int { n := -1; return &n }(),
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
