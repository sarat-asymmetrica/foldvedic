package physics

import (
	"math"
	"testing"
)

// TestAngleDiff tests the periodic boundary handling for angular differences
func TestAngleDiff(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"Same angle", 0, 0, 0},
		{"Simple difference", 45, 30, 15},
		{"Negative difference", 30, 45, -15},
		{"Wraparound positive", 170, -170, -20},
		{"Wraparound negative", -170, 170, 20},
		{"Full circle", 0, 360, 0},
		{"Half circle", 0, 180, 180},
		{"Quarter circle", 0, 90, -90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := angleDiff(tt.a, tt.b)
			if math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("angleDiff(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestGaussianPotential tests the 2D Gaussian energy function
func TestGaussianPotential(t *testing.T) {
	tests := []struct {
		name     string
		phi, psi float64
		phi0, psi0 float64
		sigPhi, sigPsi float64
		maxEnergy float64 // Maximum expected energy
	}{
		{
			name: "At center (should be ~0)",
			phi: -60, psi: -45,
			phi0: -60, psi0: -45,
			sigPhi: 30, sigPsi: 30,
			maxEnergy: 0.01, // Should be very close to 0
		},
		{
			name: "One sigma away",
			phi: -30, psi: -45, // 30° away in phi (1σ)
			phi0: -60, psi0: -45,
			sigPhi: 30, sigPsi: 30,
			maxEnergy: 0.5, // Should be less than 0.5
		},
		{
			name: "Two sigma away (should be higher energy)",
			phi: 0, psi: -45, // 60° away in phi (2σ)
			phi0: -60, psi0: -45,
			sigPhi: 30, sigPsi: 30,
			maxEnergy: 1.0, // Should approach 1.0
		},
		{
			name: "Across periodic boundary",
			phi: 170, psi: 0,
			phi0: -170, psi0: 0, // 20° away with wraparound
			sigPhi: 30, sigPsi: 30,
			maxEnergy: 0.3, // Should be low energy (close angles)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			energy := gaussianPotential(tt.phi, tt.psi, tt.phi0, tt.psi0, tt.sigPhi, tt.sigPsi)
			if energy < 0 || energy > tt.maxEnergy {
				t.Errorf("gaussianPotential() = %v, expected < %v", energy, tt.maxEnergy)
			}
		})
	}
}

// TestGeneralRamachandran tests the general amino acid Ramachandran energy
func TestGeneralRamachandran(t *testing.T) {
	tests := []struct {
		name     string
		phi, psi float64
		maxEnergy float64 // Energy should be below this
		minEnergy float64 // Energy should be above this
		description string
	}{
		{
			name: "α-helix ideal angles",
			phi: -60, psi: -45,
			maxEnergy: 2.0, minEnergy: 0.0,
			description: "Should be very favorable (0-2 kcal/mol)",
		},
		{
			name: "β-sheet ideal angles",
			phi: -120, psi: 120,
			maxEnergy: 2.0, minEnergy: 0.0,
			description: "Should be very favorable (0-2 kcal/mol)",
		},
		{
			name: "PPII ideal angles",
			phi: -75, psi: 145,
			maxEnergy: 2.0, minEnergy: 0.0,
			description: "Should be very favorable (0-2 kcal/mol)",
		},
		{
			name: "Forbidden region (0, 0)",
			phi: 0, psi: 0,
			maxEnergy: 15.0, minEnergy: 5.0,
			description: "Should be unfavorable (5-15 kcal/mol)",
		},
		{
			name: "Borderline region",
			phi: -30, psi: 0,
			maxEnergy: 13.0, minEnergy: 2.0,
			description: "Should be moderately unfavorable (2-13 kcal/mol)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			energy := generalRamachandran(tt.phi, tt.psi)
			if energy < tt.minEnergy || energy > tt.maxEnergy {
				t.Errorf("%s: energy = %.2f kcal/mol, expected [%.1f, %.1f] kcal/mol",
					tt.description, energy, tt.minEnergy, tt.maxEnergy)
			} else {
				t.Logf("%s: energy = %.2f kcal/mol ✓", tt.description, energy)
			}
		})
	}
}

// TestGlycineRamachandran tests that glycine has lower energy penalties
func TestGlycineRamachandran(t *testing.T) {
	// Test same angles for glycine vs general amino acid in a region that's
	// forbidden for general AAs but allowed for glycine
	phi, psi := 90.0, 0.0 // Forbidden for general AAs, but glycine can access

	glyEnergy := glycineRamachandran(phi, psi)
	generalEnergy := generalRamachandran(phi, psi)

	if glyEnergy >= generalEnergy {
		t.Errorf("Glycine energy (%.2f) should be lower than general amino acid energy (%.2f) for extended region",
			glyEnergy, generalEnergy)
	} else {
		t.Logf("Glycine is more permissive: Gly = %.2f kcal/mol, General = %.2f kcal/mol ✓",
			glyEnergy, generalEnergy)
	}

	// Test that glycine still has low energy in α-helix region
	alphaEnergy := glycineRamachandran(-60, -45)
	if alphaEnergy > 3.0 {
		t.Errorf("Glycine in α-helix region should be favorable, got %.2f kcal/mol", alphaEnergy)
	} else {
		t.Logf("Glycine in α-helix: %.2f kcal/mol ✓", alphaEnergy)
	}

	// Test left-handed helix (should be favorable for glycine)
	leftHelixEnergy := glycineRamachandran(60, 45)
	if leftHelixEnergy > 2.0 {
		t.Errorf("Glycine in left-handed helix should be favorable, got %.2f kcal/mol", leftHelixEnergy)
	} else {
		t.Logf("Glycine in left-handed helix: %.2f kcal/mol ✓", leftHelixEnergy)
	}
}

