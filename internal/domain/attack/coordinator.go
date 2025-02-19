package attack

import (
	"context"
	"fmt"

	"github.com/aitoroses/battlestation-codetest/internal/domain/cannon"
	"github.com/aitoroses/battlestation-codetest/internal/domain/protocol"
	"github.com/aitoroses/battlestation-codetest/internal/domain/target"
)

// Request represents an attack request
type Request struct {
	Protocols []string    `json:"protocols"`
	Scan      []ScanPoint `json:"scan"`
}

// ScanPoint represents a single point in the scan data
type ScanPoint struct {
	Coordinates target.Position   `json:"coordinates"`
	Enemies     target.EnemyGroup `json:"enemies"`
	Allies      *int              `json:"allies,omitempty"`
}

// Response represents the attack response
type Response struct {
	Target     target.Position `json:"target"`
	Casualties int             `json:"casualties"`
	Generation int             `json:"generation"`
}

// CannonManager defines the interface for managing ion cannons
type CannonManager interface {
	GetBestAvailable(ctx context.Context) (*cannon.IonCannon, error)
	Fire(ctx context.Context, c *cannon.IonCannon, req *cannon.FireRequest) (*cannon.FireResponse, error)
}

// Coordinator orchestrates the attack process
type Coordinator struct {
	cannonManager CannonManager
}

// NewCoordinator creates a new attack coordinator
func NewCoordinator(cannonManager CannonManager) *Coordinator {
	return &Coordinator{
		cannonManager: cannonManager,
	}
}

// ProcessAttack handles the complete attack sequence
func (c *Coordinator) ProcessAttack(ctx context.Context, req *Request) (*Response, error) {
	// 1. Create protocol chain
	chain, err := protocol.CreateProtocolChain(req.Protocols)
	if err != nil {
		return nil, fmt.Errorf("invalid protocols: %w", err)
	}

	// 2. Convert scan points to targets
	targets := make([]*target.Target, 0, len(req.Scan))
	for _, point := range req.Scan {
		t := target.NewTarget(point.Coordinates, point.Enemies, point.Allies)
		if t.IsValid() {
			targets = append(targets, t)
		}
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no valid targets in range")
	}

	// 3. Apply protocol chain to select target
	selectedTargets, err := protocol.ApplyProtocolChain(chain, targets)
	if err != nil {
		return nil, fmt.Errorf("target selection failed: %w", err)
	}

	// Always select first target after protocol application
	selectedTarget := selectedTargets[0]

	// 4. Get best available cannon
	selectedCannon, err := c.cannonManager.GetBestAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("no cannon available: %w", err)
	}

	// 5. Fire cannon at target
	fireReq := &cannon.FireRequest{
		Target:  selectedTarget.Coordinates,
		Enemies: selectedTarget.Enemies.Number,
	}

	fireResp, err := c.cannonManager.Fire(ctx, selectedCannon, fireReq)
	if err != nil {
		return nil, fmt.Errorf("cannon fire failed: %w", err)
	}

	// 6. Prepare response
	return &Response{
		Target:     selectedTarget.Coordinates,
		Casualties: fireResp.Casualties,
		Generation: fireResp.Generation,
	}, nil
}

// ValidateRequest checks if the attack request is valid
func ValidateRequest(req *Request) error {
	if len(req.Protocols) == 0 {
		return fmt.Errorf("no protocols specified")
	}

	if len(req.Scan) == 0 {
		return fmt.Errorf("no scan points provided")
	}

	// Validate each scan point
	for i, point := range req.Scan {
		if err := validateScanPoint(point); err != nil {
			return fmt.Errorf("invalid scan point at index %d: %w", i, err)
		}
	}

	return nil
}

// validateScanPoint checks if a scan point is valid
func validateScanPoint(point ScanPoint) error {
	// Validate enemy type
	switch point.Enemies.Type {
	case target.EnemyTypeSoldier, target.EnemyTypeMech:
		// Valid types
	default:
		return fmt.Errorf("invalid enemy type: %s", point.Enemies.Type)
	}

	// Validate enemy number
	if point.Enemies.Number <= 0 {
		return fmt.Errorf("invalid enemy number: %d", point.Enemies.Number)
	}

	// Validate allies if present
	if point.Allies != nil && *point.Allies < 0 {
		return fmt.Errorf("invalid allies number: %d", *point.Allies)
	}

	return nil
}
