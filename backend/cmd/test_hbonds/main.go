package main

import (
	"fmt"
	"log"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/physics"
)

func main() {
	fmt.Println("=== H-Bond Detection Test ===")

	// Load PDB
	protein, err := parser.ParsePDB("testdata/1L2Y.pdb")
	if err != nil {
		log.Fatalf("Failed to load PDB: %v", err)
	}

	fmt.Printf("Loaded %d residues\n", len(protein.Residues))
	fmt.Printf("Total atoms: %d\n", len(protein.Atoms))

	// Count H atoms
	hCount := 0
	backboneH := 0
	for _, atom := range protein.Atoms {
		if atom.Element == "H" {
			hCount++
			if atom.Name == "H" || atom.Name == "HN" {
				backboneH++
				fmt.Printf("  Found backbone H: residue %d (%s) at (%.2f, %.2f, %.2f)\n",
					atom.ResSeq, atom.ResName, atom.X, atom.Y, atom.Z)
			}
		}
	}

	fmt.Printf("\nTotal H atoms: %d\n", hCount)
	fmt.Printf("Backbone H atoms (H or HN): %d\n", backboneH)
	fmt.Println()

	// Detect H-bonds
	hbonds := physics.DetectHydrogenBonds(protein)
	fmt.Printf("Detected H-bonds: %d\n", len(hbonds))

	for i, hb := range hbonds {
		fmt.Printf("  H-bond %d: %s%d → %s%d, dist=%.2f Å, angle=%.1f°, E=%.2f kcal/mol\n",
			i+1,
			hb.DonorResidue.Name, hb.DonorResidue.SeqNum,
			hb.AcceptorResidue.Name, hb.AcceptorResidue.SeqNum,
			hb.Distance, hb.Angle, hb.Energy)
	}

	// Manual check: Try one pair manually
	fmt.Println("\n=== Manual H-bond check ===")
	// Check residue 6 (TRP) H → residue 2 (LEU) O (typical i→i+4 helix H-bond)
	var res2O, res6N, res6H *parser.Atom
	for _, atom := range protein.Atoms {
		if atom.ResSeq == 2 && atom.Name == "O" {
			res2O = atom
		}
		if atom.ResSeq == 6 && atom.Name == "N" {
			res6N = atom
		}
		if atom.ResSeq == 6 && atom.Name == "H" {
			res6H = atom
		}
	}

	if res2O != nil && res6H != nil && res6N != nil {
		dx := res6H.X - res2O.X
		dy := res6H.Y - res2O.Y
		dz := res6H.Z - res2O.Z
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

		fmt.Printf("Residue 6 (TRP) H → Residue 2 (LEU) O:\n")
		fmt.Printf("  H position: (%.2f, %.2f, %.2f)\n", res6H.X, res6H.Y, res6H.Z)
		fmt.Printf("  O position: (%.2f, %.2f, %.2f)\n", res2O.X, res2O.Y, res2O.Z)
		fmt.Printf("  H···O distance: %.2f Å (should be 1.5-2.5 Å for H-bond)\n", dist)
	}
}
