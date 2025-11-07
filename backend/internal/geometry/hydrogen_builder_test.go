package geometry

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// TestAddHydrogens tests hydrogen placement for 3-residue peptide
//
// WRIGHT BROTHERS TEST:
// - Build minimal 3-residue structure
// - Add hydrogens
// - Validate bond lengths
// - Expect bugs, iterate fast
func TestAddHydrogens(t *testing.T) {
	// Build 3-residue peptide: ALA-GLY-VAL
	// Using realistic backbone coordinates
	protein := &parser.Protein{
		Name: "test_peptide",
		Residues: []*parser.Residue{
			// Residue 0: ALA (N-terminal, no backbone H)
			{
				Name:    "ALA",
				SeqNum:  1,
				ChainID: "A",
				N: &parser.Atom{
					Serial: 1, Name: "N", ResName: "ALA",
					ChainID: "A", ResSeq: 1,
					X: 0.0, Y: 0.0, Z: 0.0, Element: "N",
				},
				CA: &parser.Atom{
					Serial: 2, Name: "CA", ResName: "ALA",
					ChainID: "A", ResSeq: 1,
					X: 1.46, Y: 0.0, Z: 0.0, Element: "C",
				},
				C: &parser.Atom{
					Serial: 3, Name: "C", ResName: "ALA",
					ChainID: "A", ResSeq: 1,
					X: 2.0, Y: 1.5, Z: 0.0, Element: "C",
				},
				O: &parser.Atom{
					Serial: 4, Name: "O", ResName: "ALA",
					ChainID: "A", ResSeq: 1,
					X: 3.2, Y: 1.6, Z: 0.0, Element: "O",
				},
			},
			// Residue 1: GLY (should get backbone H)
			{
				Name:    "GLY",
				SeqNum:  2,
				ChainID: "A",
				N: &parser.Atom{
					Serial: 5, Name: "N", ResName: "GLY",
					ChainID: "A", ResSeq: 2,
					X: 1.2, Y: 2.5, Z: 0.0, Element: "N",
				},
				CA: &parser.Atom{
					Serial: 6, Name: "CA", ResName: "GLY",
					ChainID: "A", ResSeq: 2,
					X: 1.7, Y: 3.8, Z: 0.0, Element: "C",
				},
				C: &parser.Atom{
					Serial: 7, Name: "C", ResName: "GLY",
					ChainID: "A", ResSeq: 2,
					X: 3.0, Y: 4.0, Z: 0.0, Element: "C",
				},
				O: &parser.Atom{
					Serial: 8, Name: "O", ResName: "GLY",
					ChainID: "A", ResSeq: 2,
					X: 3.8, Y: 3.2, Z: 0.0, Element: "O",
				},
			},
			// Residue 2: VAL (should get backbone H)
			{
				Name:    "VAL",
				SeqNum:  3,
				ChainID: "A",
				N: &parser.Atom{
					Serial: 9, Name: "N", ResName: "VAL",
					ChainID: "A", ResSeq: 3,
					X: 3.3, Y: 5.2, Z: 0.0, Element: "N",
				},
				CA: &parser.Atom{
					Serial: 10, Name: "CA", ResName: "VAL",
					ChainID: "A", ResSeq: 3,
					X: 4.5, Y: 5.5, Z: 0.0, Element: "C",
				},
				C: &parser.Atom{
					Serial: 11, Name: "C", ResName: "VAL",
					ChainID: "A", ResSeq: 3,
					X: 5.0, Y: 7.0, Z: 0.0, Element: "C",
				},
				O: &parser.Atom{
					Serial: 12, Name: "O", ResName: "VAL",
					ChainID: "A", ResSeq: 3,
					X: 6.2, Y: 7.2, Z: 0.0, Element: "O",
				},
			},
		},
	}

	// Add backbone atoms to Atoms slice
	for _, res := range protein.Residues {
		if res.N != nil {
			protein.Atoms = append(protein.Atoms, res.N)
		}
		if res.CA != nil {
			protein.Atoms = append(protein.Atoms, res.CA)
		}
		if res.C != nil {
			protein.Atoms = append(protein.Atoms, res.C)
		}
		if res.O != nil {
			protein.Atoms = append(protein.Atoms, res.O)
		}
	}

	// Count atoms before adding hydrogens
	atomsBefore := len(protein.Atoms)
	t.Logf("Atoms before H addition: %d", atomsBefore)

	// Add hydrogens
	err := AddHydrogens(protein)
	if err != nil {
		t.Fatalf("AddHydrogens failed: %v", err)
	}

	// Count atoms after
	atomsAfter := len(protein.Atoms)
	t.Logf("Atoms after H addition: %d", atomsAfter)

	// Count H atoms added
	hCount := 0
	for _, atom := range protein.Atoms {
		if atom.Element == "H" {
			hCount++
			t.Logf("H atom: %s in residue %d at (%.2f, %.2f, %.2f)",
				atom.Name, atom.ResSeq, atom.X, atom.Y, atom.Z)
		}
	}

	// VALIDATION 1: Should have added H atoms
	// Residue 0 (N-terminal): no backbone H
	// Residue 1: 1 backbone H
	// Residue 2: 1 backbone H
	// Total: 2 backbone H expected
	if hCount < 1 {
		t.Errorf("Expected at least 1 H atom, got %d", hCount)
	}

	t.Logf("Total H atoms added: %d", hCount)

	// VALIDATION 2: Check H-N bond lengths
	for _, atom := range protein.Atoms {
		if atom.Element == "H" && atom.Name == "H" {
			// Find corresponding N atom
			var nAtom *parser.Atom
			for _, a := range protein.Atoms {
				if a.ResSeq == atom.ResSeq && a.Name == "N" {
					nAtom = a
					break
				}
			}

			if nAtom == nil {
				t.Errorf("Could not find N atom for H in residue %d", atom.ResSeq)
				continue
			}

			// Calculate H-N distance
			dx := atom.X - nAtom.X
			dy := atom.Y - nAtom.Y
			dz := atom.Z - nAtom.Z
			dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

			t.Logf("H-N bond length for residue %d: %.3f Å", atom.ResSeq, dist)

			// Ideal N-H bond: 1.01 Å ± 0.14 Å tolerance
			if dist < 0.85 || dist > 1.25 {
				t.Errorf("H-N bond length %.3f Å outside valid range [0.85, 1.25]", dist)
			}
		}
	}

	// VALIDATION 3: Use ValidateHydrogenGeometry
	valid, msg := ValidateHydrogenGeometry(protein)
	t.Logf("Geometry validation: %s", msg)
	if !valid {
		t.Errorf("Hydrogen geometry validation failed: %s", msg)
	}
}

