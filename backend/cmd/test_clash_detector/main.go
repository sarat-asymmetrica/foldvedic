package main

import (
	"fmt"
	"log"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/physics"
)

func main() {
	fmt.Println("=== Clash Detector Validation Test ===")
	fmt.Println()

	// Test 1: Valid structure (should pass all checks)
	fmt.Println("Test 1: Valid Structure")
	validProtein := createValidProtein()
	quality1, report1 := physics.ScoreStructureQuality(validProtein)
	fmt.Printf("  Valid: %v\n", report1.IsValid)
	fmt.Printf("  Clashes: %d\n", report1.ClashCount)
	fmt.Printf("  Quality Score: %.3f\n", quality1)
	if !report1.IsValid || report1.HasClashes {
		log.Fatal("❌ Valid structure should pass all checks")
	}
	fmt.Println("  ✅ PASS")
	fmt.Println()

	// Test 2: Structure with clashes
	fmt.Println("Test 2: Structure with Severe Clashes")
	clashProtein := createClashProtein()
	quality2, report2 := physics.ScoreStructureQuality(clashProtein)
	fmt.Printf("  Valid: %v\n", report2.IsValid)
	fmt.Printf("  Clashes: %d\n", report2.ClashCount)
	fmt.Printf("  Worst Clash: %.2f Å\n", report2.WorstClashDist)
	fmt.Printf("  Quality Score: %.3f\n", quality2)
	if !report2.HasClashes {
		log.Fatal("❌ Clash structure should detect clashes")
	}
	fmt.Println("  ✅ PASS")
	fmt.Println()

	// Test 3: Structure with broken backbone
	fmt.Println("Test 3: Structure with Broken Backbone")
	brokenProtein := createBrokenBackboneProtein()
	_, report3 := physics.ScoreStructureQuality(brokenProtein)
	fmt.Printf("  Valid: %v\n", report3.IsValid)
	if report3.IsValid {
		log.Fatal("❌ Broken backbone structure should fail validation")
	}
	fmt.Printf("  Error: %s\n", report3.ValidationError)
	fmt.Println("  ✅ PASS")
	fmt.Println()

	// Test 4: Energy capping
	fmt.Println("Test 4: Energy Capping")
	clashProteinCapped := createClashProtein()
	energy := physics.CalculateTotalEnergy(clashProteinCapped, 8.0, 12.0)
	fmt.Printf("  Total Energy: %.2f kcal/mol\n", energy.Total)
	if energy.Total > 10000.0 {
		log.Fatal("❌ Energy should be capped at 10,000 kcal/mol")
	}
	if energy.Total < -10000.0 {
		log.Fatal("❌ Energy should be capped at -10,000 kcal/mol")
	}
	fmt.Println("  ✅ Energy properly capped")
	fmt.Println()

	fmt.Println("=== ALL TESTS PASSED ===")
}

func createValidProtein() *parser.Protein {
	// Create a simple 3-residue protein with proper spacing
	atomN1 := &parser.Atom{Name: "N", Element: "N", ResSeq: 1, X: 0, Y: 0, Z: 0}
	atomCA1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, X: 1.5, Y: 0, Z: 0}
	atomC1 := &parser.Atom{Name: "C", Element: "C", ResSeq: 1, X: 3.0, Y: 0, Z: 0}
	atomN2 := &parser.Atom{Name: "N", Element: "N", ResSeq: 2, X: 4.3, Y: 0, Z: 0}
	atomCA2 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 2, X: 5.8, Y: 0, Z: 0}
	atomC2 := &parser.Atom{Name: "C", Element: "C", ResSeq: 2, X: 7.3, Y: 0, Z: 0}
	atomN3 := &parser.Atom{Name: "N", Element: "N", ResSeq: 3, X: 8.6, Y: 0, Z: 0}
	atomCA3 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 3, X: 10.1, Y: 0, Z: 0}
	atomC3 := &parser.Atom{Name: "C", Element: "C", ResSeq: 3, X: 11.6, Y: 0, Z: 0}

	return &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", N: atomN1, CA: atomCA1, C: atomC1},
			{Name: "GLY", SeqNum: 2, ChainID: "A", N: atomN2, CA: atomCA2, C: atomC2},
			{Name: "VAL", SeqNum: 3, ChainID: "A", N: atomN3, CA: atomCA3, C: atomC3},
		},
		Atoms: []*parser.Atom{atomN1, atomCA1, atomC1, atomN2, atomCA2, atomC2, atomN3, atomCA3, atomC3},
	}
}

func createClashProtein() *parser.Protein {
	// Create a protein with severe clashes (atoms too close)
	atomCA1 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 1, ChainID: "A", X: 0, Y: 0, Z: 0}
	atomCA2 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 3, ChainID: "A", X: 0.5, Y: 0, Z: 0} // Too close!
	atomCA3 := &parser.Atom{Name: "CA", Element: "C", ResSeq: 5, ChainID: "A", X: 0.6, Y: 0, Z: 0} // Too close!

	return &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", CA: atomCA1},
			{Name: "GLY", SeqNum: 3, ChainID: "A", CA: atomCA2},
			{Name: "VAL", SeqNum: 5, ChainID: "A", CA: atomCA3},
		},
		Atoms: []*parser.Atom{atomCA1, atomCA2, atomCA3},
	}
}

func createBrokenBackboneProtein() *parser.Protein {
	// Create a protein with broken peptide bond
	atomC1 := &parser.Atom{Name: "C", Element: "C", ResSeq: 1, X: 0, Y: 0, Z: 0}
	atomN2 := &parser.Atom{Name: "N", Element: "N", ResSeq: 2, X: 10, Y: 0, Z: 0} // Way too far!

	return &parser.Protein{
		Residues: []*parser.Residue{
			{Name: "ALA", SeqNum: 1, ChainID: "A", C: atomC1},
			{Name: "GLY", SeqNum: 2, ChainID: "A", N: atomN2},
		},
		Atoms: []*parser.Atom{atomC1, atomN2},
	}
}
