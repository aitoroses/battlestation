package target

import "math"

// Position represents x,y coordinates
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Distance calculates the distance from origin (0,0)
func (p Position) Distance() float64 {
	return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}

// EnemyType represents the type of enemy (soldier or mech)
type EnemyType string

const (
	EnemyTypeSoldier EnemyType = "soldier"
	EnemyTypeMech    EnemyType = "mech"
)

// EnemyGroup represents a group of enemies of the same type
type EnemyGroup struct {
	Type   EnemyType `json:"type"`
	Number int       `json:"number"`
}

// Target represents a potential target with its position and enemy information
type Target struct {
	Coordinates Position   `json:"coordinates"`
	Enemies     EnemyGroup `json:"enemies"`
	Allies      *int       `json:"allies,omitempty"`
	distance    float64    // cached distance value
}

// NewTarget creates a new Target and pre-calculates its distance
func NewTarget(coords Position, enemies EnemyGroup, allies *int) *Target {
	t := &Target{
		Coordinates: coords,
		Enemies:     enemies,
		Allies:      allies,
	}
	t.distance = coords.Distance()
	return t
}

// Distance returns the pre-calculated distance from origin
func (t *Target) Distance() float64 {
	return t.distance
}

// IsValid checks if the target is within valid range (<=100km)
func (t *Target) IsValid() bool {
	return t.distance <= 100
}

// HasAllies returns true if there are allies present at this target
func (t *Target) HasAllies() bool {
	return t.Allies != nil && *t.Allies > 0
}

// IsMech returns true if the enemy type is mech
func (t *Target) IsMech() bool {
	return t.Enemies.Type == EnemyTypeMech
}

// ScanData represents the complete scan information from probe droids
type ScanData struct {
	Protocols []string `json:"protocols"`
	Scan      []struct {
		Coordinates Position   `json:"coordinates"`
		Enemies     EnemyGroup `json:"enemies"`
		Allies      *int       `json:"allies,omitempty"`
	} `json:"scan"`
}
