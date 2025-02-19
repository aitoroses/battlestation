package protocol

import (
	"reflect"
	"testing"

	"github.com/aitoroses/battlestation-codetest/internal/domain/target"
)

func TestValidateProtocols(t *testing.T) {
	tests := []struct {
		name      string
		protocols []string
		wantErr   bool
	}{
		{
			name:      "valid single protocol",
			protocols: []string{"avoid-mech"},
			wantErr:   false,
		},
		{
			name:      "valid multiple protocols",
			protocols: []string{"avoid-mech", "assist-allies"},
			wantErr:   false,
		},
		{
			name:      "invalid protocol",
			protocols: []string{"invalid-protocol"},
			wantErr:   true,
		},
		{
			name:      "incompatible protocols",
			protocols: []string{"closest-enemies", "furthest-enemies"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProtocols(tt.protocols)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProtocols() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createTestTargets() []*target.Target {
	allies := 5
	return []*target.Target{
		target.NewTarget(
			target.Position{X: 0, Y: 10},
			target.EnemyGroup{Type: target.EnemyTypeMech, Number: 1},
			nil,
		),
		target.NewTarget(
			target.Position{X: 0, Y: 20},
			target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 10},
			nil,
		),
		target.NewTarget(
			target.Position{X: 0, Y: 30},
			target.EnemyGroup{Type: target.EnemyTypeSoldier, Number: 20},
			&allies,
		),
	}
}

func TestAvoidMechProtocol(t *testing.T) {
	p := NewAvoidMechProtocol()
	targets := createTestTargets()

	result, err := p.Apply(targets)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 targets, got %d", len(result))
	}

	for _, r := range result {
		if r.IsMech() {
			t.Errorf("Found mech target after avoid-mech protocol: %+v", r)
		}
	}
}

func TestPrioritizeMechProtocol(t *testing.T) {
	p := NewPrioritizeMechProtocol()
	targets := createTestTargets()

	result, err := p.Apply(targets)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 target, got %d", len(result))
	}

	if !result[0].IsMech() {
		t.Error("First target should be mech type")
	}
}

func TestClosestEnemiesProtocol(t *testing.T) {
	p := NewClosestEnemiesProtocol()
	targets := createTestTargets()

	result, err := p.Apply(targets)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Expected at least one target")
	}

	if result[0].Distance() != 10 {
		t.Errorf("Expected closest target at distance 10, got %f", result[0].Distance())
	}
}

func TestFurthestEnemiesProtocol(t *testing.T) {
	p := NewFurthestEnemiesProtocol()
	targets := createTestTargets()

	result, err := p.Apply(targets)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Expected at least one target")
	}

	if result[0].Distance() != 30 {
		t.Errorf("Expected furthest target at distance 30, got %f", result[0].Distance())
	}
}

func TestAssistAlliesProtocol(t *testing.T) {
	p := NewAssistAlliesProtocol()
	targets := createTestTargets()

	result, err := p.Apply(targets)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 target with allies, got %d", len(result))
	}

	if !result[0].HasAllies() {
		t.Error("Selected target should have allies")
	}
}

func TestAvoidCrossfireProtocol(t *testing.T) {
	p := NewAvoidCrossfireProtocol()
	targets := createTestTargets()

	result, err := p.Apply(targets)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 targets without allies, got %d", len(result))
	}

	for _, r := range result {
		if r.HasAllies() {
			t.Error("Selected target should not have allies")
		}
	}
}

func TestCreateProtocolChain(t *testing.T) {
	tests := []struct {
		name      string
		protocols []string
		want      []string
		wantErr   bool
	}{
		{
			name:      "valid chain",
			protocols: []string{"avoid-mech", "closest-enemies", "assist-allies"},
			want:      []string{"avoid-mech", "closest-enemies", "assist-allies"},
			wantErr:   false,
		},
		{
			name:      "invalid protocol",
			protocols: []string{"invalid-protocol"},
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateProtocolChain(tt.protocols)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProtocolChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				var gotNames []string
				for _, p := range got {
					gotNames = append(gotNames, p.Name())
				}
				if !reflect.DeepEqual(gotNames, tt.want) {
					t.Errorf("CreateProtocolChain() = %v, want %v", gotNames, tt.want)
				}
			}
		})
	}
}
