package physics

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

func TestDetectClashes(t *testing.T) {
	// Create structure with known clash
	atom1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, ChainID: "A", X: 0, Y: 0, Z: 0}
	atom2 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 3, ChainID: "A", X: 0.5, Y: 0, Z: 0} // Too close!
	atom3 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 5, ChainID: "A", X: 10, Y: 0, Z: 0}  // Far away

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", CA: atom1},
			{Name: "GLY", SeqNum: 3, ChainID: "A", CA: atom2},
			{Name: "VAL", SeqNum: 5, ChainID: "A", CA: atom3},
		},
		Atoms: []*parser.Atom{atom1, atom2, atom3},
	}

	report := DetectClashes(protein)

	if !report.HasClashes {
		t.Error("Should detect clash at 0.5 Å distance")
	}

	if report.ClashCount != 1 {
		t.Errorf("Expected 1 clash, got %d", report.ClashCount)
	}

	if report.WorstClashDist > 1.0 {
		t.Errorf("Worst clash distance should be ~0.5 Å, got %.2f", report.WorstClashDist)
	}
}

func TestDetectClashes_NoClash(t *testing.T) {
	// Create structure with normal spacing
	atom1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, ChainID: "A", X: 0, Y: 0, Z: 0}
	atom2 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 3, ChainID: "A", X: 4.0, Y: 0, Z: 0} // Normal distance

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", CA: atom1},
			{Name: "GLY", SeqNum: 3, ChainID: "A", CA: atom2},
		},
		Atoms: []*parser.Atom{atom1, atom2},
	}

	report := DetectClashes(protein)

	if report.HasClashes {
		t.Error("Should not detect clash at 4.0 Å distance")
	}

	if report.ClashCount != 0 {
		t.Errorf("Expected 0 clashes, got %d", report.ClashCount)
	}
}

func TestValidateCoordinates_NaN(t *testing.T) {
	// Test NaN detection
	atom1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, ResName: "ALA", X: math.NaN(), Y: 0, Z: 0}

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, CA: atom1},
		},
		Atoms: []*parser.Atom{atom1},
	}

	report := ValidateCoordinates(protein)

	if report.IsValid {
		t.Error("Should detect NaN coordinate")
	}

	if report.ValidationError == "" {
		t.Error("Should provide error message")
	}

	if len(report.ValidationError) < 10 {
		t.Errorf("Error message too short: %s", report.ValidationError)
	}
}

func TestValidateCoordinates_Infinity(t *testing.T) {
	// Test infinity detection
	atom1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, ResName: "ALA", X: math.Inf(1), Y: 0, Z: 0}

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, CA: atom1},
		},
		Atoms: []*parser.Atom{atom1},
	}

	report := ValidateCoordinates(protein)

	if report.IsValid {
		t.Error("Should detect Inf coordinate")
	}
}

func TestValidateCoordinates_TooFar(t *testing.T) {
	// Test unrealistic distance
	atom1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, ResName: "ALA", X: 2000, Y: 0, Z: 0}

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, CA: atom1},
		},
		Atoms: []*parser.Atom{atom1},
	}

	report := ValidateCoordinates(protein)

	if report.IsValid {
		t.Error("Should detect atom too far from origin")
	}
}

func TestValidateCoordinates_BrokenBackbone(t *testing.T) {
	// Test broken peptide bond
	atomC := &parser.Atom{Name: "C", Element: "C", ResSeq: 1, X: 0, Y: 0, Z: 0}
	atomN := &parser.Atom{Name: "N", Element: "N", ResSeq: 2, X: 10, Y: 0, Z: 0} // Way too far!

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", C: atomC},
			{Name: "GLY", SeqNum: 2, ChainID: "A", N: atomN},
		},
		Atoms: []*parser.Atom{atomC, atomN},
	}

	report := ValidateCoordinates(protein)

	if report.IsValid {
		t.Error("Should detect broken peptide bond")
	}

	if report.ValidationError == "" || len(report.ValidationError) < 20 {
		t.Errorf("Should provide detailed error message, got: %s", report.ValidationError)
	}
}

