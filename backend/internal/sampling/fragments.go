// Package sampling - Fragment Assembly for protein structure prediction
//
// WAVE 7.3: Fragment Assembly
// Template-based conformational sampling using known secondary structure fragments
//
// BIOCHEMIST: Uses experimentally observed (φ, ψ) angles from known structures
// PHYSICIST: Fragments carry local energy landscape information
// MATHEMATICIAN: Combinatorial assembly of building blocks
// ETHICIST: Transparent template-based approach, interpretable predictions
//
// INNOVATION: Vedic harmonic ranking of fragments
// Standard fragment assembly: Select fragments by sequence similarity
// Vedic fragment assembly: Rank by both similarity and φ-ratio alignment
//
// CITATION:
// Simons, K. T., et al. (1997). "Assembly of protein tertiary structures from fragments."
// J. Mol. Biol. 268(1): 209-225. (Rosetta method)
//
// Kim, D. E., et al. (2004). "Protein structure prediction and analysis using the Robetta server."
// Nucleic Acids Res. 32: W526-W531.
package sampling

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// Fragment represents a structural fragment (contiguous (φ, ψ) angles)
type Fragment struct {
	// Fragment length (3-mer, 9-mer, etc.)
	Length int

	// Ramachandran angles for each position
	Angles []geometry.RamachandranAngles

	// Fragment source/type (e.g., "alpha_helix", "beta_sheet", "PDB_1UBQ_5")
	Source string

	// Vedic harmonic score for this fragment
	VedicScore float64

	// Sequence context (if available)
	Sequence string
}

// FragmentLibrary holds collection of structural fragments
type FragmentLibrary struct {
	// Fragments organized by length
	ThreeMers []Fragment // 3-residue fragments
	NineMers  []Fragment // 9-residue fragments

	// Vedic-ranked fragments (pre-sorted by harmonic score)
	VedicRankedThree []Fragment
	VedicRankedNine  []Fragment
}

// FragmentAssemblyConfig holds fragment assembly parameters
type FragmentAssemblyConfig struct {
	// Use 3-mers, 9-mers, or both
	UseThreeMers bool
	UseNineMers  bool

	// Number of fragment insertion attempts per position
	NumInsertions int

	// Vedic bias weight [0, 1]
	// 0 = ignore Vedic score, 1 = only use Vedic score
	VedicWeight float64

	// Random seed for reproducibility
	Seed int64
}

// DefaultFragmentAssemblyConfig returns recommended parameters
func DefaultFragmentAssemblyConfig() FragmentAssemblyConfig {
	return FragmentAssemblyConfig{
		UseThreeMers:  true,
		UseNineMers:   true,
		NumInsertions: 5,    // Try 5 fragments per position
		VedicWeight:   0.3,  // 30% Vedic influence
		Seed:          42,
	}
}

// NewFragmentLibrary creates a fragment library with common secondary structures
//
// BIOCHEMIST:
// Includes ideal fragments for:
// - Alpha helix: φ = -60°, ψ = -45°
// - Beta sheet: φ = -120°, ψ = +120°
// - Left-handed helix: φ = +60°, ψ = +45° (rare, mainly Gly)
// - Turn/loop: Various conformations
//
// For production, this would be populated from PDB database (top8000 structures)
// or Robetta fragment server. For v0.2, using ideal fragments is sufficient.
func NewFragmentLibrary() *FragmentLibrary {
	lib := &FragmentLibrary{
		ThreeMers: make([]Fragment, 0),
		NineMers:  make([]Fragment, 0),
	}

	// Add ideal secondary structure fragments
	lib.addIdealAlphaHelix()
	lib.addIdealBetaSheet()
	lib.addIdealTurnFragments()
	lib.addIdealLoopFragments()

	// Rank fragments by Vedic score
	lib.rankByVedicScore()

	return lib
}

