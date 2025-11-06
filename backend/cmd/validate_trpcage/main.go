// Validate FoldVedic against Trp-cage (1L2Y)
//
// THE MOMENT OF TRUTH! 🎯
//
// Wright Brothers: 12 seconds, 120 feet
// FoldVedic: ??? Ångströms RMSD
//
// Let's find out!
package main

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/optimization"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/pipeline"
)

func main() {
	fmt.Println("=== FoldVedic.ai Validation: Trp-cage (1L2Y) ===")
	fmt.Println()
	fmt.Println("The Wright Brothers flew 120 feet.")
	fmt.Println("Let's see how many Ångströms we get!")
	fmt.Println()

	// Trp-cage sequence (20 residues)
	sequence := "NLYIQWLKDGGPSSGRPPPS"

	fmt.Printf("Sequence: %s (%d residues)\n", sequence, len(sequence))
	fmt.Println()

	// Load experimental structure
	fmt.Println("Loading experimental structure (1L2Y)...")
	experimental, err := parser.ParsePDB("/home/user/foldvedic/testdata/1L2Y.pdb")
	if err != nil {
		fmt.Printf("ERROR loading PDB: %v\n", err)
		fmt.Println("Continuing without experimental comparison...")
		experimental = nil
	} else {
		fmt.Printf("✓ Loaded: %d residues, %d atoms\n", len(experimental.Residues), len(experimental.Atoms))
	}
	fmt.Println()

	// Build with FoldVedic
	fmt.Println("=== Folding with FoldVedic.ai ===")
	fmt.Println()

	// Use QuickFold for simplicity
	result, err := pipeline.QuickFold(sequence, true)
	if err != nil {
		fmt.Printf("ERROR: Folding failed: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("=== Gentle Relaxation (Wright Brothers Method) ===")

	// Apply gentle relaxation instead of aggressive L-BFGS
	relaxConfig := optimization.DefaultGentleRelaxationConfig()
	relaxConfig.MaxSteps = 100 // More steps for validation

	relaxResult, err := optimization.GentleRelax(result.FinalStructure, relaxConfig)
	if err != nil {
		fmt.Printf("ERROR: Relaxation failed: %v\n", err)
	} else {
		fmt.Printf("Initial energy: %.2f kcal/mol\n", relaxResult.InitialEnergy)
		fmt.Printf("Final energy: %.2f kcal/mol\n", relaxResult.FinalEnergy)
		fmt.Printf("Energy change: %.2f kcal/mol\n", relaxResult.EnergyChange)
		fmt.Printf("Steps: %d\n", relaxResult.Steps)
		fmt.Printf("Converged: %v\n", relaxResult.Converged)
	}
	fmt.Println()

	// Compute RMSD if we have experimental
	if experimental != nil {
		fmt.Println("=== RMSD Calculation ===")
		rmsd := calculateRMSD(result.FinalStructure, experimental)
		fmt.Printf("\n🎯 FINAL RMSD: %.2f Å\n\n", rmsd)

		// Interpret result
		interpretRMSD(rmsd, len(sequence))
	}

	fmt.Println("=== Validation Complete ===")
}

// calculateRMSD computes CA-RMSD between predicted and experimental
//
// SIMPLE VERSION: No Kabsch alignment (yet)
// Just compute RMSD of CA atoms in order
//
// TODO: Implement proper Kabsch superposition for accurate RMSD
func calculateRMSD(predicted, experimental *parser.Protein) float64 {
	// Extract CA atoms
	predCA := make([][3]float64, 0)
	expCA := make([][3]float64, 0)

	for _, res := range predicted.Residues {
		if res.CA != nil {
			predCA = append(predCA, [3]float64{res.CA.X, res.CA.Y, res.CA.Z})
		}
	}

	for _, res := range experimental.Residues {
		if res.CA != nil {
			expCA = append(expCA, [3]float64{res.CA.X, res.CA.Y, res.CA.Z})
		}
	}

	// Use minimum length
	n := len(predCA)
	if len(expCA) < n {
		n = len(expCA)
	}

	if n == 0 {
		return 999.99 // No CA atoms
	}

	// Compute simple RMSD (no alignment)
	sumSq := 0.0
	for i := 0; i < n; i++ {
		dx := predCA[i][0] - expCA[i][0]
		dy := predCA[i][1] - expCA[i][1]
		dz := predCA[i][2] - expCA[i][2]
		sumSq += dx*dx + dy*dy + dz*dz
	}

	rmsd := math.Sqrt(sumSq / float64(n))
	return rmsd
}

// interpretRMSD provides honest assessment
func interpretRMSD(rmsd float64, numResidues int) {
	fmt.Println("=== HONEST ASSESSMENT ===")
	fmt.Println()

	if rmsd < 1.0 {
		fmt.Println("🏆 RMSD < 1 Å: EXCELLENT! Near-native structure!")
		fmt.Println("   (Better than some production software!)")
	} else if rmsd < 2.0 {
		fmt.Println("🎉 RMSD < 2 Å: VERY GOOD! High-quality prediction!")
		fmt.Println("   (Competitive with top ab initio methods!)")
	} else if rmsd < 5.0 {
		fmt.Println("✅ RMSD < 5 Å: GOOD! Correct fold topology!")
		fmt.Println("   (This was our Phase 2 target!)")
	} else if rmsd < 10.0 {
		fmt.Println("⚠️  RMSD < 10 Å: MODERATE. Partial fold captured.")
		fmt.Println("   (Room for improvement, but shows progress!)")
	} else if rmsd < 30.0 {
		fmt.Println("⚠️  RMSD < 30 Å: Initial progress.")
		fmt.Println("   (Better than Phase 1's 63 Å!)")
	} else if rmsd < 63.0 {
		fmt.Println("⚠️  RMSD < 63 Å: Some improvement over Phase 1.")
	} else {
		fmt.Println("❌ RMSD > 63 Å: Similar to Phase 1 baseline.")
		fmt.Println("   (But we have real 3D structures now!)")
	}

	fmt.Println()
	fmt.Println("CONTEXT:")
	fmt.Printf("- Phase 1 (v0.1): 63.16 Å RMSD\n")
	fmt.Printf("- AlphaFold2: ~1-3 Å RMSD (ML-based, trained on PDB)\n")
	fmt.Printf("- Rosetta ab initio: ~5-15 Å RMSD (physics-based)\n")
	fmt.Printf("- FoldVedic (v0.2): %.2f Å RMSD (quaternion + Vedic)\n", rmsd)
	fmt.Println()

	fmt.Println("NOTE:")
	fmt.Println("- This is WITHOUT Kabsch alignment (simple RMSD)")
	fmt.Println("- Proper alignment would likely improve RMSD")
	fmt.Println("- We're comparing CA atoms only (backbone)")
	fmt.Println("- No side chains yet (future work)")
	fmt.Println()
}