// TestProlineRamachandran tests that proline has higher energy penalties
func TestProlineRamachandran(t *testing.T) {
	// Test proline at constrained phi
	phi, psi := -60.0, -30.0 // Helix-like proline conformation

	proEnergy := prolineRamachandran(phi, psi)

	if proEnergy > 2.0 {
		t.Errorf("Proline at ideal angles should be favorable, got %.2f kcal/mol", proEnergy)
	} else {
		t.Logf("Proline at ideal angles: %.2f kcal/mol ✓", proEnergy)
	}

	// Test proline at wrong phi (should be high energy)
	wrongPhiEnergy := prolineRamachandran(0, -30)
	if wrongPhiEnergy < 10.0 {
		t.Errorf("Proline at wrong phi should be unfavorable, got %.2f kcal/mol", wrongPhiEnergy)
	} else {
		t.Logf("Proline at wrong phi: %.2f kcal/mol (high penalty) ✓", wrongPhiEnergy)
	}
}

// TestGetRamachandranRegion tests secondary structure classification
func TestGetRamachandranRegion(t *testing.T) {
	tests := []struct {
		name     string
		phi, psi float64 // In radians
		expected string
	}{
		{
			name: "α-helix region",
			phi: -60 * math.Pi / 180, psi: -45 * math.Pi / 180,
			expected: "alpha-helix",
		},
		{
			name: "β-sheet region",
			phi: -120 * math.Pi / 180, psi: 120 * math.Pi / 180,
			expected: "beta-sheet",
		},
		{
			name: "PPII region",
			phi: -75 * math.Pi / 180, psi: 145 * math.Pi / 180,
			expected: "PPII",
		},
		{
			name: "Left-handed helix",
			phi: 60 * math.Pi / 180, psi: 45 * math.Pi / 180,
			expected: "left-helix",
		},
		{
			name: "Undefined region",
			phi: 0, psi: 0,
			expected: "other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region := GetRamachandranRegion(tt.phi, tt.psi)
			if region != tt.expected {
				t.Errorf("GetRamachandranRegion() = %v, expected %v", region, tt.expected)
			} else {
				t.Logf("Classified as %s ✓", region)
			}
		})
	}
}

// TestRamachandranEnergy tests the main energy function for different residue types
func TestRamachandranEnergy(t *testing.T) {
	tests := []struct {
		name        string
		phi, psi    float64 // In radians
		residueName string
		maxEnergy   float64
		description string
	}{
		{
			name: "Alanine in α-helix",
			phi: -60 * math.Pi / 180, psi: -45 * math.Pi / 180,
			residueName: "ALA",
			maxEnergy: 2.0,
			description: "Should be favorable",
		},
		{
			name: "Glycine in left-handed helix",
			phi: 60 * math.Pi / 180, psi: 45 * math.Pi / 180,
			residueName: "GLY",
			maxEnergy: 2.0,
			description: "Glycine allows left-handed helix",
		},
		{
			name: "Proline at constrained phi",
			phi: -60 * math.Pi / 180, psi: -30 * math.Pi / 180,
			residueName: "PRO",
			maxEnergy: 2.0,
			description: "Proline at ideal conformation",
		},
		{
			name: "Valine in forbidden region",
			phi: 0, psi: 0,
			residueName: "VAL",
			maxEnergy: 15.0,
			description: "Should be highly unfavorable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			energy := ramachandranEnergy(tt.phi, tt.psi, tt.residueName)
			if energy > tt.maxEnergy {
				t.Errorf("%s: energy = %.2f kcal/mol, expected < %.1f kcal/mol",
					tt.description, energy, tt.maxEnergy)
			} else {
				t.Logf("%s: energy = %.2f kcal/mol ✓", tt.description, energy)
			}
		})
	}
}

// BenchmarkRamachandranEnergy benchmarks the energy calculation performance
func BenchmarkRamachandranEnergy(b *testing.B) {
	phi := -60.0 * math.Pi / 180
	psi := -45.0 * math.Pi / 180

	b.Run("GeneralAA", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ramachandranEnergy(phi, psi, "ALA")
		}
	})

	b.Run("Glycine", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ramachandranEnergy(phi, psi, "GLY")
		}
	})

	b.Run("Proline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ramachandranEnergy(phi, psi, "PRO")
		}
	})
}
