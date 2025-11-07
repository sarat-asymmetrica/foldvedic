package folding

import (
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/physics"
)

// PUBLIC wrappers for Wave 4 integration

// NewProteinFromSequence creates a new protein structure from amino acid sequence
// Returns an extended chain conformation
func NewProteinFromSequence(sequence string) *parser.Protein {
	protein, err := buildExtendedChain(sequence)
	if err != nil {
		// Return empty protein on error
		return &parser.Protein{
			Name:     "failed",
			Residues: make([]*parser.Residue, 0),
			Atoms:    make([]*parser.Atom, 0),
		}
	}
	return protein
}

// CloneProtein creates a deep copy of a protein structure
// Exported version of cloneProtein for Wave 4
func CloneProtein(protein *parser.Protein) *parser.Protein {
	return cloneProtein(protein)
}

// CalculateEnergy calculates total energy of a protein structure
// Wrapper around physics.CalculateTotalEnergy for convenience
func CalculateEnergy(protein *parser.Protein) float64 {
	energy := physics.CalculateTotalEnergy(protein, 10.0, 12.0)
	return energy.Total
}

// GetSequence extracts amino acid sequence from protein structure
func GetSequence(protein *parser.Protein) string {
	if protein == nil || len(protein.Residues) == 0 {
		return ""
	}

	sequence := make([]byte, len(protein.Residues))
	for i, res := range protein.Residues {
		// Convert three-letter code to one-letter
		// Simplified: just take first letter for now
		if len(res.Name) > 0 {
			sequence[i] = res.Name[0]
		} else {
			sequence[i] = 'X' // Unknown
		}
	}
	return string(sequence)
}
