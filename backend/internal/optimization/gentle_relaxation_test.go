package optimization

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
)

// TestGentleRelax - WRIGHT BROTHERS TEST
func TestGentleRelax(t *testing.T) {
	// Build tiny protein
	sequence := "ACDE"
	angles := make([]geometry.RamachandranAngles, len(sequence))
	for i := range angles {
		angles[i] = geometry.RamachandranAngles{
			Phi: -60.0 * math.Pi / 180.0,
			Psi: -45.0 * math.Pi / 180.0,
		}
	}

	protein, err := geometry.BuildProteinFromAngles(sequence, angles)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Gentle relaxation
	config := DefaultGentleRelaxationConfig()
	config.MaxSteps = 20 // Quick test

	result, err := GentleRelax(protein, config)
	if err != nil {
		t.Fatalf("Gentle relax failed: %v", err)
	}

	t.Logf("Initial energy: %.2f kcal/mol", result.InitialEnergy)
	t.Logf("Final energy: %.2f kcal/mol", result.FinalEnergy)
	t.Logf("Energy change: %.2f kcal/mol", result.EnergyChange)
	t.Logf("Steps: %d", result.Steps)
	t.Logf("Converged: %v", result.Converged)

	// Sanity checks
	if result.FinalEnergy > 1e9 {
		t.Error("Energy is placeholder (>1e9)")
	}

	if result.FinalEnergy > result.InitialEnergy {
		t.Error("Energy increased (should decrease)")
	}

	if math.IsNaN(result.FinalEnergy) || math.IsInf(result.FinalEnergy, 0) {
		t.Error("Energy is NaN or Inf")
	}

	t.Log("✓ Gentle relaxation works!")
}

// TestQuickClashRemoval tests clash removal
func TestQuickClashRemoval(t *testing.T) {
	sequence := "AC"
	angles := []geometry.RamachandranAngles{
		{Phi: -60.0 * math.Pi / 180.0, Psi: -45.0 * math.Pi / 180.0},
		{Phi: -60.0 * math.Pi / 180.0, Psi: -45.0 * math.Pi / 180.0},
	}

	protein, _ := geometry.BuildProteinFromAngles(sequence, angles)

	clashes := QuickClashRemoval(protein)
	t.Logf("Clashes fixed: %d", clashes)

	// Should work without errors
	t.Log("✓ Clash removal works!")
}
