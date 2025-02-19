package tests

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aitoroses/battlestation-codetest/internal/domain/attack"
	"github.com/aitoroses/battlestation-codetest/internal/domain/cannon"
	httpPlatform "github.com/aitoroses/battlestation-codetest/internal/platform/http"
)

// TestCase represents a test case from test_cases.txt
type TestCase struct {
	Input  string
	Output string
}

func loadTestCases(t *testing.T) []TestCase {
	file, err := os.Open("../test_cases.txt")
	if err != nil {
		t.Fatalf("Failed to open test_cases.txt: %v", err)
	}
	defer file.Close()

	var testCases []TestCase
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			t.Fatalf("Invalid test case format: %s", line)
		}
		testCases = append(testCases, TestCase{
			Input:  strings.TrimSpace(parts[0]),
			Output: strings.TrimSpace(parts[1]),
		})
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading test_cases.txt: %v", err)
	}

	return testCases
}

// MockCannonManager implements attack.CannonManager for testing
type MockCannonManager struct {
	expectedOutput string
}

func NewMockCannonManager(expectedOutput string) *MockCannonManager {
	return &MockCannonManager{
		expectedOutput: expectedOutput,
	}
}

func (m *MockCannonManager) GetBestAvailable(ctx context.Context) (*cannon.IonCannon, error) {
	return &cannon.IonCannon{}, nil
}

func (m *MockCannonManager) Fire(ctx context.Context, c *cannon.IonCannon, req *cannon.FireRequest) (*cannon.FireResponse, error) {
	// Parse expected output to determine casualties and generation
	var expected attack.Response
	if err := json.Unmarshal([]byte(m.expectedOutput), &expected); err != nil {
		return nil, fmt.Errorf("failed to parse expected output: %w", err)
	}

	return &cannon.FireResponse{
		Casualties: expected.Casualties,
		Generation: expected.Generation,
	}, nil
}

func TestIntegrationTestCases(t *testing.T) {
	testCases := loadTestCases(t)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase_%d", i+1), func(t *testing.T) {
			// Create mock manager with expected output
			mockManager := NewMockCannonManager(tc.Output)

			// Create coordinator with mock manager
			coordinator := attack.NewCoordinator(mockManager)

			// Create handler
			handler := httpPlatform.NewHandler(coordinator, nil)

			// Create server mux and register routes
			mux := http.NewServeMux()
			handler.RegisterRoutes(mux)

			// Create test server
			server := httptest.NewServer(mux)
			defer server.Close()

			// Send request
			resp, err := http.Post(server.URL+"/attack", "application/json", bytes.NewBufferString(tc.Input))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
			}

			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			// Compare responses
			var got, want attack.Response
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}
			if err := json.Unmarshal([]byte(tc.Output), &want); err != nil {
				t.Fatalf("Failed to parse expected output: %v", err)
			}

			if got.Target != want.Target ||
				got.Casualties != want.Casualties ||
				got.Generation != want.Generation {
				t.Errorf("Case %d failed:\nGot:  %+v\nWant: %+v", i+1, got, want)
			}
		})
	}
}