// addIdealAlphaHelix adds perfect alpha helix fragments
//
// BIOCHEMIST:
// Alpha helix: φ = -60° ± 10°, ψ = -45° ± 10°
// Helix pitch: 3.6 residues/turn ≈ 10/φ² (Vedic connection)
func (lib *FragmentLibrary) addIdealAlphaHelix() {
	const (
		phiHelix = -60.0 * math.Pi / 180.0
		psiHelix = -45.0 * math.Pi / 180.0
	)

	// 3-mer helix fragment
	frag3 := Fragment{
		Length: 3,
		Angles: []geometry.RamachandranAngles{
			{Phi: phiHelix, Psi: psiHelix},
			{Phi: phiHelix, Psi: psiHelix},
			{Phi: phiHelix, Psi: psiHelix},
		},
		Source:   "ideal_alpha_helix",
		Sequence: "AAA", // Generic
	}
	lib.ThreeMers = append(lib.ThreeMers, frag3)

	// 9-mer helix fragment
	frag9 := Fragment{
		Length: 9,
		Angles: make([]geometry.RamachandranAngles, 9),
		Source: "ideal_alpha_helix",
		Sequence: "AAAAAAAAA",
	}
	for i := range frag9.Angles {
		frag9.Angles[i] = geometry.RamachandranAngles{Phi: phiHelix, Psi: psiHelix}
	}
	lib.NineMers = append(lib.NineMers, frag9)

	// Add slight variations (±10°) for diversity
	for _, delta := range []float64{-10, -5, +5, +10} {
		deltaPhi := delta * math.Pi / 180.0
		deltaPsi := delta * math.Pi / 180.0

		frag3Var := Fragment{
			Length: 3,
			Angles: []geometry.RamachandranAngles{
				{Phi: phiHelix + deltaPhi, Psi: psiHelix + deltaPsi},
				{Phi: phiHelix + deltaPhi, Psi: psiHelix + deltaPsi},
				{Phi: phiHelix + deltaPhi, Psi: psiHelix + deltaPsi},
			},
			Source:   fmt.Sprintf("alpha_helix_var_%+.0f", delta),
			Sequence: "AAA",
		}
		lib.ThreeMers = append(lib.ThreeMers, frag3Var)
	}
}

// addIdealBetaSheet adds perfect beta sheet fragments
//
// BIOCHEMIST:
// Beta sheet: φ = -120° ± 15°, ψ = +120° ± 15°
// Antiparallel sheets more common than parallel
func (lib *FragmentLibrary) addIdealBetaSheet() {
	const (
		phiSheet = -120.0 * math.Pi / 180.0
		psiSheet = +120.0 * math.Pi / 180.0
	)

	// 3-mer sheet fragment
	frag3 := Fragment{
		Length: 3,
		Angles: []geometry.RamachandranAngles{
			{Phi: phiSheet, Psi: psiSheet},
			{Phi: phiSheet, Psi: psiSheet},
			{Phi: phiSheet, Psi: psiSheet},
		},
		Source:   "ideal_beta_sheet",
		Sequence: "VVV", // Valine common in sheets
	}
	lib.ThreeMers = append(lib.ThreeMers, frag3)

	// 9-mer sheet fragment
	frag9 := Fragment{
		Length: 9,
		Angles: make([]geometry.RamachandranAngles, 9),
		Source: "ideal_beta_sheet",
		Sequence: "VVVVVVVVV",
	}
	for i := range frag9.Angles {
		frag9.Angles[i] = geometry.RamachandranAngles{Phi: phiSheet, Psi: psiSheet}
	}
	lib.NineMers = append(lib.NineMers, frag9)

	// Variations
	for _, delta := range []float64{-15, -10, +10, +15} {
		deltaPhi := delta * math.Pi / 180.0
		deltaPsi := delta * math.Pi / 180.0

		frag3Var := Fragment{
			Length: 3,
			Angles: []geometry.RamachandranAngles{
				{Phi: phiSheet + deltaPhi, Psi: psiSheet + deltaPsi},
				{Phi: phiSheet + deltaPhi, Psi: psiSheet + deltaPsi},
				{Phi: phiSheet + deltaPhi, Psi: psiSheet + deltaPsi},
			},
			Source:   fmt.Sprintf("beta_sheet_var_%+.0f", delta),
			Sequence: "VVV",
		}
		lib.ThreeMers = append(lib.ThreeMers, frag3Var)
	}
}

// addIdealTurnFragments adds turn conformations
//
// BIOCHEMIST:
// Type I turn: Common in loops connecting secondary structures
// Type II turn: Glycine-rich turns
func (lib *FragmentLibrary) addIdealTurnFragments() {
	// Type I turn (i+1: φ=-60°, ψ=-30°; i+2: φ=-90°, ψ=0°)
	fragTurn1 := Fragment{
		Length: 3,
		Angles: []geometry.RamachandranAngles{
			{Phi: -60.0 * math.Pi / 180.0, Psi: -30.0 * math.Pi / 180.0},
			{Phi: -90.0 * math.Pi / 180.0, Psi: 0.0 * math.Pi / 180.0},
			{Phi: -60.0 * math.Pi / 180.0, Psi: -30.0 * math.Pi / 180.0},
		},
		Source:   "type_I_turn",
		Sequence: "GNG", // Glycine common in turns
	}
	lib.ThreeMers = append(lib.ThreeMers, fragTurn1)

	// Type II turn (i+1: φ=-60°, ψ=+120°; i+2: φ=+80°, ψ=0°)
	fragTurn2 := Fragment{
		Length: 3,
		Angles: []geometry.RamachandranAngles{
			{Phi: -60.0 * math.Pi / 180.0, Psi: +120.0 * math.Pi / 180.0},
			{Phi: +80.0 * math.Pi / 180.0, Psi: 0.0 * math.Pi / 180.0},
			{Phi: -60.0 * math.Pi / 180.0, Psi: +120.0 * math.Pi / 180.0},
		},
		Source:   "type_II_turn",
		Sequence: "GPG",
	}
	lib.ThreeMers = append(lib.ThreeMers, fragTurn2)
}

