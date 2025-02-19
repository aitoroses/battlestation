package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aitoroses/battlestation-codetest/internal/domain/attack"
	"github.com/aitoroses/battlestation-codetest/internal/platform/metrics"
)

// Handler handles HTTP requests for the battle station
type Handler struct {
	coordinator *attack.Coordinator
	logger      *slog.Logger
}

// NewHandler creates a new HTTP handler
func NewHandler(coordinator *attack.Coordinator, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
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

	// Start request timing
	start := time.Now()

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
	resp, err := h.coordinator.ProcessAttack(ctx, &req)
	duration := time.Since(start)

	// Record metrics
	for _, protocol := range req.Protocols {
		metrics.RecordRequestDuration(protocol, duration.Seconds())
		if err != nil {
			metrics.RecordRequestComplete(protocol, "error")
		} else {
			metrics.RecordRequestComplete(protocol, "success")
		}
	}

	// Log request details
	h.logger.Info("Attack request processed",
		slog.Duration("duration", duration),
		slog.Any("protocols", req.Protocols),
		slog.Int("targets", len(req.Scan)),
		slog.Any("error", err),
	)

	if err != nil {
		// Determine appropriate status code based on error
		statusCode := h.determineStatusCode(err)
		h.writeError(w, err, statusCode)
		return
	}

	// Write successful response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to write response",
			slog.String("error", err.Error()),
		)
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

	h.logger.Error("Request error",
		slog.String("error", err.Error()),
		slog.Int("status_code", statusCode),
	)

	metrics.RecordError("http", err.Error())

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to write error response",
			slog.String("error", err.Error()),
		)
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
