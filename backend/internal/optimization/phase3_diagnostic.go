// Phase 3 Diagnostic - Debug NaN RMSD issue
//
// Purpose: Find out why RMSD becomes NaN after Quaternion L-BFGS
//
package optimization

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/physics"
)

// DiagnoseLBFGS performs detailed diagnostic on L-BFGS step
func DiagnoseLBFGS(protein *parser.Protein) {
	fmt.Println("\n=== DIAGNOSTIC: L-BFGS STEP-BY-STEP ===")

	// Step 1: Extract current angles
	fmt.Println("Step 1: Extracting current dihedral angles...")
	angles := ExtractDihedrals(protein)
	fmt.Printf("  Extracted %d angle pairs\n", len(angles))

	// Print first 3 residues
	for i := 0; i < min(3, len(angles)); i++ {
		fmt.Printf("  Residue %d: φ=%.2f°, ψ=%.2f°\n",
			i, angles[i].Phi*180/math.Pi, angles[i].Psi*180/math.Pi)
	}

	// Step 2: Check current coordinates
	fmt.Println("\nStep 2: Checking current coordinates...")
	validPre := checkCoordinates(protein)
	if !validPre {
		fmt.Println("  ❌ Coordinates INVALID before SetDihedrals!")
		return
	}
	fmt.Println("  ✅ Coordinates valid before SetDihedrals")

	// Step 3: Try SetDihedrals with same angles (should be no-op)
	fmt.Println("\nStep 3: Testing SetDihedrals with same angles...")
	err := SetDihedrals(protein, angles)
	if err != nil {
		fmt.Printf("  ❌ SetDihedrals failed: %v\n", err)
		return
	}
	fmt.Println("  ✅ SetDihedrals succeeded")

	// Step 4: Check coordinates after SetDihedrals
	fmt.Println("\nStep 4: Checking coordinates after SetDihedrals...")
	validPost := checkCoordinates(protein)
	if !validPost {
		fmt.Println("  ❌ Coordinates INVALID after SetDihedrals!")
		printCoordinateDetails(protein)
		return
	}
	fmt.Println("  ✅ Coordinates still valid after SetDihedrals")

	// Step 5: Calculate energy
	fmt.Println("\nStep 5: Calculating energy...")
	config := DefaultQuaternionLBFGSConfig()
	energy := evaluateEnergyForProtein(protein, config)
	if math.IsNaN(energy) || math.IsInf(energy, 0) {
		fmt.Printf("  ❌ Energy is %v\n", energy)
		return
	}
	fmt.Printf("  ✅ Energy: %.2f kcal/mol\n", energy)

	// Step 6: Calculate gradient
	fmt.Println("\nStep 6: Calculating gradient...")
	gradient := computeDihedralGradient(protein, angles, config)
	gradNorm := vectorNormFloat(gradient)
	fmt.Printf("  Gradient norm: %.6f\n", gradNorm)
	if gradNorm < 1e-10 {
		fmt.Println("  ⚠️  WARNING: Gradient is nearly zero!")
		fmt.Println("     This means energy doesn't change with dihedral angles")
		fmt.Println("     Possible causes:")
		fmt.Println("     - Cutoffs too large (atoms too far apart)")
		fmt.Println("     - Energy function not sensitive to angles")
		fmt.Println("     - Finite difference delta too small")
	}

	// Step 7: Check if gradient contains NaN
	hasNaN := false
	for _, g := range gradient {
		if math.IsNaN(g) {
			hasNaN = true
			break
		}
	}
	if hasNaN {
		fmt.Println("  ❌ Gradient contains NaN!")
		return
	}
	fmt.Println("  ✅ Gradient is valid (no NaN)")

	fmt.Println("\n=== DIAGNOSTIC COMPLETE ===")
}

