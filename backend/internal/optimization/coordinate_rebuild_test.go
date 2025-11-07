// Package optimization - Test for SetDihedrals() coordinate rebuilding
//
// WRIGHT BROTHERS TEST: Simple 3-residue peptide (Ala-Ala-Ala)
//
// Test that SetDihedrals():
// 1. Rebuilds coordinates without NaN
// 2. Preserves bond lengths (~1.458 Å N-CA, ~1.523 Å CA-C)
// 3. Doesn't cause protein to explode
//
package optimization

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// TestSetDihedralsCoordinateRebuild tests the fixed SetDihedrals() function
func TestSetDihedralsCoordinateRebuild(t *testing.T) {
	// Build simple 3-residue peptide: Ala-Ala-Ala
	sequence := "AAA"

	// Start with extended conformation
	angles := []geometry.RamachandranAngles{
		{Phi: -120.0 * math.Pi / 180.0, Psi: 120.0 * math.Pi / 180.0},
		{Phi: -120.0 * math.Pi / 180.0, Psi: 120.0 * math.Pi / 180.0},
		{Phi: -120.0 * math.Pi / 180.0, Psi: 120.0 * math.Pi / 180.0},
	}

	// Build initial protein
	protein, err := geometry.BuildProteinFromAngles(sequence, angles)
	if err != nil {
		t.Fatalf("Failed to build initial protein: %v", err)
	}

	// Check initial coordinates are valid
	if !validateCoordinates(protein, t) {
		t.Fatal("Initial coordinates invalid")
	}

	t.Logf("Initial coordinates valid")
	printProteinStats(protein, t)

	// Now modify angles and rebuild (this tests SetDihedrals)
	newAngles := []geometry.RamachandranAngles{
		{Phi: -60.0 * math.Pi / 180.0, Psi: -40.0 * math.Pi / 180.0}, // Alpha helix region
		{Phi: -60.0 * math.Pi / 180.0, Psi: -40.0 * math.Pi / 180.0},
		{Phi: -60.0 * math.Pi / 180.0, Psi: -40.0 * math.Pi / 180.0},
	}

	err = SetDihedrals(protein, newAngles)
	if err != nil {
		t.Fatalf("SetDihedrals failed: %v", err)
	}

	// Validate rebuilt coordinates
	if !validateCoordinates(protein, t) {
		t.Fatal("Rebuilt coordinates invalid (NaN or explosion)")
	}

	t.Logf("Rebuilt coordinates valid")
	printProteinStats(protein, t)

	// Validate bond lengths
	if !validateBondLengths(protein, t) {
		t.Fatal("Bond lengths out of range")
	}

	t.Logf("✅ SetDihedrals coordinate rebuild TEST PASSED")
}

// validateCoordinates checks for NaN and explosions
func validateCoordinates(protein *parser.Protein, t *testing.T) bool {
	for i, res := range protein.Residues {
		atoms := []*parser.Atom{res.N, res.CA, res.C, res.O}
		names := []string{"N", "CA", "C", "O"}

		for j, atom := range atoms {
			if atom == nil {
				continue
			}

			// Check for NaN
			if math.IsNaN(atom.X) || math.IsNaN(atom.Y) || math.IsNaN(atom.Z) {
				t.Errorf("Residue %d %s: NaN coordinate (%.2f, %.2f, %.2f)",
					i, names[j], atom.X, atom.Y, atom.Z)
				return false
			}

			// Check for explosion (atoms should be within 100 Å of origin)
			r := math.Sqrt(atom.X*atom.X + atom.Y*atom.Y + atom.Z*atom.Z)
			if r > 100.0 {
				t.Errorf("Residue %d %s: Exploded! Distance from origin: %.2f Å",
					i, names[j], r)
				return false
			}
		}
	}

	return true
}

// validateBondLengths checks if backbone bond lengths are reasonable
func validateBondLengths(protein *parser.Protein, t *testing.T) bool {
	allValid := true

	for i, res := range protein.Residues {
		// Check N-CA bond
		if res.N != nil && res.CA != nil {
			nca := distance(res.N, res.CA)
			expected := 1.458 // Å
			if math.Abs(nca-expected) > 0.2 {
				t.Errorf("Residue %d N-CA bond: %.3f Å (expected %.3f ± 0.2)",
					i, nca, expected)
				allValid = false
			}
		}

		// Check CA-C bond
		if res.CA != nil && res.C != nil {
			cac := distance(res.CA, res.C)
			expected := 1.523 // Å
			if math.Abs(cac-expected) > 0.2 {
				t.Errorf("Residue %d CA-C bond: %.3f Å (expected %.3f ± 0.2)",
					i, cac, expected)
				allValid = false
			}
		}

		// Check C-N peptide bond (if not last residue)
		if i < len(protein.Residues)-1 {
			nextN := protein.Residues[i+1].N
			if res.C != nil && nextN != nil {
				cn := distance(res.C, nextN)
				expected := 1.329 // Å
				if math.Abs(cn-expected) > 0.5 { // More tolerance for peptide bond
					t.Errorf("Residue %d C-N bond: %.3f Å (expected %.3f ± 0.5)",
						i, cn, expected)
					allValid = false
				}
			}
		}

		// Check C=O bond
		if res.C != nil && res.O != nil {
			co := distance(res.C, res.O)
			expected := 1.231 // Å
			if math.Abs(co-expected) > 0.2 {
				t.Errorf("Residue %d C=O bond: %.3f Å (expected %.3f ± 0.2)",
					i, co, expected)
				allValid = false
			}
		}
	}

	return allValid
}

// distance calculates Euclidean distance between two atoms
func distance(a1, a2 *parser.Atom) float64 {
	dx := a2.X - a1.X
	dy := a2.Y - a1.Y
	dz := a2.Z - a1.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// printProteinStats prints diagnostic info
func printProteinStats(protein *parser.Protein, t *testing.T) {
	for i, res := range protein.Residues {
		if res.N != nil && res.CA != nil && res.C != nil {
			t.Logf("  Residue %d: N(%.2f, %.2f, %.2f) CA(%.2f, %.2f, %.2f) C(%.2f, %.2f, %.2f)",
				i,
				res.N.X, res.N.Y, res.N.Z,
				res.CA.X, res.CA.Y, res.CA.Z,
				res.C.X, res.C.Y, res.C.Z)
		}
	}
}
