package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aitoroses/battlestation-codetest/internal/domain/attack"
)

// Handler handles HTTP requests for the battle station
type Handler struct {
	coordinator *attack.Coordinator
	logger      *log.Logger
}

// NewHandler creates a new HTTP handler
func NewHandler(coordinator *attack.Coordinator, logger *log.Logger) *Handler {
	return &Handler{
		coordinator: coordinator,
		logger:      logger,
	}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /attack", h.handleAttack)
}

// handleAttack processes attack requests
func (h *Handler) handleAttack(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, fmt.Errorf("failed to read request body: %w", err), http.StatusBadRequest)
		return
	}

	// Parse request
	var req attack.Request
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(w, fmt.Errorf("failed to parse request: %w", err), http.StatusBadRequest)
		return
	}

	// Validate request
	if err := attack.ValidateRequest(&req); err != nil {
		h.writeError(w, fmt.Errorf("invalid request: %w", err), http.StatusBadRequest)
		return
	}

	// Create context with timeout
	ctx := r.Context()

	// Process attack
	start := time.Now()
	resp, err := h.coordinator.ProcessAttack(ctx, &req)
	duration := time.Since(start)

	// Log request details
	h.logger.Printf("Attack request processed in %v - Protocols: %v, Targets: %d, Error: %v",
		duration, req.Protocols, len(req.Scan), err)

	if err != nil {
		// Determine appropriate status code based on error
		statusCode := h.determineStatusCode(err)
		h.writeError(w, err, statusCode)
		return
	}

	// Write successful response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Printf("Failed to write response: %v", err)
	}
}

// writeError writes an error response in JSON format
func (h *Handler) writeError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	response := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Failed to write error response: %v", err)
	}
}

// determineStatusCode maps errors to appropriate HTTP status codes
func (h *Handler) determineStatusCode(err error) int {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout
	case errors.Is(err, context.Canceled):
		return http.StatusServiceUnavailable
	default:
		// Check error string for specific cases
		switch err.Error() {
		case "no cannons available":
			return http.StatusServiceUnavailable
		case "no valid targets in range":
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}
}
