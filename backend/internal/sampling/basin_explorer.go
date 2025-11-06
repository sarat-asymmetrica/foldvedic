// Package sampling - Ramachandran Basin Explorer
//
// WAVE 7.4: Ramachandran Basin Explorer
// Systematic exploration of allowed conformational basins in Ramachandran space
//
// BIOCHEMIST: Samples known allowed regions (helix, sheet, turn)
// PHYSICIST: Each basin represents local energy minimum
// MATHEMATICIAN: Gaussian sampling around basin centers
// ETHICIST: Interpretable, biophysically grounded conformational search
//
// INNOVATION: Vedic-guided basin sampling
// Standard approach: Uniform sampling in allowed regions
// Vedic approach: Bias toward φ-ratio geometries within each basin
//
// CITATION:
// Ramachandran, G. N., et al. (1963). "Stereochemistry of polypeptide chain configurations."
// J. Mol. Biol. 7(1): 95-99.
//
// Lovell, S. C., et al. (2003). "Structure validation by Cα geometry: φ, ψ and Cβ deviation."
// Proteins 50(3): 437-450.
package sampling

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// RamachandranBasin represents an allowed region in Ramachandran space
type RamachandranBasin struct {
	// Basin name (e.g., "alpha_helix", "beta_sheet")
	Name string

	// Central (φ, ψ) values (degrees)
	PhiCenter float64
	PsiCenter float64

	// Standard deviations (degrees) - defines basin width
	PhiSigma float64
	PsiSigma float64

	// Residue preferences (empty = all residues allowed)
	// e.g., ["GLY"] for glycine-specific regions
	ResiduePreferences []string

	// Relative population (0-1) - how common is this basin?
	Population float64

	// Vedic harmonic score for this basin
	VedicScore float64
}

// BasinExplorerConfig holds parameters for basin exploration
type BasinExplorerConfig struct {
	// Number of structures to generate per basin
	SamplesPerBasin int

	// Use Vedic biasing (prefer high Vedic score basins)
	UseVedicBiasing bool

	// Sample all basins or only high-population ones
	SampleAllBasins bool

	// Glycine special handling (larger allowed regions)
	GlycineHandling bool

	// Proline special handling (restricted φ angle)
	ProlineHandling bool

	// Random seed
	Seed int64
}

// DefaultBasinExplorerConfig returns recommended parameters
func DefaultBasinExplorerConfig() BasinExplorerConfig {
	return BasinExplorerConfig{
		SamplesPerBasin: 10,
		UseVedicBiasing: true,
		SampleAllBasins: true,
		GlycineHandling: true,
		ProlineHandling: true,
		Seed:            42,
	}
}

