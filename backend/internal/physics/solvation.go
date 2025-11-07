package physics

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// Three-letter to one-letter amino acid code mapping
var threeToOne = map[string]byte{
	"ALA": 'A', "CYS": 'C', "ASP": 'D', "GLU": 'E',
	"PHE": 'F', "GLY": 'G', "HIS": 'H', "ILE": 'I',
	"LYS": 'K', "LEU": 'L', "MET": 'M', "ASN": 'N',
	"PRO": 'P', "GLN": 'Q', "ARG": 'R', "SER": 'S',
	"THR": 'T', "VAL": 'V', "TRP": 'W', "TYR": 'Y',
}

// Hydrophobicity scale (Kyte-Doolittle)
// Positive = hydrophobic (prefers core)
// Negative = hydrophilic (prefers surface)
var hydrophobicityScale = map[byte]float64{
	'A': 1.8,  // Alanine
	'C': 2.5,  // Cysteine
	'D': -3.5, // Aspartic acid
	'E': -3.5, // Glutamic acid
	'F': 2.8,  // Phenylalanine
	'G': -0.4, // Glycine
	'H': -3.2, // Histidine
	'I': 4.5,  // Isoleucine
	'K': -3.9, // Lysine
	'L': 3.8,  // Leucine
	'M': 1.9,  // Methionine
	'N': -3.5, // Asparagine
	'P': -1.6, // Proline
	'Q': -3.5, // Glutamine
	'R': -4.5, // Arginine
	'S': -0.8, // Serine
	'T': -0.7, // Threonine
	'V': 4.2,  // Valine
	'W': -0.9, // Tryptophan
	'Y': -1.3, // Tyrosine
}

// CalculateSASA calculates Solvent-Accessible Surface Area for each residue
// Uses simplified Lee-Richards algorithm
func CalculateSASA(protein *parser.Protein) map[*parser.Residue]float64 {
	sasa := make(map[*parser.Residue]float64)

	// Probe radius (water molecule, ~1.4 Å)
	probeRadius := 1.4

	// Atomic radii (from Bondi, 1964)
	atomicRadii := map[string]float64{
		"N":  1.55,
		"CA": 1.70,
		"C":  1.70,
		"O":  1.52,
	}

	// For each residue, calculate SASA of its CA atom (simplified)
	// Real SASA would consider all atoms, but CA is good approximation
	for _, residue := range protein.Residues {
		if residue.CA == nil {
			continue
		}

		// Extended radius = atomic radius + probe radius
		caRadius := atomicRadii["CA"] + probeRadius

		// Generate probe points on sphere around CA
		numProbePoints := 100 // More points = more accurate, but slower
		exposedPoints := 0

		for i := 0; i < numProbePoints; i++ {
			// Fibonacci sphere point generation
			phi := math.Pi * (3.0 - math.Sqrt(5.0)) // Golden angle in radians
			y := 1.0 - (float64(i)/float64(numProbePoints-1))*2.0
			radius := math.Sqrt(1.0 - y*y)
			theta := phi * float64(i)

			x := math.Cos(theta) * radius
			z := math.Sin(theta) * radius

			// Scale to CA radius
			probeX := residue.CA.X + x*caRadius
			probeY := residue.CA.Y + y*caRadius
			probeZ := residue.CA.Z + z*caRadius

			// Check if this probe point is buried by other atoms
			isExposed := true
			for _, otherResidue := range protein.Residues {
				if otherResidue.SeqNum == residue.SeqNum {
					continue // Skip self
				}

				// Check distance to other CA atoms
				if otherResidue.CA != nil {
					dx := probeX - otherResidue.CA.X
					dy := probeY - otherResidue.CA.Y
					dz := probeZ - otherResidue.CA.Z
					dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

					otherRadius := atomicRadii["CA"] + probeRadius
					if dist < otherRadius {
						isExposed = false
						break
					}
				}
			}

			if isExposed {
				exposedPoints++
			}
		}

		// SASA is proportional to fraction of exposed points
		// Multiply by sphere surface area
		sphereArea := 4.0 * math.Pi * caRadius * caRadius
		residueSASA := sphereArea * float64(exposedPoints) / float64(numProbePoints)

		sasa[residue] = residueSASA
	}

	return sasa
}

