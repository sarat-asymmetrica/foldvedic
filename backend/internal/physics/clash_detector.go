package physics

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// ClashReport contains clash detection results
type ClashReport struct {
	HasClashes      bool
	ClashCount      int
	WorstClashDist  float64
	Energy          float64
	IsValid         bool
	ValidationError string
}

// DetectClashes checks for severe atomic overlaps
func DetectClashes(protein *parser.Protein) ClashReport {
	report := ClashReport{
		HasClashes:     false,
		ClashCount:     0,
		WorstClashDist: 999.9,
		IsValid:        true,
	}

	// VdW radii (Å)
	vdwRadii := map[string]float64{
		"H": 1.20,
		"C": 1.70,
		"N": 1.55,
		"O": 1.52,
		"S": 1.80,
	}

	// Collect all atoms from Protein.Atoms field (pre-populated)
	atoms := protein.Atoms
	if len(atoms) == 0 {
		// No atoms to check
		return report
	}

	// Check all atom pairs
	for i := 0; i < len(atoms); i++ {
		for j := i + 1; j < len(atoms); j++ {
			a1 := atoms[i]
			a2 := atoms[j]

			// Skip bonded atoms (same residue or adjacent residues)
			if math.Abs(float64(a1.ResSeq-a2.ResSeq)) <= 1 && a1.ChainID == a2.ChainID {
				continue
			}

			// Calculate distance
			dx := a1.X - a2.X
			dy := a1.Y - a2.Y
			dz := a1.Z - a2.Z
			dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

			// Get VdW radii
			r1 := vdwRadii[a1.Element]
			r2 := vdwRadii[a2.Element]
			if r1 == 0 {
				r1 = 1.70
			} // Default to carbon
			if r2 == 0 {
				r2 = 1.70
			}

			// Clash threshold: 0.6 × (r1 + r2)
			// Normal contact: ~1.0 × (r1 + r2)
			threshold := 0.6 * (r1 + r2)

			if dist < threshold {
				report.HasClashes = true
				report.ClashCount++
				if dist < report.WorstClashDist {
					report.WorstClashDist = dist
				}
			}
		}
	}

	return report
}

// ValidateCoordinates checks for NaN, infinity, unrealistic distances
func ValidateCoordinates(protein *parser.Protein) ClashReport {
	report := ClashReport{IsValid: true}

	// Check all atoms
	for _, atom := range protein.Atoms {
		// Check for NaN
		if math.IsNaN(atom.X) || math.IsNaN(atom.Y) || math.IsNaN(atom.Z) {
			report.IsValid = false
			report.ValidationError = fmt.Sprintf("NaN coordinate in atom %s (residue %d %s)",
				atom.Name, atom.ResSeq, atom.ResName)
			return report
		}

		// Check for infinity
		if math.IsInf(atom.X, 0) || math.IsInf(atom.Y, 0) || math.IsInf(atom.Z, 0) {
			report.IsValid = false
			report.ValidationError = fmt.Sprintf("Inf coordinate in atom %s (residue %d %s)",
				atom.Name, atom.ResSeq, atom.ResName)
			return report
		}

		// Check for unrealistic distances from origin (>1000 Å)
		dist := math.Sqrt(atom.X*atom.X + atom.Y*atom.Y + atom.Z*atom.Z)
		if dist > 1000.0 {
			report.IsValid = false
			report.ValidationError = fmt.Sprintf("Atom too far from origin: %.1f Å (residue %d %s)",
				dist, atom.ResSeq, atom.ResName)
			return report
		}
	}

	// Check backbone connectivity
	for i := 1; i < len(protein.Residues); i++ {
		prevRes := protein.Residues[i-1]
		currRes := protein.Residues[i]

		// Check if same chain
		if prevRes.ChainID != currRes.ChainID {
			continue
		}

		// Get C from previous residue and N from current residue
		prevC := prevRes.C
		currN := currRes.N

		if prevC != nil && currN != nil {
			dx := currN.X - prevC.X
			dy := currN.Y - prevC.Y
			dz := currN.Z - prevC.Z
			dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

			// Peptide bond length should be ~1.33 Å
			if dist < 1.0 || dist > 2.0 {
				report.IsValid = false
				report.ValidationError = fmt.Sprintf(
					"Broken peptide bond between residues %d-%d: %.2f Å (should be ~1.33)",
					prevRes.SeqNum, currRes.SeqNum, dist)
				return report
			}
		}
	}

	return report
}

// ScoreStructureQuality combines validation + clash detection
func ScoreStructureQuality(protein *parser.Protein) (float64, ClashReport) {
	// Validate coordinates first
	report := ValidateCoordinates(protein)
	if !report.IsValid {
		return -999999.9, report
	}

	// Check for clashes
	clashReport := DetectClashes(protein)
	report.HasClashes = clashReport.HasClashes
	report.ClashCount = clashReport.ClashCount
	report.WorstClashDist = clashReport.WorstClashDist

	// Calculate quality score (0-1)
	// 0 clashes = 1.0
	// >10 clashes = 0.0
	quality := 1.0
	if report.ClashCount > 0 {
		quality = math.Max(0.0, 1.0-float64(report.ClashCount)/10.0)
	}

	return quality, report
}

// Helper function removed - not needed with current Residue structure
// Residues directly store N, CA, C, O atom pointers
