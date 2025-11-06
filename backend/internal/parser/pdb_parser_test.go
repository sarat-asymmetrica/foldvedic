package parser

import (
	"testing"
)

func TestParsePDB(t *testing.T) {
	// Test parsing the test peptide
	protein, err := ParsePDB("../../../testdata/test_peptide.pdb")
	if err != nil {
		t.Fatalf("Failed to parse PDB: %v", err)
	}

	// Check protein was parsed
	if protein == nil {
		t.Fatal("Protein is nil")
	}

	// Check residue count (3 residues: ALA, GLY, VAL)
	if len(protein.Residues) != 3 {
		t.Errorf("Expected 3 residues, got %d", len(protein.Residues))
	}

	// Check atom count (12 backbone atoms)
	if len(protein.Atoms) != 12 {
		t.Errorf("Expected 12 atoms, got %d", len(protein.Atoms))
	}

	// Check first residue
	if len(protein.Residues) > 0 {
		res := protein.Residues[0]
		if res.Name != "ALA" {
			t.Errorf("Expected first residue to be ALA, got %s", res.Name)
		}
		if !res.HasCompleteBackbone() {
			t.Error("First residue should have complete backbone")
		}
	}

	// Check complete backbone count
	complete := protein.NumCompleteResidues()
	if complete != 3 {
		t.Errorf("Expected 3 complete residues, got %d", complete)
	}
}

func TestParseAtomLine(t *testing.T) {
	// Test parsing a typical ATOM line
	line := "ATOM      1  N   ALA A   1      11.104   6.134  -6.504  1.00  0.00           N"
	atom, err := parseAtomLine(line)

	if err != nil {
		t.Fatalf("Failed to parse atom line: %v", err)
	}

	// Check serial number
	if atom.Serial != 1 {
		t.Errorf("Expected serial 1, got %d", atom.Serial)
	}

	// Check atom name
	if atom.Name != "N" {
		t.Errorf("Expected atom name 'N', got '%s'", atom.Name)
	}

	// Check residue name
	if atom.ResName != "ALA" {
		t.Errorf("Expected residue 'ALA', got '%s'", atom.ResName)
	}

	// Check coordinates (with small tolerance for floating point)
	tolerance := 0.001
	if abs(atom.X-11.104) > tolerance {
		t.Errorf("Expected X=11.104, got %f", atom.X)
	}
	if abs(atom.Y-6.134) > tolerance {
		t.Errorf("Expected Y=6.134, got %f", atom.Y)
	}
	if abs(atom.Z-(-6.504)) > tolerance {
		t.Errorf("Expected Z=-6.504, got %f", atom.Z)
	}

	// Check chain ID
	if atom.ChainID != "A" {
		t.Errorf("Expected chain 'A', got '%s'", atom.ChainID)
	}

	// Check residue sequence number
	if atom.ResSeq != 1 {
		t.Errorf("Expected resSeq 1, got %d", atom.ResSeq)
	}
}

func TestIsBackboneAtom(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"N", true},
		{"CA", true},
		{"C", true},
		{"O", true},
		{"CB", false},
		{"CD", false},
		{"OXT", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBackboneAtom(tt.name)
			if result != tt.expected {
				t.Errorf("isBackboneAtom(%s) = %v, expected %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestResidueHasCompleteBackbone(t *testing.T) {
	// Complete backbone
	res1 := &Residue{
		Name: "ALA",
		N:    &Atom{Name: "N"},
		CA:   &Atom{Name: "CA"},
		C:    &Atom{Name: "C"},
	}
	if !res1.HasCompleteBackbone() {
		t.Error("Residue with N, CA, C should have complete backbone")
	}

	// Missing CA
	res2 := &Residue{
		Name: "ALA",
		N:    &Atom{Name: "N"},
		C:    &Atom{Name: "C"},
	}
	if res2.HasCompleteBackbone() {
		t.Error("Residue missing CA should not have complete backbone")
	}

	// Missing all
	res3 := &Residue{
		Name: "ALA",
	}
	if res3.HasCompleteBackbone() {
		t.Error("Residue with no atoms should not have complete backbone")
	}
}

// Helper function for floating point comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