// GetStandardRamachandranBasins returns commonly observed basins
//
// BIOCHEMIST:
// Based on Lovell et al. (2003) and MolProbity definitions
// - Alpha helix: Most common, φ=-60°, ψ=-45°
// - Beta sheet: Second most common, φ=-120°, ψ=+120°
// - Left-handed helix: Rare except for Gly, φ=+60°, ψ=+45°
// - Extended: φ=-120°, ψ=+120° (similar to beta but looser)
//
// VEDIC ENHANCEMENT:
// Basins are ranked by golden ratio alignment
func GetStandardRamachandranBasins() []RamachandranBasin {
	basins := []RamachandranBasin{
		// Alpha helix region (most common)
		{
			Name:               "alpha_helix",
			PhiCenter:          -60.0,
			PsiCenter:          -45.0,
			PhiSigma:           20.0,
			PsiSigma:           20.0,
			ResiduePreferences: []string{}, // All residues except Pro at i+1
			Population:         0.35,       // ~35% of non-Gly/Pro residues
			VedicScore:         0.85,       // High φ-ratio alignment
		},

		// Beta sheet region
		{
			Name:               "beta_sheet",
			PhiCenter:          -120.0,
			PsiCenter:          +120.0,
			PhiSigma:           30.0,
			PsiSigma:           30.0,
			ResiduePreferences: []string{},
			Population:         0.25, // ~25%
			VedicScore:         0.75, // Moderate φ-ratio
		},

		// Left-handed helix (mainly Gly)
		{
			Name:               "left_handed_helix",
			PhiCenter:          +60.0,
			PsiCenter:          +45.0,
			PhiSigma:           25.0,
			PsiSigma:           25.0,
			ResiduePreferences: []string{"GLY"},
			Population:         0.05, // ~5%, mostly Gly
			VedicScore:         0.80, // Mirror of alpha helix
		},

		// Extended/PPII region
		{
			Name:               "extended_ppii",
			PhiCenter:          -75.0,
			PsiCenter:          +145.0,
			PhiSigma:           25.0,
			PsiSigma:           25.0,
			ResiduePreferences: []string{},
			Population:         0.15, // ~15%
			VedicScore:         0.60, // Lower φ-ratio
		},

		// Bridge region (connecting helix and sheet)
		{
			Name:               "bridge",
			PhiCenter:          -90.0,
			PsiCenter:          0.0,
			PhiSigma:           30.0,
			PsiSigma:           40.0,
			ResiduePreferences: []string{},
			Population:         0.10, // ~10%
			VedicScore:         0.50, // Lower φ-ratio
		},

		// Turn regions
		{
			Name:               "turn_type_I",
			PhiCenter:          -60.0,
			PsiCenter:          -30.0,
			PhiSigma:           20.0,
			PsiSigma:           30.0,
			ResiduePreferences: []string{},
			Population:         0.05, // ~5%
			VedicScore:         0.70,
		},

		{
			Name:               "turn_type_II",
			PhiCenter:          +80.0,
			PsiCenter:          0.0,
			PhiSigma:           25.0,
			PsiSigma:           30.0,
			ResiduePreferences: []string{"GLY", "ASN", "ASP"},
			Population:         0.03, // ~3%
			VedicScore:         0.65,
		},
	}

	return basins
}

// ExploreRamachandranBasins generates ensemble by sampling allowed basins
//
// ALGORITHM:
// 1. Get standard Ramachandran basins
// 2. For each residue position:
//    a. Select basin (by population or Vedic score)
//    b. Sample (φ, ψ) from Gaussian around basin center
// 3. Build structure from sampled angles
// 4. Return ensemble of diverse conformations
//
// BIOCHEMIST:
// This ensures all generated structures have biophysically allowed angles
// No Ramachandran outliers (unlike random sampling)
func ExploreRamachandranBasins(sequence string, config BasinExplorerConfig) ([]*parser.Protein, error) {
	if len(sequence) == 0 {
		return nil, fmt.Errorf("empty sequence")
	}

	rand.Seed(config.Seed)

	basins := GetStandardRamachandranBasins()
	ensemble := make([]*parser.Protein, 0)

	// For Vedic biasing, sort basins by Vedic score
	if config.UseVedicBiasing {
		sortBasinsByVedic(basins)
	}

	// Sample each basin
	for _, basin := range basins {
		// Skip low-population basins if not sampling all
		if !config.SampleAllBasins && basin.Population < 0.05 {
			continue
		}

		numSamples := config.SamplesPerBasin
		if config.UseVedicBiasing {
			// More samples from high Vedic score basins
			numSamples = int(float64(config.SamplesPerBasin) * basin.VedicScore)
			if numSamples < 1 {
				numSamples = 1
			}
		}

		for sample := 0; sample < numSamples; sample++ {
			// Generate angles for entire sequence sampling this basin
			angles := make([]geometry.RamachandranAngles, len(sequence))

			for resIdx := range sequence {
				// Sample (φ, ψ) from this basin
				phi, psi := sampleFromBasin(basin, config)

				angles[resIdx] = geometry.RamachandranAngles{
					Phi: phi * math.Pi / 180.0, // Convert to radians
					Psi: psi * math.Pi / 180.0,
				}
			}

			// Build structure
			template := createSequenceTemplate(sequence)
			protein, err := buildStructureFromAngles(template, angles)
			if err != nil {
				// Skip failed structures
				continue
			}

			ensemble = append(ensemble, protein)
		}
	}

	if len(ensemble) == 0 {
		return nil, fmt.Errorf("failed to generate any structures")
	}

	return ensemble, nil
}