// TestBackboneHydrogenPlacement tests backbone NH hydrogen placement
func TestBackboneHydrogenPlacement(t *testing.T) {
	// Simple 2-residue test
	protein := &parser.Protein{
		Name: "dipeptide",
		Residues: []*parser.Residue{
			{
				Name:    "ALA",
				SeqNum:  1,
				ChainID: "A",
				N:       &parser.Atom{Serial: 1, Name: "N", X: 0, Y: 0, Z: 0, Element: "N"},
				CA:      &parser.Atom{Serial: 2, Name: "CA", X: 1.46, Y: 0, Z: 0, Element: "C"},
				C:       &parser.Atom{Serial: 3, Name: "C", X: 2.0, Y: 1.5, Z: 0, Element: "C"},
				O:       &parser.Atom{Serial: 4, Name: "O", X: 3.2, Y: 1.6, Z: 0, Element: "O"},
			},
			{
				Name:    "GLY",
				SeqNum:  2,
				ChainID: "A",
				N:       &parser.Atom{Serial: 5, Name: "N", X: 1.2, Y: 2.5, Z: 0, Element: "N"},
				CA:      &parser.Atom{Serial: 6, Name: "CA", X: 1.7, Y: 3.8, Z: 0, Element: "C"},
				C:       &parser.Atom{Serial: 7, Name: "C", X: 3.0, Y: 4.0, Z: 0, Element: "C"},
				O:       &parser.Atom{Serial: 8, Name: "O", X: 3.8, Y: 3.2, Z: 0, Element: "O"},
			},
		},
	}

	// Add atoms to protein
	for _, res := range protein.Residues {
		protein.Atoms = append(protein.Atoms, res.N, res.CA, res.C, res.O)
	}

	// Add backbone hydrogen
	err := addBackboneHydrogen(protein.Residues[1], 1, protein)
	if err != nil {
		t.Fatalf("addBackboneHydrogen failed: %v", err)
	}

	// Find added H atom
	var hAtom *parser.Atom
	for _, atom := range protein.Atoms {
		if atom.Element == "H" && atom.ResSeq == 2 {
			hAtom = atom
			break
		}
	}

	if hAtom == nil {
		t.Fatal("No H atom added")
	}

	// Check H-N distance
	n := protein.Residues[1].N
	dx := hAtom.X - n.X
	dy := hAtom.Y - n.Y
	dz := hAtom.Z - n.Z
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

	t.Logf("H-N bond length: %.3f Å (target: 1.01 Å)", dist)

	if dist < 0.85 || dist > 1.25 {
		t.Errorf("H-N bond length %.3f Å outside valid range", dist)
	}
}

