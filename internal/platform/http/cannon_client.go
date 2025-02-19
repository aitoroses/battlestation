package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aitoroses/battlestation-codetest/internal/domain/cannon"
)

// CannonClient implements the cannon.HTTPClient interface
type CannonClient struct {
	client  *http.Client
	timeout time.Duration
}

// NewCannonClient creates a new HTTP client for ion cannons
func NewCannonClient(timeout time.Duration) *CannonClient {
	return &CannonClient{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// GetStatus retrieves the current status of an ion cannon
func (c *CannonClient) GetStatus(ctx context.Context, baseURL string) (*cannon.Status, error) {
	// Create request with context
	statusURL, err := url.JoinPath(baseURL, "status")
	if err != nil {
		return nil, fmt.Errorf("failed to create status URL: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		statusURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	// Parse response
	var status cannon.Status
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &status, nil
}

// Fire sends a fire request to an ion cannon
func (c *CannonClient) Fire(ctx context.Context, baseURL string, req *cannon.FireRequest) (*cannon.FireResponse, error) {
	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request with context
	fireURL, err := url.JoinPath(baseURL, "fire")
	if err != nil {
		return nil, fmt.Errorf("failed to create fire URL: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fireURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, respBody)
	}

	// Parse response
	var fireResp cannon.FireResponse
	if err := json.Unmarshal(respBody, &fireResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &fireResp, nil
}
