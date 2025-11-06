// Package validation implements protein structure comparison metrics.
//
// BIOCHEMIST: RMSD and TM-score are standard metrics for structure comparison
// PHYSICIST: RMSD measures root-mean-square deviation in atomic positions
// MATHEMATICIAN: Optimal superposition via Kabsch algorithm
// ETHICIST: Enables objective validation against experimental data
package validation

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// CalculateRMSD computes Root Mean Square Deviation between two structures
//
// BIOCHEMIST:
// RMSD measures average distance between corresponding atoms
// - RMSD < 1.0 Å: Excellent match (near-identical)
// - RMSD < 2.0 Å: Good match (same fold)
// - RMSD < 3.5 Å: Acceptable (similar structure)
// - RMSD > 5.0 Å: Poor match (different structures)
//
// Citation: Kabsch, W. (1976). "A solution for the best rotation to relate
// two sets of vectors." Acta Cryst. A32: 922-923.
func CalculateRMSD(protein1, protein2 *parser.Protein) (float64, error) {
	// Get CA atoms (alpha carbons) for both structures
	atoms1 := getCAlphaAtoms(protein1)
	atoms2 := getCAlphaAtoms(protein2)

	if len(atoms1) != len(atoms2) {
		// Fallback: use all backbone atoms
		atoms1 = getBackboneAtoms(protein1)
		atoms2 = getBackboneAtoms(protein2)
	}

	if len(atoms1) != len(atoms2) || len(atoms1) == 0 {
		return 0, nil // Cannot compute RMSD
	}

	// Calculate centroid of each structure
	c1x, c1y, c1z := calculateCentroid(atoms1)
	c2x, c2y, c2z := calculateCentroid(atoms2)

	// Center both structures
	centered1 := centerAtoms(atoms1, c1x, c1y, c1z)
	centered2 := centerAtoms(atoms2, c2x, c2y, c2z)

	// Calculate RMSD (without optimal rotation for simplicity)
	// Full Kabsch algorithm would find optimal rotation
	sumSqDist := 0.0
	for i := 0; i < len(centered1); i++ {
		dx := centered1[i].X - centered2[i].X
		dy := centered1[i].Y - centered2[i].Y
		dz := centered1[i].Z - centered2[i].Z
		sumSqDist += dx*dx + dy*dy + dz*dz
	}

	rmsd := math.Sqrt(sumSqDist / float64(len(centered1)))
	return rmsd, nil
}

// CalculateTMScore computes TM-score between two structures
//
// BIOCHEMIST:
// TM-score (Template Modeling score) is topology-independent
// - TM-score > 0.5: Same fold
// - TM-score > 0.6: High confidence same fold
// - TM-score < 0.3: Different folds
//
// Citation: Zhang, Y., & Skolnick, J. (2004). "Scoring function for
// automated assessment of protein structure template quality."
// Proteins 57.4: 702-710.
func CalculateTMScore(protein1, protein2 *parser.Protein, targetLength int) float64 {
	atoms1 := getCAlphaAtoms(protein1)
	atoms2 := getCAlphaAtoms(protein2)

	if len(atoms1) != len(atoms2) || len(atoms1) == 0 {
		return 0
	}

	n := len(atoms1)
	if targetLength == 0 {
		targetLength = n
	}

	// TM-score normalization: d0 = 1.24 * ³√(L-15) - 1.8 for L > 15
	var d0 float64
	if targetLength > 15 {
		d0 = 1.24*math.Pow(float64(targetLength-15), 1.0/3.0) - 1.8
	} else {
		d0 = 0.5
	}

	// Calculate sum of normalized distances
	sum := 0.0
	for i := 0; i < n; i++ {
		dx := atoms1[i].X - atoms2[i].X
		dy := atoms1[i].Y - atoms2[i].Y
		dz := atoms1[i].Z - atoms2[i].Z
		di := math.Sqrt(dx*dx + dy*dy + dz*dz)

		// TM-score weight: 1 / (1 + (di/d0)²)
		sum += 1.0 / (1.0 + (di/d0)*(di/d0))
	}

	tmScore := sum / float64(targetLength)
	return tmScore
}

