package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aitoroses/battlestation-codetest/internal/domain/attack"
	"github.com/aitoroses/battlestation-codetest/internal/domain/cannon"
)

// MockCannonManager implements attack.CannonManager for testing
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

func TestHandler_HandleAttack(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockCannon   *cannon.IonCannon
		mockFireResp *cannon.FireResponse
		wantStatus   int
		wantResponse string
	}{
		{
			name: "avoid-mech protocol",
			requestBody: `{
				"protocols": ["avoid-mech"],
				"scan": [
					{
						"coordinates": {"x": 0, "y": 40},
						"enemies": {"type": "soldier", "number": 10}
					},
					{
						"coordinates": {"x": 0, "y": 80},
						"allies": 5,
						"enemies": {"type": "mech", "number": 1}
					}
				]
			}`,
			mockCannon: &cannon.IonCannon{},
			mockFireResp: &cannon.FireResponse{
				Casualties: 10,
				Generation: 1,
			},
			wantStatus: http.StatusOK,
			wantResponse: `{
				"target": {"x": 0, "y": 40},
				"casualties": 10,
				"generation": 1
			}`,
		},
		{
			name: "prioritize-mech protocol",
			requestBody: `{
				"protocols": ["prioritize-mech"],
				"scan": [
					{
						"coordinates": {"x": 0, "y": 40},
						"enemies": {"type": "soldier", "number": 10}
					},
					{
						"coordinates": {"x": 0, "y": 80},
						"allies": 5,
						"enemies": {"type": "mech", "number": 1}
					}
				]
			}`,
			mockCannon: &cannon.IonCannon{},
			mockFireResp: &cannon.FireResponse{
				Casualties: 1,
				Generation: 2,
			},
			wantStatus: http.StatusOK,
			wantResponse: `{
				"target": {"x": 0, "y": 80},
				"casualties": 1,
				"generation": 2
			}`,
		},
		{
			name: "invalid protocol",
			requestBody: `{
				"protocols": ["invalid-protocol"],
				"scan": [
					{
						"coordinates": {"x": 0, "y": 40},
						"enemies": {"type": "soldier", "number": 10}
					}
				]
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid request format",
			requestBody: `invalid json`,
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock manager
			mockManager := &MockCannonManager{
				bestCannon: tt.mockCannon,
				fireResp:   tt.mockFireResp,
			}

			// Create coordinator with mock manager
			coordinator := attack.NewCoordinator(mockManager)

			// Create handler
			handler := NewHandler(coordinator, nil)

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.handleAttack(w, r)
			}))
			defer server.Close()

			// Create request
			req, err := http.NewRequest(http.MethodPost, server.URL+"/attack", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Send request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", resp.StatusCode, tt.wantStatus)
			}

			// For successful requests, verify response
			if tt.wantStatus == http.StatusOK {
				var got, want attack.Response
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if err := json.Unmarshal([]byte(tt.wantResponse), &want); err != nil {
					t.Fatalf("Failed to decode expected response: %v", err)
				}

				if got.Target != want.Target ||
					got.Casualties != want.Casualties ||
					got.Generation != want.Generation {
					t.Errorf("Handler returned unexpected response: got %+v want %+v", got, want)
				}
			}
		})
	}
}

func TestHandler_RegisterRoutes(t *testing.T) {
	handler := NewHandler(nil, nil)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Test that the attack endpoint is registered
	server := httptest.NewServer(mux)
	defer server.Close()

	// Send request to non-existent endpoint
	resp, err := http.Get(server.URL + "/nonexistent")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for non-existent endpoint, got %d", resp.StatusCode)
	}

	// Send request with wrong method
	resp, err = http.Get(server.URL + "/attack")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 for GET /attack, got %d", resp.StatusCode)
	}
}