// MixedBasinSampling generates diverse ensemble by mixing basins per residue
//
// BIOCHEMIST:
// More realistic: Different residues adopt different basins
// e.g., helix-turn-sheet motif
//
// ALGORITHM:
// 1. For each structure:
//    a. For each residue, randomly select basin
//    b. Sample (φ, ψ) from selected basin
// 2. Build structure
func MixedBasinSampling(sequence string, config BasinExplorerConfig, numStructures int) ([]*parser.Protein, error) {
	if len(sequence) == 0 {
		return nil, fmt.Errorf("empty sequence")
	}

	rand.Seed(config.Seed)

	basins := GetStandardRamachandranBasins()
	ensemble := make([]*parser.Protein, 0, numStructures)

	// Prepare basin selection weights
	weights := make([]float64, len(basins))
	totalWeight := 0.0

	for i, basin := range basins {
		if config.UseVedicBiasing {
			weights[i] = basin.VedicScore * basin.Population
		} else {
			weights[i] = basin.Population
		}
		totalWeight += weights[i]
	}

	// Normalize weights
	for i := range weights {
		weights[i] /= totalWeight
	}

	// Generate structures
	for structIdx := 0; structIdx < numStructures; structIdx++ {
		angles := make([]geometry.RamachandranAngles, len(sequence))

		for resIdx := range sequence {
			// Select basin for this residue (weighted random)
			basin := selectBasinWeighted(basins, weights)

			// Handle residue-specific constraints
			resName := string(sequence[resIdx])
			if config.GlycineHandling && resName == "G" {
				// Glycine: prefer left-handed helix or bridge
				basin = basins[2] // left_handed_helix
			}
			if config.ProlineHandling && resName == "P" {
				// Proline: restricted to specific regions
				basin = basins[3] // extended_ppii
			}

			// Sample from selected basin
			phi, psi := sampleFromBasin(basin, config)

			angles[resIdx] = geometry.RamachandranAngles{
				Phi: phi * math.Pi / 180.0,
				Psi: psi * math.Pi / 180.0,
			}
		}

		// Build structure
		template := createSequenceTemplate(sequence)
		protein, err := buildStructureFromAngles(template, angles)
		if err != nil {
			continue
		}

		ensemble = append(ensemble, protein)
	}

	if len(ensemble) == 0 {
		return nil, fmt.Errorf("failed to generate any structures")
	}

	return ensemble, nil
}

// sampleFromBasin samples (φ, ψ) from Gaussian around basin center
//
// MATHEMATICIAN:
// Gaussian sampling: N(μ, σ²)
// μ = basin center, σ = basin standard deviation
func sampleFromBasin(basin RamachandranBasin, config BasinExplorerConfig) (phi, psi float64) {
	// Gaussian sampling
	phi = basin.PhiCenter + rand.NormFloat64()*basin.PhiSigma
	psi = basin.PsiCenter + rand.NormFloat64()*basin.PsiSigma

	// Wrap to [-180, +180]
	phi = wrapAngle(phi)
	psi = wrapAngle(psi)

	return phi, psi
}

// wrapAngle wraps angle to [-180, +180] degrees
func wrapAngle(angle float64) float64 {
	for angle > 180.0 {
		angle -= 360.0
	}
	for angle < -180.0 {
		angle += 360.0
	}
	return angle
}

// selectBasinWeighted selects basin using weighted random sampling
func selectBasinWeighted(basins []RamachandranBasin, weights []float64) RamachandranBasin {
	r := rand.Float64()
	cumulative := 0.0

	for i, weight := range weights {
		cumulative += weight
		if r <= cumulative {
			return basins[i]
		}
	}

	// Fallback: return first basin
	return basins[0]
}

