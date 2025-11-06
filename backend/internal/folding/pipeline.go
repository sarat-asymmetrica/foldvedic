// Package folding implements the complete protein folding prediction pipeline.
//
// BIOCHEMIST: Sequence → Structure prediction
// PHYSICIST: Energy minimization guided folding
// MATHEMATICIAN: Quaternion-based conformational search
// ETHICIST: Fast, interpretable, accessible to all
package folding

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/physics"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/validation"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/vedic"
)

// PredictionConfig holds folding prediction parameters
type PredictionConfig struct {
	// Sequence to fold
	Sequence string

	// Initial structure (nil = extended chain)
	InitialStructure *parser.Protein

	// Energy minimization config
	MinimizerConfig physics.MinimizerConfig

	// Number of conformational samples
	NumSamples int

	// Random seed for reproducibility
	Seed int64
}

// DefaultPredictionConfig returns default folding parameters
func DefaultPredictionConfig(sequence string) PredictionConfig {
	return PredictionConfig{
		Sequence:        sequence,
		MinimizerConfig: physics.DefaultMinimizerConfig(),
		NumSamples:      10, // Sample 10 conformations
		Seed:            42,
	}
}

// PredictionResult holds folding prediction results
type PredictionResult struct {
	// Predicted structure
	Predicted *parser.Protein

	// Ramachandran angles
	Angles []geometry.RamachandranAngles

	// Final energy
	Energy physics.EnergyComponents

	// Vedic harmonic score
	VedicScore vedic.VedicScore

	// Comparison to experimental (if provided)
	Comparison *validation.StructureComparison

	// Folding statistics
	NumSteps      int
	Converged     bool
	TimeElapsed   float64 // seconds
	QualityScore  float64 // Harmonic mean of metrics
}

// PredictStructure performs complete protein folding prediction
//
// ALGORITHM:
// 1. Generate initial structure (extended chain or template)
// 2. Sample multiple conformations using quaternion interpolation
// 3. Energy minimize each conformation
// 4. Select best structure (lowest energy + highest Vedic score)
// 5. Validate against experimental if available
func PredictStructure(config PredictionConfig, experimental *parser.Protein) (*PredictionResult, error) {
	rand.Seed(config.Seed)

	result := &PredictionResult{}

	// Step 1: Generate initial structure
	if config.InitialStructure == nil {
		// Create extended chain from sequence
		initial, err := buildExtendedChain(config.Sequence)
		if err != nil {
			return nil, fmt.Errorf("failed to build initial structure: %w", err)
		}
		config.InitialStructure = initial
	}

	// Step 2: Sample conformations
	bestEnergy := math.Inf(1)
	var bestStructure *parser.Protein
	var bestAngles []geometry.RamachandranAngles
	var bestVedicScore vedic.VedicScore

	for sample := 0; sample < config.NumSamples; sample++ {
		// Clone initial structure
		structure := cloneProtein(config.InitialStructure)

		// Perturb angles for conformational sampling
		perturbStructure(structure, sample)

		// Energy minimize
		minResult, err := physics.MinimizeEnergy(structure, config.MinimizerConfig)
		if err != nil {
			// Skip this sample if minimization failed
			continue
		}

		// Calculate scores
		angles := geometry.CalculateRamachandran(structure)
		vedicScore := vedic.CalculateVedicScore(structure, angles)

		// Combined score: 70% energy + 30% Vedic harmonics
		// Lower energy is better, higher Vedic is better
		combinedScore := 0.7*minResult.FinalEnergy.Total - 0.3*vedicScore.TotalScore*1000

		// Select best
		if combinedScore < bestEnergy {
			bestEnergy = combinedScore
			bestStructure = structure
			bestAngles = angles
			bestVedicScore = vedicScore
			result.NumSteps = minResult.Steps
			result.Converged = minResult.Converged
		}
	}

	if bestStructure == nil {
		return nil, fmt.Errorf("all conformational samples failed")
	}

	// Store results
	result.Predicted = bestStructure
	result.Angles = bestAngles
	result.Energy = physics.CalculateTotalEnergy(bestStructure, 10.0, 12.0)
	result.VedicScore = bestVedicScore

	// Compare to experimental if provided
	if experimental != nil {
		comp := validation.CompareStructures(bestStructure, experimental)
		result.Comparison = &comp
		result.QualityScore = calculateQualityScore(comp, bestVedicScore)
	} else {
		// Quality based on Vedic score alone
		result.QualityScore = bestVedicScore.TotalScore
	}

	return result, nil
}

