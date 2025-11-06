// Quick energy diagnostic - is 2212 kcal/mol realistic?
package main

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/optimization"
)

func main() {
	fmt.Println("=== Energy Diagnostic ===")
	fmt.Println()

	// Build tiny protein with quaternion builder
	sequence := "ACDEF"
	angles := make([]geometry.RamachandranAngles, len(sequence))
	for i := range angles {
		// Helix conformation
		angles[i] = geometry.RamachandranAngles{
			Phi: -60.0 * math.Pi / 180.0,
			Psi: -45.0 * math.Pi / 180.0,
		}
	}

	protein, err := geometry.BuildProteinFromAngles(sequence, angles)
	if err != nil {
		fmt.Printf("Build failed: %v\n", err)
		return
	}

	fmt.Printf("Protein built: %d residues, %d atoms\n", len(protein.Residues), len(protein.Atoms))
	fmt.Println()

	// Try to compute energy
	config := optimization.DefaultAdaptiveOptimizationConfig()
	result, err := optimization.OptimizeProtein(protein, config)

	if err != nil {
		fmt.Printf("Optimization failed: %v\n", err)
	} else {
		fmt.Printf("Initial energy: %.2f kcal/mol\n", result.InitialEnergy)
		fmt.Printf("Final energy: %.2f kcal/mol\n", result.FinalEnergy)
		fmt.Printf("Energy change: %.2f kcal/mol\n", result.EnergyChange)
		fmt.Printf("Iterations: %d\n", result.Iterations)
		fmt.Printf("Strategy: %v\n", result.Strategy)
	}

	fmt.Println()

	// Energy sanity check
	// Typical protein: -5 to -15 kcal/mol per residue
	// Folded proteins: strongly negative (favorable)
	// Unfolded/misfolded: weakly negative or positive

	energyPerResidue := result.FinalEnergy / float64(len(sequence))
	fmt.Printf("Energy per residue: %.2f kcal/mol\n", energyPerResidue)
	fmt.Println()

	if result.FinalEnergy > 1e9 {
		fmt.Println("❌ Energy is placeholder (>1e9)")
	} else if result.FinalEnergy > 10000 {
		fmt.Println("⚠️  Energy is very high (likely clashes or bad geometry)")
	} else if result.FinalEnergy > 1000 {
		fmt.Println("⚠️  Energy is high (unfolded or poor structure)")
	} else if result.FinalEnergy > 0 {
		fmt.Println("⚠️  Energy is positive (unfavorable, likely incorrect)")
	} else if result.FinalEnergy > -100*float64(len(sequence)) {
		fmt.Println("✓ Energy is reasonable (weak to moderate favorability)")
	} else {
		fmt.Println("✓ Energy is strongly negative (well-folded)")
	}
}