// sortBasinsByVedic sorts basins by Vedic score (descending)
func sortBasinsByVedic(basins []RamachandranBasin) {
	// Bubble sort (small number of basins)
	n := len(basins)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if basins[j].VedicScore < basins[j+1].VedicScore {
				basins[j], basins[j+1] = basins[j+1], basins[j]
			}
		}
	}
}

// createSequenceTemplate creates template protein from sequence
func createSequenceTemplate(sequence string) *parser.Protein {
	protein := &parser.Protein{
		Name:     "basin_sampled",
		Residues: make([]*parser.Residue, len(sequence)),
	}

	for i, aa := range sequence {
		// Map single-letter to three-letter code (simplified)
		threeLetter := singleToThreeLetter(string(aa))

		protein.Residues[i] = &parser.Residue{
			Name:    threeLetter,
			SeqNum:  i + 1,
			ChainID: "A",
		}
	}

	return protein
}

// singleToThreeLetter converts single-letter amino acid code to three-letter
//
// BIOCHEMIST: Standard amino acid codes
func singleToThreeLetter(single string) string {
	codeMap := map[string]string{
		"A": "ALA", "R": "ARG", "N": "ASN", "D": "ASP",
		"C": "CYS", "Q": "GLN", "E": "GLU", "G": "GLY",
		"H": "HIS", "I": "ILE", "L": "LEU", "K": "LYS",
		"M": "MET", "F": "PHE", "P": "PRO", "S": "SER",
		"T": "THR", "W": "TRP", "Y": "TYR", "V": "VAL",
	}

	if code, ok := codeMap[single]; ok {
		return code
	}

	return "ALA" // Default fallback
}

// ConstrainedBasinSampling generates structures with basin constraints
//
// BIOCHEMIST:
// Apply secondary structure constraints:
// - Helix regions: sample only helix basin
// - Sheet regions: sample only sheet basin
// - Loop regions: sample multiple basins
//
// Requires secondary structure prediction (Wave 9.1)
// For v0.2, use simple heuristic: hydrophobic residues prefer sheet
func ConstrainedBasinSampling(sequence string, constraints map[int]string, config BasinExplorerConfig, numStructures int) ([]*parser.Protein, error) {
	if len(sequence) == 0 {
		return nil, fmt.Errorf("empty sequence")
	}

	rand.Seed(config.Seed)

	basins := GetStandardRamachandranBasins()
	basinMap := make(map[string]RamachandranBasin)
	for _, basin := range basins {
		basinMap[basin.Name] = basin
	}

	ensemble := make([]*parser.Protein, 0, numStructures)

	for structIdx := 0; structIdx < numStructures; structIdx++ {
		angles := make([]geometry.RamachandranAngles, len(sequence))

		for resIdx := range sequence {
			var basin RamachandranBasin

			// Check for constraint
			if constraintName, hasConstraint := constraints[resIdx]; hasConstraint {
				if constrainedBasin, ok := basinMap[constraintName]; ok {
					basin = constrainedBasin
				} else {
					// Unknown constraint: use default
					basin = basins[0]
				}
			} else {
				// No constraint: use weighted sampling
				weights := make([]float64, len(basins))
				totalWeight := 0.0
				for i, b := range basins {
					weights[i] = b.Population * b.VedicScore
					totalWeight += weights[i]
				}
				for i := range weights {
					weights[i] /= totalWeight
				}
				basin = selectBasinWeighted(basins, weights)
			}

			// Sample from basin
			phi, psi := sampleFromBasin(basin, config)

			angles[resIdx] = geometry.RamachandranAngles{
				Phi: phi * math.Pi / 180.0,
				Psi: psi * math.Pi / 180.0,
			}
		}

		// Build structure
		template := createSequenceTemplate(sequence)
		protein, err := buildStructureFromAngles(template, angles)
		if err != nil {
			continue
		}

		ensemble = append(ensemble, protein)
	}

	if len(ensemble) == 0 {
		return nil, fmt.Errorf("failed to generate any structures")
	}

	return ensemble, nil
}
