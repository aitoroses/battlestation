package target

import (
	"math"
	"testing"
)

func TestPosition_Distance(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     float64
	}{
		{
			name:     "origin",
			position: Position{X: 0, Y: 0},
			want:     0,
		},
		{
			name:     "positive coordinates",
			position: Position{X: 3, Y: 4},
			want:     5,
		},
		{
			name:     "negative coordinates",
			position: Position{X: -3, Y: -4},
			want:     5,
		},
		{
			name:     "mixed coordinates",
			position: Position{X: -3, Y: 4},
			want:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.position.Distance()
			if math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("Position.Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTarget_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		target *Target
		want   bool
	}{
		{
			name: "within range",
			target: NewTarget(
				Position{X: 60, Y: 80},
				EnemyGroup{Type: EnemyTypeSoldier, Number: 10},
				nil,
			),
			want: true,
		},
		{
			name: "at range limit",
			target: NewTarget(
				Position{X: 80, Y: 60},
				EnemyGroup{Type: EnemyTypeSoldier, Number: 10},
				nil,
			),
			want: true,
		},
		{
			name: "out of range",
			target: NewTarget(
				Position{X: 80, Y: 80},
				EnemyGroup{Type: EnemyTypeSoldier, Number: 10},
				nil,
			),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.IsValid(); got != tt.want {
				t.Errorf("Target.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTarget_HasAllies(t *testing.T) {
	allies := 5
	tests := []struct {
		name   string
		target *Target
		want   bool
	}{
		{
			name: "with allies",
			target: NewTarget(
				Position{X: 0, Y: 0},
				EnemyGroup{Type: EnemyTypeSoldier, Number: 10},
				&allies,
			),
			want: true,
		},
		{
			name: "without allies",
			target: NewTarget(
				Position{X: 0, Y: 0},
				EnemyGroup{Type: EnemyTypeSoldier, Number: 10},
				nil,
			),
			want: false,
		},
		{
			name: "zero allies",
			target: NewTarget(
				Position{X: 0, Y: 0},
				EnemyGroup{Type: EnemyTypeSoldier, Number: 10},
				new(int), // zero value
			),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.HasAllies(); got != tt.want {
				t.Errorf("Target.HasAllies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTarget_IsMech(t *testing.T) {
	tests := []struct {
		name   string
		target *Target
		want   bool
	}{
		{
			name: "is mech",
			target: NewTarget(
				Position{X: 0, Y: 0},
				EnemyGroup{Type: EnemyTypeMech, Number: 1},
				nil,
			),
			want: true,
		},
		{
			name: "is soldier",
			target: NewTarget(
				Position{X: 0, Y: 0},
				EnemyGroup{Type: EnemyTypeSoldier, Number: 10},
				nil,
			),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.IsMech(); got != tt.want {
				t.Errorf("Target.IsMech() = %v, want %v", got, tt.want)
			}
		})
	}
}
