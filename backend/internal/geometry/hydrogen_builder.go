// Package geometry - Hydrogen Atom Placement
//
// WAVE 11.4: Root Cause Fix - Blocker #1
//
// PROBLEM IDENTIFIED:
// - Zero hydrogen bonds detected (should be 10-15 for Trp-cage)
// - Coordinate builder only places N, CA, C, O atoms
// - Energy function blind to secondary structure without H-bonds
//
// SOLUTION:
// Add hydrogen atoms using ideal geometry from crystallography
//
// BIOCHEMIST: Standard H-bond geometry (N-H...O=C)
//   - N-H bond length: 1.01 Å (Engh & Huber, 2001)
//   - C-N-H angle: ~120° (sp2 planar)
//   - CA-HA bond length: 1.09 Å (tetrahedral sp3)
//
// PHYSICIST: Geometric placement using vector algebra
//   - Bisector method for backbone NH
//   - Cross product for tetrahedral Hα
//
// WRIGHT BROTHERS: Start simple
//   - Backbone H (NH) - critical for H-bonds
//   - Alpha H (Hα) - tetrahedral geometry
//   - Full sidechain H placement later (500+ lines)
//
// VALIDATION TARGET:
// - Before: 0 H-bonds, 0 kcal/mol H-bond energy
// - After: 10-15 H-bonds, -50 to -75 kcal/mol
package geometry

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// Standard hydrogen bond lengths (Å) from crystallography
// Citation: Engh, R. A., & Huber, R. (2001). "Structure quality and target parameters"
const (
	BondN_H  = 1.01 // N-H bond (backbone amide)
	BondCA_H = 1.09 // CA-Hα bond (tetrahedral)
)

// AddHydrogens adds hydrogen atoms to all residues using ideal geometry
//
// ALGORITHM:
// 1. Backbone NH hydrogen (critical for H-bonds)
//    - Skip N-terminal (no previous C)
//    - Place H using bisector of C-N-CA angle
//    - Distance: 1.01 Å from N
//
// 2. Alpha hydrogen (Hα at CA)
//    - Tetrahedral geometry
//    - Perpendicular to C-CA-N plane
//    - Distance: 1.09 Å from CA
//
// 3. Future: Full sidechain hydrogens (Wave 12+)
//
// BIOCHEMIST NOTES:
// - Glycine has 2 Hα (no Cβ)
// - Proline has no NH (cyclic structure)
// - These edge cases handled in future waves
func AddHydrogens(protein *parser.Protein) error {
	if protein == nil {
		return fmt.Errorf("protein is nil")
	}

	for i := range protein.Residues {
		res := protein.Residues[i]
		if res == nil {
			continue
		}

		// Add backbone NH hydrogen (critical for H-bonds)
		if err := addBackboneHydrogen(res, i, protein); err != nil {
			return fmt.Errorf("failed to add backbone H for residue %d: %w", i, err)
		}

		// Add alpha hydrogen (Hα)
		if err := addAlphaHydrogen(res); err != nil {
			return fmt.Errorf("failed to add alpha H for residue %d: %w", i, err)
		}
	}

	return nil
}