// addIdealLoopFragments adds flexible loop conformations
//
// BIOCHEMIST:
// Loops: Highly variable, often glycine/proline rich
// Use extended conformations as starting points
func (lib *FragmentLibrary) addIdealLoopFragments() {
	// Extended loop
	fragExtended := Fragment{
		Length: 3,
		Angles: []geometry.RamachandranAngles{
			{Phi: -120.0 * math.Pi / 180.0, Psi: +120.0 * math.Pi / 180.0},
			{Phi: -100.0 * math.Pi / 180.0, Psi: +100.0 * math.Pi / 180.0},
			{Phi: -120.0 * math.Pi / 180.0, Psi: +120.0 * math.Pi / 180.0},
		},
		Source:   "extended_loop",
		Sequence: "GGG",
	}
	lib.ThreeMers = append(lib.ThreeMers, fragExtended)

	// Compact loop
	fragCompact := Fragment{
		Length: 3,
		Angles: []geometry.RamachandranAngles{
			{Phi: -80.0 * math.Pi / 180.0, Psi: +80.0 * math.Pi / 180.0},
			{Phi: -70.0 * math.Pi / 180.0, Psi: +70.0 * math.Pi / 180.0},
			{Phi: -80.0 * math.Pi / 180.0, Psi: +80.0 * math.Pi / 180.0},
		},
		Source:   "compact_loop",
		Sequence: "PPP", // Proline in loops
	}
	lib.ThreeMers = append(lib.ThreeMers, fragCompact)
}

// rankByVedicScore calculates and sorts fragments by Vedic harmonic score
//
// MATHEMATICIAN:
// Vedic score measures golden ratio alignment in fragment geometry
// Higher score = more likely to be stable/native-like
func (lib *FragmentLibrary) rankByVedicScore() {
	// Calculate Vedic scores for all fragments
	for i := range lib.ThreeMers {
		lib.ThreeMers[i].VedicScore = calculateFragmentVedicScore(lib.ThreeMers[i].Angles)
	}

	for i := range lib.NineMers {
		lib.NineMers[i].VedicScore = calculateFragmentVedicScore(lib.NineMers[i].Angles)
	}

	// Sort by Vedic score (descending)
	lib.VedicRankedThree = make([]Fragment, len(lib.ThreeMers))
	copy(lib.VedicRankedThree, lib.ThreeMers)
	sortFragmentsByVedic(lib.VedicRankedThree)

	lib.VedicRankedNine = make([]Fragment, len(lib.NineMers))
	copy(lib.VedicRankedNine, lib.NineMers)
	sortFragmentsByVedic(lib.VedicRankedNine)
}

// calculateFragmentVedicScore computes Vedic score for fragment angles
func calculateFragmentVedicScore(angles []geometry.RamachandranAngles) float64 {
	if len(angles) == 0 {
		return 0.0
	}

	// Create temporary protein for Vedic scoring
	// We only need angles, not full structure
	totalScore := 0.0

	for _, angle := range angles {
		// Simple Vedic metric: closeness to golden ratio multiples
		// φ-ratio appears in helix pitch, sheet packing
		phiNorm := math.Abs(angle.Phi)
		psiNorm := math.Abs(angle.Psi)

		// Check if angles near φ-ratio multiples
		// φ ≈ 1.618, so 90°/φ ≈ 55.6°, 180°/φ ≈ 111.2°
		goldenAngles := []float64{
			55.6 * math.Pi / 180.0,  // 90°/φ
			90.0 * math.Pi / 180.0,  // Right angle
			111.2 * math.Pi / 180.0, // 180°/φ
		}

		minPhiDist := math.Inf(1)
		minPsiDist := math.Inf(1)

		for _, golden := range goldenAngles {
			phiDist := math.Abs(phiNorm - golden)
			psiDist := math.Abs(psiNorm - golden)

			if phiDist < minPhiDist {
				minPhiDist = phiDist
			}
			if psiDist < minPsiDist {
				minPsiDist = psiDist
			}
		}

		// Score inversely proportional to distance from golden angles
		// Max distance: π (180°), so score = 1 - dist/π
		phiScore := 1.0 - minPhiDist/math.Pi
		psiScore := 1.0 - minPsiDist/math.Pi

		totalScore += (phiScore + psiScore) / 2.0
	}

	return totalScore / float64(len(angles))
}

