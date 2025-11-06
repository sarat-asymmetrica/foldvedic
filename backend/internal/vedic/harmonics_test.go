package vedic

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

func TestCalculateVedicScore(t *testing.T) {
	// Load test protein
	protein, err := parser.ParsePDB("../../../testdata/test_peptide.pdb")
	if err != nil {
		t.Fatalf("Failed to parse test PDB: %v", err)
	}

	// Calculate Ramachandran angles
	angles := geometry.CalculateRamachandran(protein)

	// Calculate Vedic score
	score := CalculateVedicScore(protein, angles)

	t.Logf("Vedic Score Results:")
	t.Logf("  Golden Ratio Score: %.3f", score.GoldenRatioScore)
	t.Logf("  Digital Root Score: %.3f", score.DigitalRootScore)
	t.Logf("  Breathing Score:    %.3f", score.BreathingScore)
	t.Logf("  Total Score:        %.3f", score.TotalScore)
	t.Logf("  Helix Residues:     %d", score.NumHelixResidues)
	t.Logf("  Sheet Residues:     %d", score.NumSheetResidues)

	// Scores should be in [0, 1]
	if score.GoldenRatioScore < 0 || score.GoldenRatioScore > 1 {
		t.Errorf("Golden ratio score out of range: %.3f", score.GoldenRatioScore)
	}

	if score.DigitalRootScore < 0 || score.DigitalRootScore > 1 {
		t.Errorf("Digital root score out of range: %.3f", score.DigitalRootScore)
	}

	if score.BreathingScore < 0 || score.BreathingScore > 1 {
		t.Errorf("Breathing score out of range: %.3f", score.BreathingScore)
	}

	if score.TotalScore < 0 || score.TotalScore > 1 {
		t.Errorf("Total score out of range: %.3f", score.TotalScore)
	}
}

func TestGoldenRatioAlignment(t *testing.T) {
	// Test with pure helix angles
	helixAngles := []geometry.RamachandranAngles{
		{Phi: -60 * math.Pi / 180, Psi: -45 * math.Pi / 180}, // Helix
		{Phi: -60 * math.Pi / 180, Psi: -45 * math.Pi / 180}, // Helix
		{Phi: -60 * math.Pi / 180, Psi: -45 * math.Pi / 180}, // Helix
	}

	protein := &parser.Protein{Residues: make([]*parser.Residue, 3)}

	score := calculateGoldenRatioAlignment(protein, helixAngles)

	t.Logf("Pure helix golden ratio score: %.3f", score)

	// Pure helix should score well (high secondary structure content)
	if score < 0.5 {
		t.Errorf("Pure helix should score >= 0.5, got %.3f", score)
	}

	// Test with pure sheet angles
	sheetAngles := []geometry.RamachandranAngles{
		{Phi: -120 * math.Pi / 180, Psi: +120 * math.Pi / 180}, // Sheet
		{Phi: -120 * math.Pi / 180, Psi: +120 * math.Pi / 180}, // Sheet
		{Phi: -120 * math.Pi / 180, Psi: +120 * math.Pi / 180}, // Sheet
	}

	scoreSheet := calculateGoldenRatioAlignment(protein, sheetAngles)
	t.Logf("Pure sheet golden ratio score: %.3f", scoreSheet)

	// Pure sheet should also score well
	if scoreSheet < 0.5 {
		t.Errorf("Pure sheet should score >= 0.5, got %.3f", scoreSheet)
	}

	// Test with random coil (low structure)
	coilAngles := []geometry.RamachandranAngles{
		{Phi: 0, Psi: 0},                                      // Forbidden region
		{Phi: +90 * math.Pi / 180, Psi: -90 * math.Pi / 180}, // Unusual
		{Phi: +45 * math.Pi / 180, Psi: +45 * math.Pi / 180}, // Extended
	}

	scoreCoil := calculateGoldenRatioAlignment(protein, coilAngles)
	t.Logf("Random coil golden ratio score: %.3f", scoreCoil)

	// Random coil should score lower than structured
	if scoreCoil > score {
		t.Error("Random coil should score lower than helix")
	}
}