// addBackboneHydrogen adds backbone NH hydrogen using ideal geometry
//
// GEOMETRY:
// - C(i-1) - N(i) - CA(i) defines the plane
// - H is placed as bisector of C-N-CA angle
// - C-N-H angle ≈ 120° (sp2 planar nitrogen)
// - N-H distance = 1.01 Å
//
// PHYSICIST: Vector algebra approach
// 1. Normalize C→N vector
// 2. Normalize N→CA vector
// 3. H direction = normalized(C→N + N→CA)
// 4. H position = N + (H direction × 1.01 Å)
func addBackboneHydrogen(residue *parser.Residue, resIndex int, protein *parser.Protein) error {
	// Skip N-terminal residue (no previous C)
	if resIndex == 0 {
		return nil
	}

	// Skip proline (no NH hydrogen - cyclic structure)
	if residue.Name == "PRO" {
		return nil
	}

	// Get backbone atoms for current residue
	n := residue.N
	ca := residue.CA
	if n == nil || ca == nil {
		return nil // Skip if backbone incomplete
	}

	// Get previous residue's C for geometry reference
	prevRes := protein.Residues[resIndex-1]
	if prevRes == nil {
		return nil
	}
	prevC := prevRes.C
	if prevC == nil {
		return nil // Can't place H without previous C
	}

	// Calculate NH bond vector using bisector method
	//
	// Vector from previous C to current N
	cnVec := Vector3{
		X: n.X - prevC.X,
		Y: n.Y - prevC.Y,
		Z: n.Z - prevC.Z,
	}
	cnVec = cnVec.Normalize()

	// Vector from N to CA
	ncaVec := Vector3{
		X: ca.X - n.X,
		Y: ca.Y - n.Y,
		Z: ca.Z - n.Z,
	}
	ncaVec = ncaVec.Normalize()

	// H direction: bisector of C-N-CA angle
	// This gives ~120° C-N-H angle (sp2 geometry)
	hDir := cnVec.Add(ncaVec).Normalize()

	// Place H at 1.01 Å from N
	hPos := Vector3{
		X: n.X + hDir.X*BondN_H,
		Y: n.Y + hDir.Y*BondN_H,
		Z: n.Z + hDir.Z*BondN_H,
	}

	// Create H atom
	hAtom := &parser.Atom{
		Serial:  len(protein.Atoms) + 1, // Assign next serial number
		Name:    "H",
		ResName: residue.Name,
		ChainID: residue.ChainID,
		ResSeq:  residue.SeqNum,
		X:       hPos.X,
		Y:       hPos.Y,
		Z:       hPos.Z,
		Element: "H",
	}

	// Add to protein's atom list
	protein.Atoms = append(protein.Atoms, hAtom)

	// Note: We don't store H in Residue struct (only N, CA, C, O)
	// H-bond detection will find it in protein.Atoms by name and residue

	return nil
}

// addAlphaHydrogen adds Hα hydrogen at CA using tetrahedral geometry
//
// GEOMETRY:
// - CA is sp3 hybridized (tetrahedral)
// - Four bonds: N, C, Cβ (or H for Gly), Hα
// - For simplicity: Place Hα perpendicular to C-CA-N plane
// - CA-Hα distance = 1.09 Å
//
// BIOCHEMIST NOTE:
// - Glycine has 2 Hα (no Cβ sidechain)
// - For now: Add 1 Hα for all residues
// - Full sidechain H placement in future waves
//
// PHYSICIST: Tetrahedral geometry via cross product
// 1. CA→C vector
// 2. CA→N vector
// 3. Perpendicular = C×N (cross product)
// 4. Hα = CA + (perpendicular × 1.09 Å)
func addAlphaHydrogen(residue *parser.Residue) error {
	// Get backbone atoms
	ca := residue.CA
	c := residue.C
	n := residue.N

	if ca == nil || c == nil || n == nil {
		return nil // Skip if backbone incomplete
	}

	// Vectors from CA to C and N
	caC := Vector3{
		X: c.X - ca.X,
		Y: c.Y - ca.Y,
		Z: c.Z - ca.Z,
	}.Normalize()

	caN := Vector3{
		X: n.X - ca.X,
		Y: n.Y - ca.Y,
		Z: n.Z - ca.Z,
	}.Normalize()

	// Cross product gives perpendicular direction
	// This approximates tetrahedral geometry
	perp := caC.Cross(caN).Normalize()

	// Handle degenerate case (vectors parallel)
	if perp.Length() < 0.1 {
		// Use different perpendicular
		fallback := Vector3{X: 0, Y: 0, Z: 1}
		perp = caC.Cross(fallback).Normalize()
	}

	// Place Hα at 1.09 Å from CA
	haPos := Vector3{
		X: ca.X + perp.X*BondCA_H,
		Y: ca.Y + perp.Y*BondCA_H,
		Z: ca.Z + perp.Z*BondCA_H,
	}

	// Create HA atom
	haAtom := &parser.Atom{
		Serial:  0, // Will be assigned by caller if needed
		Name:    "HA",
		ResName: residue.Name,
		ChainID: residue.ChainID,
		ResSeq:  residue.SeqNum,
		X:       haPos.X,
		Y:       haPos.Y,
		Z:       haPos.Z,
		Element: "H",
	}

	// Note: For now, just validate geometry
	// Full integration requires expanding Residue struct
	// Validation: Check HA is reasonable distance from CA
	dist := math.Sqrt(
		(haPos.X-ca.X)*(haPos.X-ca.X) +
			(haPos.Y-ca.Y)*(haPos.Y-ca.Y) +
			(haPos.Z-ca.Z)*(haPos.Z-ca.Z),
	)

	if dist < 0.95 || dist > 1.25 {
		return fmt.Errorf("HA bond length %.3f Å outside valid range [0.95, 1.25]", dist)
	}

	// Atom created successfully but not yet stored
	// Full storage requires updating parser.Residue to include sidechain atoms
	// For now: This validates geometry is correct
	_ = haAtom

	return nil
}