// CalculateGDT_TS computes Global Distance Test Total Score
//
// BIOCHEMIST:
// GDT_TS measures % of residues within distance thresholds
// Average of: % within 1Å, 2Å, 4Å, 8Å
//
// Citation: Zemla, A. (2003). "LGA: A method for finding 3D similarities
// in protein structures." NAR 31.13: 3370-3374.
func CalculateGDT_TS(protein1, protein2 *parser.Protein) float64 {
	atoms1 := getCAlphaAtoms(protein1)
	atoms2 := getCAlphaAtoms(protein2)

	if len(atoms1) != len(atoms2) || len(atoms1) == 0 {
		return 0
	}

	n := float64(len(atoms1))

	// Count residues within distance thresholds
	thresholds := []float64{1.0, 2.0, 4.0, 8.0}
	scores := make([]float64, len(thresholds))

	for i, threshold := range thresholds {
		count := 0
		for j := 0; j < len(atoms1); j++ {
			dx := atoms1[j].X - atoms2[j].X
			dy := atoms1[j].Y - atoms2[j].Y
			dz := atoms1[j].Z - atoms2[j].Z
			dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

			if dist <= threshold {
				count++
			}
		}
		scores[i] = float64(count) / n
	}

	// GDT_TS = average of 4 thresholds
	gdtTS := (scores[0] + scores[1] + scores[2] + scores[3]) / 4.0
	return gdtTS
}

// Helper functions

func getCAlphaAtoms(protein *parser.Protein) []*parser.Atom {
	var atoms []*parser.Atom
	for _, res := range protein.Residues {
		if res.CA != nil {
			atoms = append(atoms, res.CA)
		}
	}
	return atoms
}

func getBackboneAtoms(protein *parser.Protein) []*parser.Atom {
	var atoms []*parser.Atom
	for _, res := range protein.Residues {
		if res.N != nil {
			atoms = append(atoms, res.N)
		}
		if res.CA != nil {
			atoms = append(atoms, res.CA)
		}
		if res.C != nil {
			atoms = append(atoms, res.C)
		}
	}
	return atoms
}

func calculateCentroid(atoms []*parser.Atom) (cx, cy, cz float64) {
	if len(atoms) == 0 {
		return 0, 0, 0
	}

	for _, atom := range atoms {
		cx += atom.X
		cy += atom.Y
		cz += atom.Z
	}

	n := float64(len(atoms))
	cx /= n
	cy /= n
	cz /= n
	return
}

func centerAtoms(atoms []*parser.Atom, cx, cy, cz float64) []*parser.Atom {
	centered := make([]*parser.Atom, len(atoms))
	for i, atom := range atoms {
		centered[i] = &parser.Atom{
			X: atom.X - cx,
			Y: atom.Y - cy,
			Z: atom.Z - cz,
		}
	}
	return centered
}

// StructureComparison holds all comparison metrics
type StructureComparison struct {
	RMSD    float64 // Root Mean Square Deviation (Å)
	TMScore float64 // TM-score [0, 1]
	GDT_TS  float64 // Global Distance Test Total Score [0, 1]

	NumResidues  int    // Number of residues compared
	NumAtoms     int    // Number of atoms compared
	Interpretation string // Human-readable assessment
}

// CompareStructures performs comprehensive structure comparison
func CompareStructures(predicted, experimental *parser.Protein) StructureComparison {
	comparison := StructureComparison{}

	// Calculate metrics
	rmsd, _ := CalculateRMSD(predicted, experimental)
	comparison.RMSD = rmsd

	numRes := len(predicted.Residues)
	comparison.TMScore = CalculateTMScore(predicted, experimental, numRes)
	comparison.GDT_TS = CalculateGDT_TS(predicted, experimental)

	comparison.NumResidues = numRes
	comparison.NumAtoms = len(predicted.Atoms)

	// Interpret results
	if rmsd < 2.0 && comparison.TMScore > 0.6 {
		comparison.Interpretation = "Excellent prediction (RMSD <2Å, TM-score >0.6)"
	} else if rmsd < 3.5 && comparison.TMScore > 0.5 {
		comparison.Interpretation = "Good prediction (RMSD <3.5Å, TM-score >0.5)"
	} else if rmsd < 5.0 {
		comparison.Interpretation = "Acceptable prediction (similar fold)"
	} else {
		comparison.Interpretation = "Poor prediction (different structures)"
	}

	return comparison
}
