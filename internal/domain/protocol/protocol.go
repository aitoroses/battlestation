package protocol

import (
	"fmt"

	"github.com/aitoroses/battlestation-codetest/internal/domain/target"
)

// Protocol defines the interface for target selection protocols
type Protocol interface {
	Apply(targets []*target.Target) ([]*target.Target, error)
	Name() string
}

// ValidateProtocols checks if the provided protocols are valid and compatible
func ValidateProtocols(protocols []string) error {
	hasClosest := false
	hasFurthest := false

	for _, p := range protocols {
		switch p {
		case "closest-enemies":
			if hasFurthest {
				return fmt.Errorf("incompatible protocols: closest-enemies and furthest-enemies")
			}
			hasClosest = true
		case "furthest-enemies":
			if hasClosest {
				return fmt.Errorf("incompatible protocols: closest-enemies and furthest-enemies")
			}
			hasFurthest = true
		case "assist-allies", "avoid-crossfire", "prioritize-mech", "avoid-mech":
			// These protocols are always compatible
			continue
		default:
			return fmt.Errorf("invalid protocol: %s", p)
		}
	}
	return nil
}

// CreateProtocolChain creates a chain of protocols in the correct order
func CreateProtocolChain(protocols []string) ([]Protocol, error) {
	if err := ValidateProtocols(protocols); err != nil {
		return nil, err
	}

	// Initialize protocol slices by priority
	var (
		validationProtocols []Protocol
		typeProtocols       []Protocol
		positionProtocols   []Protocol
		tacticalProtocols   []Protocol
	)

	// Create protocol instances in the correct order
	for _, p := range protocols {
		switch p {
		// Validation protocols (first)
		case "avoid-mech":
			validationProtocols = append(validationProtocols, NewAvoidMechProtocol())
		case "avoid-crossfire":
			validationProtocols = append(validationProtocols, NewAvoidCrossfireProtocol())

		// Type protocols (second)
		case "prioritize-mech":
			typeProtocols = append(typeProtocols, NewPrioritizeMechProtocol())

		// Position protocols (third)
		case "closest-enemies":
			positionProtocols = append(positionProtocols, NewClosestEnemiesProtocol())
		case "furthest-enemies":
			positionProtocols = append(positionProtocols, NewFurthestEnemiesProtocol())

		// Tactical protocols (fourth)
		case "assist-allies":
			tacticalProtocols = append(tacticalProtocols, NewAssistAlliesProtocol())
		}
	}

	// Combine all protocols in order
	result := make([]Protocol, 0,
		len(validationProtocols)+
			len(typeProtocols)+
			len(positionProtocols)+
			len(tacticalProtocols))

	result = append(result, validationProtocols...)
	result = append(result, typeProtocols...)
	result = append(result, positionProtocols...)
	result = append(result, tacticalProtocols...)

	return result, nil
}

// ApplyProtocolChain applies all protocols in sequence
func ApplyProtocolChain(chain []Protocol, targets []*target.Target) ([]*target.Target, error) {
	var err error
	current := targets

	for _, p := range chain {
		current, err = p.Apply(current)
		if err != nil {
			return nil, fmt.Errorf("protocol %s failed: %w", p.Name(), err)
		}
		if len(current) == 0 {
			return nil, fmt.Errorf("no valid targets after applying protocol %s", p.Name())
		}
	}

	return current, nil
}