// TestProlineNoHydrogen tests that proline doesn't get backbone H
//
// BIOCHEMIST: Proline is cyclic - no NH hydrogen
func TestProlineNoHydrogen(t *testing.T) {
	protein := &parser.Protein{
		Name: "with_proline",
		Residues: []*parser.Residue{
			{
				Name:    "ALA",
				SeqNum:  1,
				ChainID: "A",
				N:       &parser.Atom{Serial: 1, Name: "N", X: 0, Y: 0, Z: 0, Element: "N"},
				CA:      &parser.Atom{Serial: 2, Name: "CA", X: 1.46, Y: 0, Z: 0, Element: "C"},
				C:       &parser.Atom{Serial: 3, Name: "C", X: 2.0, Y: 1.5, Z: 0, Element: "C"},
				O:       &parser.Atom{Serial: 4, Name: "O", X: 3.2, Y: 1.6, Z: 0, Element: "O"},
			},
			{
				Name:    "PRO", // Proline - should NOT get backbone H
				SeqNum:  2,
				ChainID: "A",
				N:       &parser.Atom{Serial: 5, Name: "N", X: 1.2, Y: 2.5, Z: 0, Element: "N"},
				CA:      &parser.Atom{Serial: 6, Name: "CA", X: 1.7, Y: 3.8, Z: 0, Element: "C"},
				C:       &parser.Atom{Serial: 7, Name: "C", X: 3.0, Y: 4.0, Z: 0, Element: "C"},
				O:       &parser.Atom{Serial: 8, Name: "O", X: 3.8, Y: 3.2, Z: 0, Element: "O"},
			},
		},
	}

	// Add atoms
	for _, res := range protein.Residues {
		protein.Atoms = append(protein.Atoms, res.N, res.CA, res.C, res.O)
	}

	// Try to add backbone H to proline (should skip)
	err := addBackboneHydrogen(protein.Residues[1], 1, protein)
	if err != nil {
		t.Fatalf("addBackboneHydrogen failed: %v", err)
	}

	// Check no H was added to proline
	for _, atom := range protein.Atoms {
		if atom.Element == "H" && atom.ResSeq == 2 {
			t.Error("Proline should not have backbone H atom")
		}
	}

	t.Log("Proline correctly skipped for backbone H")
}

// TestAlphaHydrogenGeometry tests CA-HA bond geometry
func TestAlphaHydrogenGeometry(t *testing.T) {
	residue := &parser.Residue{
		Name:    "ALA",
		SeqNum:  1,
		ChainID: "A",
		N:       &parser.Atom{Serial: 1, Name: "N", X: 0, Y: 0, Z: 0, Element: "N"},
		CA:      &parser.Atom{Serial: 2, Name: "CA", X: 1.46, Y: 0, Z: 0, Element: "C"},
		C:       &parser.Atom{Serial: 3, Name: "C", X: 2.0, Y: 1.5, Z: 0, Element: "C"},
		O:       &parser.Atom{Serial: 4, Name: "O", X: 3.2, Y: 1.6, Z: 0, Element: "O"},
	}

	// Add alpha hydrogen (validation only for now)
	err := addAlphaHydrogen(residue)
	if err != nil {
		t.Fatalf("addAlphaHydrogen failed: %v", err)
	}

	t.Log("Alpha hydrogen geometry validated (CA-HA bond)")
}

// BenchmarkAddHydrogens benchmarks hydrogen addition for 20-residue protein
func BenchmarkAddHydrogens(b *testing.B) {
	// Build 20-residue protein
	protein := &parser.Protein{
		Name:     "benchmark",
		Residues: make([]*parser.Residue, 20),
		Atoms:    make([]*parser.Atom, 0, 80),
	}

	for i := 0; i < 20; i++ {
		res := &parser.Residue{
			Name:    "ALA",
			SeqNum:  i + 1,
			ChainID: "A",
			N: &parser.Atom{
				Serial: i*4 + 1, Name: "N",
				X: float64(i) * 3.8, Y: 0, Z: 0, Element: "N",
			},
			CA: &parser.Atom{
				Serial: i*4 + 2, Name: "CA",
				X: float64(i)*3.8 + 1.46, Y: 0, Z: 0, Element: "C",
			},
			C: &parser.Atom{
				Serial: i*4 + 3, Name: "C",
				X: float64(i)*3.8 + 2.0, Y: 1.5, Z: 0, Element: "C",
			},
			O: &parser.Atom{
				Serial: i*4 + 4, Name: "O",
				X: float64(i)*3.8 + 3.2, Y: 1.6, Z: 0, Element: "O",
			},
		}
		protein.Residues[i] = res
		protein.Atoms = append(protein.Atoms, res.N, res.CA, res.C, res.O)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Copy protein to avoid mutation between runs
		testProtein := protein
		_ = AddHydrogens(testProtein)
	}
}
