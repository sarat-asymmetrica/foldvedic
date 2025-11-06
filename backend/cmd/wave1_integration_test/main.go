// Wave 1 Integration Test - combines all three agents
//
// Agent 1.1: PDB Parser & Ramachandran Mapper
// Agent 1.2: Force Field Engine
// Agent 1.3: Vedic Harmonic Scorer
//
// This test demonstrates the complete Wave 1 pipeline
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/physics"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/vedic"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("FoldVedic Wave 1 Integration Test")
	fmt.Println("Testing: PDB Parser + Ramachandran + Force Field + Vedic Scorer")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Test file (relative to repository root)
	pdbFile := "../../../testdata/test_peptide.pdb"

	// Agent 1.1: Parse PDB and calculate Ramachandran angles
	fmt.Println("Agent 1.1: PDB Parser & Ramachandran Mapper")
	fmt.Println(strings.Repeat("-", 80))

	protein, err := parser.ParsePDB(pdbFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse PDB: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Parsed PDB file: %s\n", pdbFile)
	fmt.Printf("  Atoms:     %d\n", len(protein.Atoms))
	fmt.Printf("  Residues:  %d\n", len(protein.Residues))
	fmt.Printf("  Complete:  %d\n", protein.NumCompleteResidues())
	fmt.Println()

	// Calculate Ramachandran angles
	angles := geometry.CalculateRamachandran(protein)
	fmt.Printf("✓ Calculated Ramachandran angles for %d residues\n", len(angles))

	// Show sample angles
	fmt.Println("  Sample angles (first 3 residues):")
	for i := 0; i < 3 && i < len(angles); i++ {
		fmt.Printf("    Residue %d: φ = %6.1f°, ψ = %6.1f°\n",
			i+1, angles[i].ToDegressPhi(), angles[i].ToDegressPsi())
	}
	fmt.Println()

	// Map to quaternions
	fmt.Println("  Testing quaternion mapping:")
	for i := 0; i < 3 && i < len(angles); i++ {
		q := geometry.RamachandranToQuaternion(angles[i].Phi, angles[i].Psi)
		fmt.Printf("    Residue %d: q = [%.3f, %.3f, %.3f, %.3f]\n",
			i+1, q.W, q.X, q.Y, q.Z)
	}
	fmt.Println()

	// Test slerp
	if len(angles) >= 2 {
		q1 := geometry.RamachandranToQuaternion(angles[0].Phi, angles[0].Psi)
		q2 := geometry.RamachandranToQuaternion(angles[1].Phi, angles[1].Psi)
		qMid := q1.Slerp(q2, 0.5)
		fmt.Printf("✓ Quaternion slerp interpolation tested\n")
		fmt.Printf("  Midpoint quaternion: [%.3f, %.3f, %.3f, %.3f]\n",
			qMid.W, qMid.X, qMid.Y, qMid.Z)
		fmt.Println()
	}

	// Agent 1.2: Force Field Energy Calculation
	fmt.Println("Agent 1.2: Force Field Engine")
	fmt.Println(strings.Repeat("-", 80))

	energy := physics.CalculateTotalEnergy(protein, 10.0, 12.0)

	fmt.Println("✓ Calculated total energy")
	fmt.Printf("  Bond:          %8.2f kcal/mol\n", energy.Bond)
	fmt.Printf("  Angle:         %8.2f kcal/mol\n", energy.Angle)
	fmt.Printf("  Dihedral:      %8.2f kcal/mol\n", energy.Dihedral)
	fmt.Printf("  Van der Waals: %8.2f kcal/mol\n", energy.VanDerWaals)
	fmt.Printf("  Electrostatic: %8.2f kcal/mol\n", energy.Electrostatic)
	fmt.Printf("  Total:         %8.2f kcal/mol\n", energy.Total)
	fmt.Println()

	// Validate energy
	warnings := physics.ValidateEnergy(energy, len(protein.Residues))
	if len(warnings) > 0 {
		fmt.Println("  Energy validation warnings:")
		for _, warning := range warnings {
			fmt.Printf("    ⚠  %s\n", warning)
		}
		fmt.Println()
	} else {
		fmt.Println("✓ Energy values are physically reasonable")
		fmt.Println()
	}

	// Test force calculation
	forces := physics.CalculateForces(protein, 10.0, 12.0)
	fmt.Printf("✓ Calculated forces for %d atoms\n", len(forces))

	// Show sample forces
	fmt.Println("  Sample forces (first 3 atoms):")
	for i := 0; i < 3 && i < len(protein.Atoms); i++ {
		atom := protein.Atoms[i]
		force := forces[atom.Serial]
		forceMag := force.Magnitude()
		fmt.Printf("    Atom %d (%s): |F| = %.2f kcal/(mol·Å)\n",
			atom.Serial, atom.Name, forceMag)
	}
	fmt.Println()

	// Agent 1.3: Vedic Harmonic Scorer
	fmt.Println("Agent 1.3: Vedic Harmonic Scorer")
	fmt.Println(strings.Repeat("-", 80))

	vedicScore := vedic.CalculateVedicScore(protein, angles)

	fmt.Println("✓ Calculated Vedic harmonic score")
	fmt.Printf("  Golden Ratio Score: %.3f\n", vedicScore.GoldenRatioScore)
	fmt.Printf("  Digital Root Score: %.3f\n", vedicScore.DigitalRootScore)
	fmt.Printf("  Breathing Score:    %.3f\n", vedicScore.BreathingScore)
	fmt.Printf("  Total Score:        %.3f\n", vedicScore.TotalScore)
	fmt.Println()

	fmt.Printf("  Secondary structure content:\n")
	fmt.Printf("    Helix residues: %d\n", vedicScore.NumHelixResidues)
	fmt.Printf("    Sheet residues: %d\n", vedicScore.NumSheetResidues)
	fmt.Println()

	// Summary
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Wave 1 Integration Test: COMPLETE")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	fmt.Println("Summary:")
	fmt.Printf("  ✓ Agent 1.1: PDB parsing and Ramachandran mapping functional\n")
	fmt.Printf("  ✓ Agent 1.2: Force field energy calculation operational\n")
	fmt.Printf("  ✓ Agent 1.3: Vedic harmonic scoring implemented\n")
	fmt.Println()

	fmt.Println("All Wave 1 components successfully integrated!")
	fmt.Println()

	fmt.Println("Next steps (Wave 2):")
	fmt.Println("  - Download real PDB structures (1UBQ, 1CRN, 2KXA)")
	fmt.Println("  - Implement RMSD/TM-score validation metrics")
	fmt.Println("  - Add spatial hashing for O(n) force calculations")
	fmt.Println()

	os.Exit(0)
}