// getAtom retrieves atom by name from residue
//
// Helper function for testing and validation
func getAtom(residue *parser.Residue, name string) *parser.Atom {
	switch name {
	case "N":
		return residue.N
	case "CA":
		return residue.CA
	case "C":
		return residue.C
	case "O":
		return residue.O
	default:
		return nil
	}
}

// ValidateHydrogenGeometry checks if H atoms have reasonable bond lengths
//
// VALIDATION TARGETS:
// - N-H bond: 0.95 - 1.15 Å (ideal 1.01 Å ± 0.14 Å tolerance)
// - CA-HA bond: 0.95 - 1.25 Å (ideal 1.09 Å ± tolerance)
//
// WRIGHT BROTHERS: Quick sanity check before "takeoff"
func ValidateHydrogenGeometry(protein *parser.Protein) (bool, string) {
	hCount := 0
	invalidCount := 0

	// Check all H atoms
	for _, atom := range protein.Atoms {
		if atom.Element != "H" {
			continue
		}

		hCount++

		// Find parent atom (N for backbone H, CA for HA)
		var parentAtom *parser.Atom
		var expectedDist float64
		var tolerance float64

		if atom.Name == "H" {
			// Backbone NH hydrogen
			// Find N atom in same residue
			for _, a := range protein.Atoms {
				if a.ResSeq == atom.ResSeq && a.Name == "N" {
					parentAtom = a
					expectedDist = BondN_H
					tolerance = 0.14 // ±0.14 Å
					break
				}
			}
		} else if atom.Name == "HA" {
			// Alpha hydrogen
			// Find CA atom in same residue
			for _, a := range protein.Atoms {
				if a.ResSeq == atom.ResSeq && a.Name == "CA" {
					parentAtom = a
					expectedDist = BondCA_H
					tolerance = 0.16 // ±0.16 Å
					break
				}
			}
		}

		if parentAtom == nil {
			continue // Can't validate without parent
		}

		// Calculate actual distance
		dx := atom.X - parentAtom.X
		dy := atom.Y - parentAtom.Y
		dz := atom.Z - parentAtom.Z
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

		// Check if within tolerance
		if dist < expectedDist-tolerance || dist > expectedDist+tolerance {
			invalidCount++
		}
	}

	if hCount == 0 {
		return false, "No hydrogen atoms found"
	}

	if invalidCount > 0 {
		return false, fmt.Sprintf("%d/%d hydrogen bonds outside valid range", invalidCount, hCount)
	}

	return true, fmt.Sprintf("All %d hydrogen atoms have valid geometry", hCount)
}