// checkCoordinates returns true if all coordinates are valid
func checkCoordinates(protein *parser.Protein) bool {
	for _, res := range protein.Residues {
		atoms := []*parser.Atom{res.N, res.CA, res.C, res.O}
		for _, atom := range atoms {
			if atom == nil {
				continue
			}
			if math.IsNaN(atom.X) || math.IsNaN(atom.Y) || math.IsNaN(atom.Z) {
				return false
			}
			r := math.Sqrt(atom.X*atom.X + atom.Y*atom.Y + atom.Z*atom.Z)
			if r > 1000.0 {
				return false
			}
		}
	}
	return true
}

// printCoordinateDetails prints detailed info about invalid coordinates
func printCoordinateDetails(protein *parser.Protein) {
	fmt.Println("\n  Coordinate Details:")
	for i, res := range protein.Residues {
		if res.N != nil {
			fmt.Printf("    Res %d N:  (%.2f, %.2f, %.2f)\n", i, res.N.X, res.N.Y, res.N.Z)
		}
		if res.CA != nil {
			fmt.Printf("    Res %d CA: (%.2f, %.2f, %.2f)\n", i, res.CA.X, res.CA.Y, res.CA.Z)
		}
		if res.C != nil {
			fmt.Printf("    Res %d C:  (%.2f, %.2f, %.2f)\n", i, res.C.X, res.C.Y, res.C.Z)
		}
	}
}

// min returns minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DiagnoseEnergyGradient checks why gradient might be zero
func DiagnoseEnergyGradient(protein *parser.Protein, config QuaternionLBFGSConfig) {
	fmt.Println("\n=== ENERGY GRADIENT DIAGNOSTIC ===")

	// Check current energy
	energy := physics.CalculateTotalEnergy(protein, config.VdWCutoff, config.ElecCutoff)
	fmt.Printf("Total Energy: %.2f kcal/mol\n", energy.Total)
	fmt.Printf("  VdW:        %.2f kcal/mol\n", energy.VanDerWaals)
	fmt.Printf("  Elec:       %.2f kcal/mol\n", energy.Electrostatic)
	fmt.Printf("  Bond:       %.2f kcal/mol\n", energy.Bond)
	fmt.Printf("  Angle:      %.2f kcal/mol\n", energy.Angle)
	fmt.Printf("  Dihedral:   %.2f kcal/mol\n", energy.Dihedral)

	// Check if energy components are zero
	if math.Abs(energy.Total) < 1e-6 {
		fmt.Println("\n⚠️  Total energy is essentially zero!")
		fmt.Println("   This suggests no interactions are being calculated")
	}

	// Check number of atom pairs within cutoff
	atomCount := len(protein.Atoms)
	fmt.Printf("\nTotal atoms: %d\n", atomCount)
	fmt.Printf("VdW cutoff: %.1f Å\n", config.VdWCutoff)
	fmt.Printf("Elec cutoff: %.1f Å\n", config.ElecCutoff)

	// Count pairs within cutoff
	pairsWithinVdW := 0
	pairsWithinElec := 0
	for i := 0; i < len(protein.Atoms); i++ {
		for j := i + 1; j < len(protein.Atoms); j++ {
			a1 := protein.Atoms[i]
			a2 := protein.Atoms[j]
			dx := a2.X - a1.X
			dy := a2.Y - a1.Y
			dz := a2.Z - a1.Z
			r := math.Sqrt(dx*dx + dy*dy + dz*dz)

			if r < config.VdWCutoff {
				pairsWithinVdW++
			}
			if r < config.ElecCutoff {
				pairsWithinElec++
			}
		}
	}

	fmt.Printf("\nAtom pairs within VdW cutoff: %d\n", pairsWithinVdW)
	fmt.Printf("Atom pairs within Elec cutoff: %d\n", pairsWithinElec)

	if pairsWithinVdW == 0 {
		fmt.Println("\n❌ NO ATOM PAIRS within VdW cutoff!")
		fmt.Println("   This explains why gradient is zero")
		fmt.Println("   Atoms might be too far apart (extended structure)")
	}

	fmt.Println("\n=== DIAGNOSTIC COMPLETE ===")
}