func TestValidateCoordinates_ValidStructure(t *testing.T) {
	// Test valid structure
	atomN1 := &parser.Atom{Name: "N", Element: "N", ResSeq: 1, X: 0, Y: 0, Z: 0}
	atomCA1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, X: 1.5, Y: 0, Z: 0}
	atomC1 := &parser.Atom{Name: "C", Element: "C", ResSeq: 1, X: 3.0, Y: 0, Z: 0}
	atomN2 := &parser.Atom{Name: "N", Element: "N", ResSeq: 2, X: 4.3, Y: 0, Z: 0} // ~1.3 Å from previous C
	atomCA2 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 2, X: 5.8, Y: 0, Z: 0}
	atomC2 := &parser.Atom{Name: "C", Element: "C", ResSeq: 2, X: 7.3, Y: 0, Z: 0}

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", N: atomN1, CA: atomCA1, C: atomC1},
			{Name: "GLY", SeqNum: 2, ChainID: "A", N: atomN2, CA: atomCA2, C: atomC2},
		},
		Atoms: []*parser.Atom{atomN1, atomCA1, atomC1, atomN2, atomCA2, atomC2},
	}

	report := ValidateCoordinates(protein)

	if !report.IsValid {
		t.Errorf("Should validate correct structure, got error: %s", report.ValidationError)
	}
}

func TestScoreStructureQuality(t *testing.T) {
	// Test quality scoring
	atom1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, ChainID: "A", X: 0, Y: 0, Z: 0}
	atom2 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 3, ChainID: "A", X: 4.0, Y: 0, Z: 0}

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", CA: atom1},
			{Name: "GLY", SeqNum: 3, ChainID: "A", CA: atom2},
		},
		Atoms: []*parser.Atom{atom1, atom2},
	}

	quality, report := ScoreStructureQuality(protein)

	if !report.IsValid {
		t.Errorf("Should validate structure, got error: %s", report.ValidationError)
	}

	if quality < 0.9 {
		t.Errorf("Expected high quality score for clash-free structure, got %.2f", quality)
	}

	if report.HasClashes {
		t.Error("Should not report clashes for well-spaced atoms")
	}
}

func TestScoreStructureQuality_WithClashes(t *testing.T) {
	// Create structure with multiple clashes
	atom1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, ChainID: "A", X: 0, Y: 0, Z: 0}
	atom2 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 3, ChainID: "A", X: 0.5, Y: 0, Z: 0}
	atom3 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 5, ChainID: "A", X: 0.6, Y: 0, Z: 0}

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", CA: atom1},
			{Name: "GLY", SeqNum: 3, ChainID: "A", CA: atom2},
			{Name: "VAL", SeqNum: 5, ChainID: "A", CA: atom3},
		},
		Atoms: []*parser.Atom{atom1, atom2, atom3},
	}

	quality, report := ScoreStructureQuality(protein)

	if quality > 0.9 {
		t.Errorf("Expected low quality score for structure with clashes, got %.2f", quality)
	}

	if !report.HasClashes {
		t.Error("Should detect clashes")
	}

	if report.ClashCount < 2 {
		t.Errorf("Expected at least 2 clashes, got %d", report.ClashCount)
	}
}

func TestScoreStructureQuality_Invalid(t *testing.T) {
	// Test invalid structure (NaN)
	atom1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, X: math.NaN(), Y: 0, Z: 0}

	protein := &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, CA: atom1},
		},
		Atoms: []*parser.Atom{atom1},
	}

	quality, report := ScoreStructureQuality(protein)

	if quality > -999999.0 {
		t.Errorf("Expected very negative quality for invalid structure, got %.2f", quality)
	}

	if report.IsValid {
		t.Error("Should mark structure as invalid")
	}
}