// CalculateSolvationEnergy calculates implicit solvation energy
// Uses SASA-based model (similar to EEF1)
func CalculateSolvationEnergy(protein *parser.Protein) float64 {
	sasa := CalculateSASA(protein)

	totalEnergy := 0.0

	for residue, residueSASA := range sasa {
		// Get hydrophobicity parameter for this residue
		if residue == nil {
			continue
		}

		aa, ok := threeToOne[residue.Name]
		if !ok {
			continue
		}
		hydrophobicity, ok := hydrophobicityScale[aa]
		if !ok {
			hydrophobicity = 0.0 // Unknown residue
		}

		// Solvation free energy = σ × SASA
		// σ is the atomic solvation parameter (kcal/mol/Ų)
		// Derived from hydrophobicity: more hydrophobic = more penalty for exposure

		// Scale factor (empirical, from EEF1 force field)
		sigma := hydrophobicity * 0.012 // kcal/mol/Ų

		// Energy contribution
		energy := sigma * residueSASA

		totalEnergy += energy
	}

	return totalEnergy
}

// GetBurialStatistics calculates statistics about residue burial
type BurialStatistics struct {
	NumBuried    int     // SASA < 20 Ų
	NumPartial   int     // 20 < SASA < 100 Ų
	NumExposed   int     // SASA > 100 Ų
	AvgSASA      float64
	TotalSASA    float64
	HydrophobicBuried   int // Hydrophobic residues in core (good)
	HydrophilicBuried   int // Hydrophilic residues in core (bad)
	HydrophobicExposed  int // Hydrophobic residues on surface (bad)
	HydrophilicExposed  int // Hydrophilic residues on surface (good)
}

func GetBurialStatistics(protein *parser.Protein) BurialStatistics {
	sasa := CalculateSASA(protein)

	stats := BurialStatistics{}

	sumSASA := 0.0

	for residue, residueSASA := range sasa {
		sumSASA += residueSASA

		// Classify burial level
		if residueSASA < 20.0 {
			stats.NumBuried++
		} else if residueSASA < 100.0 {
			stats.NumPartial++
		} else {
			stats.NumExposed++
		}

		// Get hydrophobicity
		if residue == nil {
			continue
		}

		aa, ok := threeToOne[residue.Name]
		if !ok {
			continue
		}
		hydrophobicity, ok := hydrophobicityScale[aa]
		if !ok {
			continue
		}

		isHydrophobic := hydrophobicity > 0.5
		isBuried := residueSASA < 50.0

		if isHydrophobic && isBuried {
			stats.HydrophobicBuried++ // Good: hydrophobic in core
		} else if !isHydrophobic && isBuried {
			stats.HydrophilicBuried++ // Bad: hydrophilic in core
		} else if isHydrophobic && !isBuried {
			stats.HydrophobicExposed++ // Bad: hydrophobic on surface
		} else {
			stats.HydrophilicExposed++ // Good: hydrophilic on surface
		}
	}

	if len(sasa) > 0 {
		stats.AvgSASA = sumSASA / float64(len(sasa))
	}
	stats.TotalSASA = sumSASA

	return stats
}

// CalculateHydrophobicEffect calculates hydrophobic collapse energy
// Rewards buried hydrophobic residues, penalizes exposed ones
func CalculateHydrophobicEffect(protein *parser.Protein) float64 {
	sasa := CalculateSASA(protein)

	totalEnergy := 0.0

	for residue, residueSASA := range sasa {
		if residue == nil {
			continue
		}

		aa, ok := threeToOne[residue.Name]
		if !ok {
			continue
		}
		hydrophobicity, ok := hydrophobicityScale[aa]
		if !ok {
			continue
		}

		// For hydrophobic residues (positive scale), burial is favored
		// For hydrophilic residues (negative scale), exposure is favored
		// Energy = hydrophobicity × SASA
		// This creates driving force for hydrophobic collapse

		energy := hydrophobicity * residueSASA * 0.05 // Scale factor

		totalEnergy += energy
	}

	return totalEnergy
}

// CalculateEntropyPenalty calculates entropy loss upon folding
// Simplified: proportional to number of buried residues
func CalculateEntropyPenalty(protein *parser.Protein) float64 {
	sasa := CalculateSASA(protein)

	numBuried := 0
	for _, residueSASA := range sasa {
		if residueSASA < 50.0 {
			numBuried++
		}
	}

	// Entropy penalty: ~1 kcal/mol per buried residue (rough estimate)
	// This opposes folding and must be overcome by favorable interactions
	entropyPenalty := float64(numBuried) * 1.0

	return entropyPenalty
}

// CalculateTotalSolvationFreeEnergy combines all solvation terms
func CalculateTotalSolvationFreeEnergy(protein *parser.Protein) float64 {
	// ΔG_solvation = ΔG_transfer + ΔG_hydrophobic - TΔS_burial

	solvationEnergy := CalculateSolvationEnergy(protein)
	hydrophobicEnergy := CalculateHydrophobicEffect(protein)
	entropyPenalty := CalculateEntropyPenalty(protein)

	totalEnergy := solvationEnergy + hydrophobicEnergy + entropyPenalty

	return totalEnergy
}
