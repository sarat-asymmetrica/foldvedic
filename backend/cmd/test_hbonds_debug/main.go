package main

import (
	"fmt"
	"log"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

func main() {
	fmt.Println("=== Detailed H-Bond Geometry Check ===")

	// Load PDB
	protein, err := parser.ParsePDB("testdata/1L2Y.pdb")
	if err != nil {
		log.Fatalf("Failed to load PDB: %v", err)
	}

	// Build H atom map
	hAtomMap := make(map[int]*parser.Atom)
	for _, atom := range protein.Atoms {
		if atom.Element == "H" && (atom.Name == "H" || atom.Name == "HN") {
			hAtomMap[atom.ResSeq] = atom
			fmt.Printf("H atom for residue %d: (%.2f, %.2f, %.2f)\n",
				atom.ResSeq, atom.X, atom.Y, atom.Z)
		}
	}

	fmt.Printf("\nFound %d backbone H atoms\n\n", len(hAtomMap))

	// Check a few specific pairs
	testPairs := [][2]int{
		{6, 2},  // i+4 helix pattern
		{7, 3},  // i+4
		{8, 4},  // i+4
		{13, 9}, // possible longer-range
	}

	for _, pair := range testPairs {
		donorSeq := pair[0]
		acceptorSeq := pair[1]

		// Find atoms
		var donorN, donorH, acceptorO *parser.Atom
		for _, atom := range protein.Atoms {
			if atom.ResSeq == donorSeq && atom.Name == "N" {
				donorN = atom
			}
			if atom.ResSeq == acceptorSeq && atom.Name == "O" {
				acceptorO = atom
			}
		}
		donorH = hAtomMap[donorSeq]

		if donorN == nil || donorH == nil || acceptorO == nil {
			fmt.Printf("Residue %d → %d: Missing atoms (N=%v, H=%v, O=%v)\n",
				donorSeq, acceptorSeq,
				donorN != nil, donorH != nil, acceptorO != nil)
			continue
		}

		// Calculate H···O distance
		dx := donorH.X - acceptorO.X
		dy := donorH.Y - acceptorO.Y
		dz := donorH.Z - acceptorO.Z
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

		// Calculate N-H···O angle
		// Angle is measured at H in the N-H-O chain
		// Vector from H to N (backwards)
		hnX := donorN.X - donorH.X
		hnY := donorN.Y - donorH.Y
		hnZ := donorN.Z - donorH.Z
		hnLen := math.Sqrt(hnX*hnX + hnY*hnY + hnZ*hnZ)

		// Vector from H to O (forward)
		hoX := acceptorO.X - donorH.X
		hoY := acceptorO.Y - donorH.Y
		hoZ := acceptorO.Z - donorH.Z
		hoLen := math.Sqrt(hoX*hoX + hoY*hoY + hoZ*hoZ)

		// Dot product
		dotProduct := (hnX*hoX + hnY*hoY + hnZ*hoZ) / (hnLen * hoLen)
		if dotProduct > 1.0 {
			dotProduct = 1.0
		}
		if dotProduct < -1.0 {
			dotProduct = -1.0
		}

		angleRad := math.Acos(dotProduct)
		angleDeg := angleRad * 180.0 / math.Pi

		fmt.Printf("Residue %d → %d:\n", donorSeq, acceptorSeq)
		fmt.Printf("  H···O distance: %.2f Å (need: 1.5-2.5 Å)\n", dist)
		fmt.Printf("  N-H···O angle:  %.1f° (need: >120°)\n", angleDeg)

		passDistance := dist >= 1.5 && dist <= 2.5
		passAngle := angleDeg >= 120.0

		if passDistance && passAngle {
			fmt.Printf("  ✅ VALID H-BOND!\n")
		} else {
			if !passDistance {
				fmt.Printf("  ❌ Distance out of range\n")
			}
			if !passAngle {
				fmt.Printf("  ❌ Angle too small\n")
			}
		}
		fmt.Println()
	}
}