func TestDigitalRootConsistency(t *testing.T) {
	// Test with typical helix angles
	helixAngles := []geometry.RamachandranAngles{
		{Phi: -60 * math.Pi / 180, Psi: -45 * math.Pi / 180},
		{Phi: -60 * math.Pi / 180, Psi: -45 * math.Pi / 180},
	}

	score := calculateDigitalRootConsistency(helixAngles)

	t.Logf("Helix digital root consistency: %.3f", score)

	// Should have some consistency (exact value depends on DR patterns)
	if score < 0 || score > 1 {
		t.Errorf("Digital root score out of range: %.3f", score)
	}

	// Test with empty angles
	emptyScore := calculateDigitalRootConsistency([]geometry.RamachandranAngles{})
	if emptyScore != 0 {
		t.Errorf("Empty angles should give score 0, got %.3f", emptyScore)
	}
}

func TestDigitalRoot(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{9, 9},
		{10, 1},  // 1 + 0 = 1
		{38, 2},  // 3 + 8 = 11 → 1 + 1 = 2
		{123, 6}, // 1 + 2 + 3 = 6
		{999, 9}, // 9 + 9 + 9 = 27 → 2 + 7 = 9
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := digitalRoot(tt.input)
			if result != tt.expected {
				t.Errorf("digitalRoot(%d) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStructuralBreathingScore(t *testing.T) {
	// Load test protein
	protein, err := parser.ParsePDB("../../../testdata/test_peptide.pdb")
	if err != nil {
		t.Fatalf("Failed to parse test PDB: %v", err)
	}

	score := calculateStructuralBreathingScore(protein)

	t.Logf("Structural breathing score: %.3f", score)

	// Score should be in [0, 1]
	if score < 0 || score > 1 {
		t.Errorf("Breathing score out of range: %.3f", score)
	}

	// Folded protein should have reasonable compactness
	if score == 0 {
		t.Error("Folded protein should have non-zero breathing score")
	}
}

func TestGoldenRatioConstants(t *testing.T) {
	// Verify golden ratio constants are correct
	tolerance := 0.000001

	// φ = (1 + √5) / 2
	expectedPhi := (1.0 + math.Sqrt(5.0)) / 2.0
	if math.Abs(Phi-expectedPhi) > tolerance {
		t.Errorf("Phi constant incorrect: %.15f, expected %.15f", Phi, expectedPhi)
	}

	// 1/φ
	expectedPhiInverse := 1.0 / Phi
	if math.Abs(PhiInverse-expectedPhiInverse) > tolerance {
		t.Errorf("PhiInverse incorrect: %.15f, expected %.15f", PhiInverse, expectedPhiInverse)
	}

	// φ²
	expectedPhiSquared := Phi * Phi
	if math.Abs(PhiSquared-expectedPhiSquared) > tolerance {
		t.Errorf("PhiSquared incorrect: %.15f, expected %.15f", PhiSquared, expectedPhiSquared)
	}

	// Verify Fibonacci property: φ² = φ + 1
	if math.Abs(PhiSquared-(Phi+1.0)) > tolerance {
		t.Error("Golden ratio should satisfy φ² = φ + 1")
	}
}

func TestHelixPitchGoldenRelation(t *testing.T) {
	// Verify helix pitch relates to golden ratio
	// Standard helix: 3.6 residues per turn
	// Golden helix: 10/φ² = 3.819 residues per turn

	standardPitch := 3.6
	goldenPitch := 10.0 / PhiSquared

	t.Logf("Standard helix pitch: %.3f residues/turn", standardPitch)
	t.Logf("Golden helix pitch: %.3f residues/turn", goldenPitch)

	// Difference should be ~6%
	difference := math.Abs(goldenPitch - standardPitch)
	percentDiff := difference / standardPitch * 100

	t.Logf("Difference: %.3f (%.1f%%)", difference, percentDiff)

	if percentDiff > 10 {
		t.Errorf("Helix pitch difference too large: %.1f%% (expect ~6%%)", percentDiff)
	}
}

func BenchmarkCalculateVedicScore(b *testing.B) {
	protein, _ := parser.ParsePDB("../../../testdata/test_peptide.pdb")
	angles := geometry.CalculateRamachandran(protein)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalculateVedicScore(protein, angles)
	}
}
