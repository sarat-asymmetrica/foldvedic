// Diagnostic program for Phase 3 NaN issue
package main

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/optimization"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  🔍 PHASE 3 NaN DIAGNOSTIC                                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Build simple 5-residue peptide
	sequence := "AAAAA"
	fmt.Printf("Building test protein: %s (%d residues)\n", sequence, len(sequence))

	angles := make([]geometry.RamachandranAngles, len(sequence))
	for i := range angles {
		angles[i] = geometry.RamachandranAngles{
			Phi: -120.0 * math.Pi / 180.0,
			Psi: +120.0 * math.Pi / 180.0,
		}
	}

	protein, err := geometry.BuildProteinFromAngles(sequence, angles)
	if err != nil {
		fmt.Printf("❌ Failed to build protein: %v\n", err)
		return
	}

	fmt.Printf("✅ Built protein: %d residues, %d atoms\n\n", len(protein.Residues), len(protein.Atoms))

	// Run diagnostic
	optimization.DiagnoseLBFGS(protein)

	// Run energy gradient diagnostic
	config := optimization.DefaultQuaternionLBFGSConfig()
	optimization.DiagnoseEnergyGradient(protein, config)
}