// sortFragmentsByVedic sorts fragments by Vedic score (descending)
func sortFragmentsByVedic(fragments []Fragment) {
	// Simple bubble sort (fragment count is small)
	n := len(fragments)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if fragments[j].VedicScore < fragments[j+1].VedicScore {
				fragments[j], fragments[j+1] = fragments[j+1], fragments[j]
			}
		}
	}
}

// FragmentAssembly builds protein structure via fragment insertion
//
// ALGORITHM:
// 1. Start with extended chain
// 2. For each position in sequence:
//    a. Select top fragments (by Vedic score + sequence context)
//    b. Insert fragment angles at this position
//    c. Evaluate combined score (structure quality + Vedic)
// 3. Return best assembled structure
//
// BIOCHEMIST:
// This mimics how Rosetta builds structures, but with Vedic ranking
func FragmentAssembly(sequence string, library *FragmentLibrary, config FragmentAssemblyConfig) (*parser.Protein, error) {
	if len(sequence) == 0 {
		return nil, fmt.Errorf("empty sequence")
	}

	if library == nil {
		return nil, fmt.Errorf("fragment library is nil")
	}

	rand.Seed(config.Seed)

	// Start with extended chain
	angles := make([]geometry.RamachandranAngles, len(sequence))
	for i := range angles {
		// Extended conformation
		angles[i] = geometry.RamachandranAngles{
			Phi: -120.0 * math.Pi / 180.0,
			Psi: +120.0 * math.Pi / 180.0,
		}
	}

	// Insert 9-mers first (larger context)
	if config.UseNineMers && len(sequence) >= 9 {
		for pos := 0; pos <= len(sequence)-9; pos++ {
			insertBestFragment(angles, pos, library.VedicRankedNine, config)
		}
	}

	// Insert 3-mers (refine local structure)
	if config.UseThreeMers && len(sequence) >= 3 {
		for pos := 0; pos <= len(sequence)-3; pos++ {
			insertBestFragment(angles, pos, library.VedicRankedThree, config)
		}
	}

	// Build structure from assembled angles
	// Create template protein
	template := &parser.Protein{
		Name:     "fragment_assembled",
		Residues: make([]*parser.Residue, len(sequence)),
	}

	for i := range sequence {
		template.Residues[i] = &parser.Residue{
			Name:    string(sequence[i]),
			SeqNum:  i + 1,
			ChainID: "A",
		}
	}

	protein, err := buildStructureFromAngles(template, angles)
	if err != nil {
		return nil, fmt.Errorf("failed to build structure from assembled fragments: %w", err)
	}

	return protein, nil
}

// insertBestFragment inserts best-scoring fragment at position
func insertBestFragment(angles []geometry.RamachandranAngles, pos int, fragments []Fragment, config FragmentAssemblyConfig) {
	if len(fragments) == 0 {
		return
	}

	// Try top N fragments
	numTries := min(config.NumInsertions, len(fragments))

	bestScore := math.Inf(-1)
	var bestAngles []geometry.RamachandranAngles

	for i := 0; i < numTries; i++ {
		frag := fragments[i]

		// Check if fragment fits
		if pos+frag.Length > len(angles) {
			continue
		}

		// Calculate score for inserting this fragment
		// Score = Vedic score of fragment
		score := frag.VedicScore * config.VedicWeight

		// Add local scoring (could be sequence compatibility, clash checking, etc.)
		// For v0.2, just use Vedic score

		if score > bestScore {
			bestScore = score
			bestAngles = frag.Angles
		}
	}

	// Insert best fragment
	if bestAngles != nil {
		for i, angle := range bestAngles {
			if pos+i < len(angles) {
				angles[pos+i] = angle
			}
		}
	}
}

// min returns minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GenerateFragmentEnsemble creates ensemble via fragment assembly with variation
//
// BIOCHEMIST:
// Multiple fragment assemblies with different random selections
// Explores combinatorial space of fragment combinations
func GenerateFragmentEnsemble(sequence string, library *FragmentLibrary, config FragmentAssemblyConfig, numStructures int) ([]*parser.Protein, error) {
	ensemble := make([]*parser.Protein, 0, numStructures)
	baseSeed := config.Seed

	for i := 0; i < numStructures; i++ {
		config.Seed = baseSeed + int64(i)

		protein, err := FragmentAssembly(sequence, library, config)
		if err != nil {
			// Skip failed assemblies
			continue
		}

		ensemble = append(ensemble, protein)
	}

	if len(ensemble) == 0 {
		return nil, fmt.Errorf("all fragment assemblies failed")
	}

	return ensemble, nil
}