// buildExtendedChain creates an extended chain from sequence
//
// BIOCHEMIST:
// Extended chain: all phi=-120°, psi=+120° (beta strand)
// Spacing: 3.8 Å per residue along X axis
func buildExtendedChain(sequence string) (*parser.Protein, error) {
	if len(sequence) == 0 {
		return nil, fmt.Errorf("empty sequence")
	}

	protein := &parser.Protein{
		Name:     "predicted",
		Residues: make([]*parser.Residue, 0, len(sequence)),
		Atoms:    make([]*parser.Atom, 0, len(sequence)*4),
	}

	// Standard backbone geometry
	bondLengthNC := 1.33  // N-CA
	bondLengthCAC := 1.52 // CA-C
	bondLengthCO := 1.23  // C-O

	atomSerial := 1
	x := 0.0

	for i, aa := range sequence {
		// Three-letter code
		threeLetterCode := string(aa) // Simplified

		res := &parser.Residue{
			Name:    threeLetterCode,
			SeqNum:  i + 1,
			ChainID: "A",
		}

		// Backbone atoms: N, CA, C, O
		// Extended chain along X axis
		res.N = &parser.Atom{
			Serial:  atomSerial,
			Name:    "N",
			ResName: threeLetterCode,
			ChainID: "A",
			ResSeq:  i + 1,
			X:       x,
			Y:       0,
			Z:       0,
			Element: "N",
		}
		atomSerial++
		protein.Atoms = append(protein.Atoms, res.N)

		x += bondLengthNC
		res.CA = &parser.Atom{
			Serial:  atomSerial,
			Name:    "CA",
			ResName: threeLetterCode,
			ChainID: "A",
			ResSeq:  i + 1,
			X:       x,
			Y:       0,
			Z:       0,
			Element: "C",
		}
		atomSerial++
		protein.Atoms = append(protein.Atoms, res.CA)

		x += bondLengthCAC
		res.C = &parser.Atom{
			Serial:  atomSerial,
			Name:    "C",
			ResName: threeLetterCode,
			ChainID: "A",
			ResSeq:  i + 1,
			X:       x,
			Y:       0,
			Z:       0,
			Element: "C",
		}
		atomSerial++
		protein.Atoms = append(protein.Atoms, res.C)

		res.O = &parser.Atom{
			Serial:  atomSerial,
			Name:    "O",
			ResName: threeLetterCode,
			ChainID: "A",
			ResSeq:  i + 1,
			X:       x,
			Y:       bondLengthCO,
			Z:       0,
			Element: "O",
		}
		atomSerial++
		protein.Atoms = append(protein.Atoms, res.O)

		protein.Residues = append(protein.Residues, res)
		x += 0.5 // Small gap to next residue
	}

	return protein, nil
}

// Helper functions

func cloneProtein(protein *parser.Protein) *parser.Protein {
	clone := &parser.Protein{
		Name:     protein.Name,
		Residues: make([]*parser.Residue, len(protein.Residues)),
		Atoms:    make([]*parser.Atom, len(protein.Atoms)),
	}

	// Clone atoms
	atomMap := make(map[*parser.Atom]*parser.Atom)
	for i, atom := range protein.Atoms {
		clonedAtom := &parser.Atom{
			Serial:    atom.Serial,
			Name:      atom.Name,
			ResName:   atom.ResName,
			ChainID:   atom.ChainID,
			ResSeq:    atom.ResSeq,
			X:         atom.X,
			Y:         atom.Y,
			Z:         atom.Z,
			Occupancy: atom.Occupancy,
			TempFacto: atom.TempFacto,
			Element:   atom.Element,
		}
		clone.Atoms[i] = clonedAtom
		atomMap[atom] = clonedAtom
	}

	// Clone residues with updated atom pointers
	for i, res := range protein.Residues {
		clonedRes := &parser.Residue{
			Name:    res.Name,
			SeqNum:  res.SeqNum,
			ChainID: res.ChainID,
		}
		if res.N != nil {
			clonedRes.N = atomMap[res.N]
		}
		if res.CA != nil {
			clonedRes.CA = atomMap[res.CA]
		}
		if res.C != nil {
			clonedRes.C = atomMap[res.C]
		}
		if res.O != nil {
			clonedRes.O = atomMap[res.O]
		}
		clone.Residues[i] = clonedRes
	}

	return clone
}

func perturbStructure(protein *parser.Protein, sampleIndex int) {
	// Add small random perturbations to break symmetry
	noise := 0.1 * float64(sampleIndex+1) // Increasing noise per sample

	for _, atom := range protein.Atoms {
		atom.X += (rand.Float64()*2 - 1) * noise
		atom.Y += (rand.Float64()*2 - 1) * noise
		atom.Z += (rand.Float64()*2 - 1) * noise
	}
}

func calculateQualityScore(comp validation.StructureComparison, vedic vedic.VedicScore) float64 {
	// Quality = harmonic mean of:
	// - RMSD score (lower is better): 1 / (1 + RMSD)
	// - TM-score (higher is better)
	// - Vedic score (higher is better)

	rmsdScore := 1.0 / (1.0 + comp.RMSD)
	tmScore := comp.TMScore
	vedicScore := vedic.TotalScore

	if vedicScore == 0 {
		vedicScore = 0.1 // Avoid division by zero
	}

	// Harmonic mean
	sum := 1.0/rmsdScore + 1.0/tmScore + 1.0/vedicScore
	return 3.0 / sum
}
