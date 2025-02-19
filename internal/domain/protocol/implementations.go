package protocol

import (
	"sort"

	"github.com/aitoroses/battlestation-codetest/internal/domain/target"
)

// AvoidMechProtocol filters out mech targets
type AvoidMechProtocol struct{}

func NewAvoidMechProtocol() *AvoidMechProtocol {
	return &AvoidMechProtocol{}
}

func (p *AvoidMechProtocol) Name() string {
	return "avoid-mech"
}

func (p *AvoidMechProtocol) Apply(targets []*target.Target) ([]*target.Target, error) {
	result := make([]*target.Target, 0, len(targets))
	for _, t := range targets {
		if !t.IsMech() {
			result = append(result, t)
		}
	}
	return result, nil
}

// AvoidCrossfireProtocol filters out targets with allies
type AvoidCrossfireProtocol struct{}

func NewAvoidCrossfireProtocol() *AvoidCrossfireProtocol {
	return &AvoidCrossfireProtocol{}
}

func (p *AvoidCrossfireProtocol) Name() string {
	return "avoid-crossfire"
}

func (p *AvoidCrossfireProtocol) Apply(targets []*target.Target) ([]*target.Target, error) {
	result := make([]*target.Target, 0, len(targets))
	for _, t := range targets {
		if !t.HasAllies() {
			result = append(result, t)
		}
	}
	return result, nil
}

// PrioritizeMechProtocol prioritizes mech targets
type PrioritizeMechProtocol struct{}

func NewPrioritizeMechProtocol() *PrioritizeMechProtocol {
	return &PrioritizeMechProtocol{}
}

func (p *PrioritizeMechProtocol) Name() string {
	return "prioritize-mech"
}

func (p *PrioritizeMechProtocol) Apply(targets []*target.Target) ([]*target.Target, error) {
	// First check if there are any mech targets
	hasMech := false
	for _, t := range targets {
		if t.IsMech() {
			hasMech = true
			break
		}
	}

	// If there are mech targets, filter out non-mech targets
	if hasMech {
		result := make([]*target.Target, 0, len(targets))
		for _, t := range targets {
			if t.IsMech() {
				result = append(result, t)
			}
		}
		return result, nil
	}

	// If no mech targets, return all targets unchanged
	return targets, nil
}

// ClosestEnemiesProtocol selects closest targets
type ClosestEnemiesProtocol struct{}

func NewClosestEnemiesProtocol() *ClosestEnemiesProtocol {
	return &ClosestEnemiesProtocol{}
}

func (p *ClosestEnemiesProtocol) Name() string {
	return "closest-enemies"
}

func (p *ClosestEnemiesProtocol) Apply(targets []*target.Target) ([]*target.Target, error) {
	if len(targets) <= 1 {
		return targets, nil
	}

	// Sort by distance
	sorted := make([]*target.Target, len(targets))
	copy(sorted, targets)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Distance() < sorted[j].Distance()
	})

	// Return only the closest targets (those with same minimum distance)
	minDist := sorted[0].Distance()
	result := make([]*target.Target, 0)
	for _, t := range sorted {
		if t.Distance() == minDist {
			result = append(result, t)
		} else {
			break
		}
	}
	return result, nil
}

// FurthestEnemiesProtocol selects furthest targets
type FurthestEnemiesProtocol struct{}

func NewFurthestEnemiesProtocol() *FurthestEnemiesProtocol {
	return &FurthestEnemiesProtocol{}
}

func (p *FurthestEnemiesProtocol) Name() string {
	return "furthest-enemies"
}

func (p *FurthestEnemiesProtocol) Apply(targets []*target.Target) ([]*target.Target, error) {
	if len(targets) <= 1 {
		return targets, nil
	}

	// Sort by distance in descending order
	sorted := make([]*target.Target, len(targets))
	copy(sorted, targets)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Distance() > sorted[j].Distance()
	})

	// Return only the furthest targets (those with same maximum distance)
	maxDist := sorted[0].Distance()
	result := make([]*target.Target, 0)
	for _, t := range sorted {
		if t.Distance() == maxDist {
			result = append(result, t)
		} else {
			break
		}
	}
	return result, nil
}

// AssistAlliesProtocol prioritizes targets with allies
type AssistAlliesProtocol struct{}

func NewAssistAlliesProtocol() *AssistAlliesProtocol {
	return &AssistAlliesProtocol{}
}

func (p *AssistAlliesProtocol) Name() string {
	return "assist-allies"
}

func (p *AssistAlliesProtocol) Apply(targets []*target.Target) ([]*target.Target, error) {
	// First check if there are any targets with allies
	hasAllies := false
	for _, t := range targets {
		if t.HasAllies() {
			hasAllies = true
			break
		}
	}

	// If there are targets with allies, filter out targets without allies
	if hasAllies {
		result := make([]*target.Target, 0, len(targets))
		for _, t := range targets {
			if t.HasAllies() {
				result = append(result, t)
			}
		}
		return result, nil
	}

	// If no targets with allies, return all targets unchanged
	return targets, nil
}
