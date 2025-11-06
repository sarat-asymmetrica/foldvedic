// Package pipeline - Unified Pipeline v2
//
// WAVE 10.1: Unified Pipeline v2 - Integration of All Phase 2 Enhancements
// Combines Waves 7-9 into cohesive protein folding pipeline
//
// PHASE 1 (v0.1):
// - Extended chain initialization
// - Random perturbations
// - 100-step steepest descent
// - Result: 63.16 Å RMSD (poor accuracy, excellent infrastructure)
//
// PHASE 2 (v0.2):
// - Secondary structure prediction → intelligent initialization
// - Advanced sampling → quaternion slerp, Monte Carlo, fragments, basin explorer
// - Advanced optimization → L-BFGS, simulated annealing, 1000+ steps
// - Structural priors → contact maps, Vedic biasing
// - Target: <5 Å RMSD (competitive with Rosetta)
//
// INTEGRATION STRATEGY:
// 1. Predict secondary structure from sequence
// 2. Predict contact map
// 3. Initialize structure using SS prediction + basin sampling
// 4. Generate ensemble via advanced sampling (Wave 7)
// 5. Optimize each structure with advanced methods (Wave 8)
// 6. Apply Vedic biasing and contact restraints throughout
// 7. Select best structure, validate
//
// BIOCHEMIST: Complete ab initio folding pipeline
// PHYSICIST: Combines global search with local refinement
// MATHEMATICIAN: Multi-scale optimization strategy
// ETHICIST: Transparent, interpretable, reproducible
package pipeline

import (
	"fmt"
	"time"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/optimization"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/prediction"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/sampling"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/validation"
)

// UnifiedPipelineV2Config holds all configuration parameters
type UnifiedPipelineV2Config struct {
	// Input
	Sequence string

	// Secondary structure prediction
	UseSSprediction bool
	SSMethod        prediction.PredictionMethod

	// Contact map prediction
	UseContactMap bool
	ContactConfig prediction.ContactMapConfig

	// Sampling strategy (multiple can be enabled)
	UseQuaternionSlerp bool
	UseMonteCarlo      bool
	UseFragmentAssembly bool
	UseBasinExplorer   bool
	NumSamplesPerMethod int

	// Optimization strategy
	OptimizationStrategy optimization.OptimizationStrategy
	OptimizationConfig   optimization.AdaptiveOptimizationConfig

	// Vedic biasing
	UseVedicBiasing bool
	VedicBias       prediction.VedicStructuralBias

	// Output
	Verbose bool
}

// DefaultUnifiedPipelineV2Config returns recommended Phase 2 parameters
func DefaultUnifiedPipelineV2Config(sequence string) UnifiedPipelineV2Config {
	return UnifiedPipelineV2Config{
		Sequence:            sequence,
		UseSSprediction:     true,
		SSMethod:            prediction.MethodChouFasman,
		UseContactMap:       true,
		ContactConfig:       prediction.DefaultContactMapConfig(),
		UseQuaternionSlerp:  true,
		UseMonteCarlo:       true,
		UseFragmentAssembly: true,
		UseBasinExplorer:    true,
		NumSamplesPerMethod: 5, // 5 samples × 4 methods = 20 total
		OptimizationStrategy: optimization.StrategyHybrid,
		OptimizationConfig:   optimization.DefaultAdaptiveOptimizationConfig(),
		UseVedicBiasing:      true,
		VedicBias:            prediction.DefaultVedicStructuralBias(),
		Verbose:              false,
	}
}

// UnifiedPipelineV2Result holds comprehensive pipeline results
type UnifiedPipelineV2Result struct {
	// Predictions
	SecondaryStructure []prediction.SecondaryStructurePrediction
	ContactMap         []prediction.ContactPrediction
	VedicReport        prediction.VedicHarmonicReport

	// Final structure
	FinalStructure *parser.Protein
	FinalAngles    []geometry.RamachandranAngles

	// Energetics
	FinalEnergy      float64
	FinalVedicScore  float64
	CombinedScore    float64

	// Optimization statistics
	OptimizationResult *optimization.OptimizationResult

	// Comparison to experimental (if provided)
	Validation *validation.StructureComparison

	// Pipeline statistics
	TotalSamplesGenerated int
	TotalTimeSeconds      float64
	SuccessRate           float64

	// Quality assessment
	QualityScore float64 // Harmonic mean of all metrics
}

// RunUnifiedPipelineV2 executes complete Phase 2 folding pipeline
//
// ALGORITHM:
// Phase A: Prediction & Planning (Wave 9)
//   1. Predict secondary structure
//   2. Predict contact map
//   3. Plan sampling strategy
//
// Phase B: Conformational Sampling (Wave 7)
//   4. Generate ensemble using multiple methods:
//      - Quaternion slerp (S³ exploration)
//      - Monte Carlo (thermodynamic sampling)
//      - Fragment assembly (template-based)
//      - Basin explorer (Ramachandran sampling)
//
// Phase C: Energy Optimization (Wave 8)
//   5. Optimize each sampled structure:
//      - Simulated annealing (global search)
//      - L-BFGS (local refinement)
//      - Hybrid (SA → L-BFGS)
//
// Phase D: Selection & Validation
//   6. Score structures (energy + Vedic + contact satisfaction)
//   7. Select best structure
//   8. Validate against experimental (if available)
//
// EXPECTED PERFORMANCE:
// - Time: 5-30 seconds per protein (depends on size)
// - RMSD: <5 Å target (vs 63.16 Å Phase 1)
// - Success rate: >90% (no crashes)
func RunUnifiedPipelineV2(config UnifiedPipelineV2Config, experimental *parser.Protein) (*UnifiedPipelineV2Result, error) {
	startTime := time.Now()

	result := &UnifiedPipelineV2Result{}

	if config.Verbose {
		fmt.Printf("=== FoldVedic.ai Unified Pipeline v2.0 ===\n")
		fmt.Printf("Sequence: %s (%d residues)\n", config.Sequence, len(config.Sequence))
		fmt.Printf("\n")
	}

	// PHASE A: PREDICTION & PLANNING
	if config.Verbose {
		fmt.Printf("Phase A: Structural Predictions\n")
	}

	// Step 1: Secondary structure prediction
	var ssPred []prediction.SecondaryStructurePrediction
	if config.UseSSprediction {
		ssConfig := prediction.DefaultPredictionConfig()
		ssConfig.Method = config.SSMethod

		var err error
		ssPred, err = prediction.PredictSecondaryStructure(config.Sequence, ssConfig)
		if err != nil {
			return nil, fmt.Errorf("secondary structure prediction failed: %w", err)
		}

		result.SecondaryStructure = ssPred

		if config.Verbose {
			ssString := prediction.GetSecondaryStructureString(ssPred)
			fmt.Printf("  Secondary Structure: %s\n", ssString)
		}
	}

	// Step 2: Contact map prediction
	var contacts []prediction.ContactPrediction
	if config.UseContactMap {
		var err error
		contacts, err = prediction.PredictContactMap(config.Sequence, config.ContactConfig)
		if err != nil {
			return nil, fmt.Errorf("contact map prediction failed: %w", err)
		}

		result.ContactMap = contacts

		if config.Verbose {
			stats := prediction.GetContactRangeStatistics(contacts)
			fmt.Printf("  Contact Map: %d contacts (Short:%d, Medium:%d, Long:%d)\n",
				stats.Total, stats.ShortRange, stats.MediumRange, stats.LongRange)
		}
	}

	if config.Verbose {
		fmt.Printf("\n")
	}

	// PHASE B: CONFORMATIONAL SAMPLING
	if config.Verbose {
		fmt.Printf("Phase B: Conformational Sampling\n")
	}

	ensemble := make([]*parser.Protein, 0)

	// Initialize base structure from secondary structure prediction
	baseStructure := initializeFromSSPrediction(config.Sequence, ssPred)

	// Method 1: Quaternion slerp sampling
	if config.UseQuaternionSlerp {
		slerpConfig := sampling.DefaultQuaternionSearchConfig()
		slerpConfig.NumSamples = config.NumSamplesPerMethod

		slerpEnsemble, err := sampling.QuaternionGuidedSearch(baseStructure, slerpConfig)
		if err == nil {
			ensemble = append(ensemble, slerpEnsemble...)
			if config.Verbose {
				fmt.Printf("  Quaternion Slerp: %d structures\n", len(slerpEnsemble))
			}
		}
	}

	// Method 2: Monte Carlo sampling
	if config.UseMonteCarlo {
		mcConfig := sampling.DefaultMonteCarloConfig()
		mcConfig.NumSteps = 500 // Quick MC runs
		mcConfig.VedicWeight = 0.3

		mcEnsemble, err := sampling.GenerateMonteCarloEnsemble(baseStructure, mcConfig, config.NumSamplesPerMethod)
		if err == nil {
			ensemble = append(ensemble, mcEnsemble...)
			if config.Verbose {
				fmt.Printf("  Monte Carlo: %d structures\n", len(mcEnsemble))
			}
		}
	}

	// Method 3: Fragment assembly
	if config.UseFragmentAssembly {
		fragmentLib := sampling.NewFragmentLibrary()
		fragConfig := sampling.DefaultFragmentAssemblyConfig()

		fragEnsemble, err := sampling.GenerateFragmentEnsemble(config.Sequence, fragmentLib, fragConfig, config.NumSamplesPerMethod)
		if err == nil {
			ensemble = append(ensemble, fragEnsemble...)
			if config.Verbose {
				fmt.Printf("  Fragment Assembly: %d structures\n", len(fragEnsemble))
			}
		}
	}

	// Method 4: Basin explorer
	if config.UseBasinExplorer {
		basinConfig := sampling.DefaultBasinExplorerConfig()
		basinConfig.SamplesPerBasin = 2 // 2 per basin × ~7 basins = 14 structures

		basinEnsemble, err := sampling.ExploreRamachandranBasins(config.Sequence, basinConfig)
		if err == nil {
			ensemble = append(ensemble, basinEnsemble...)
			if config.Verbose {
				fmt.Printf("  Basin Explorer: %d structures\n", len(basinEnsemble))
			}
		}
	}

	result.TotalSamplesGenerated = len(ensemble)

	if len(ensemble) == 0 {
		return nil, fmt.Errorf("no structures generated during sampling")
	}

	if config.Verbose {
		fmt.Printf("  Total ensemble: %d structures\n", len(ensemble))
		fmt.Printf("\n")
	}

	// PHASE C: ENERGY OPTIMIZATION
	if config.Verbose {
		fmt.Printf("Phase C: Energy Optimization\n")
	}

	bestEnergy := 1e10
	var bestStructure *parser.Protein
	var bestOptResult *optimization.OptimizationResult

	successful := 0

	for i, structure := range ensemble {
		// Optimize structure
		optResult, err := optimization.OptimizeProtein(structure, config.OptimizationConfig)
		if err != nil {
			// Skip failed optimizations
			continue
		}

		successful++

		// Apply Vedic biasing if enabled
		finalEnergy := optResult.FinalEnergy
		if config.UseVedicBiasing {
			angles := geometry.CalculateRamachandran(structure)
			vedicEnergy := prediction.CalculateVedicEnergy(structure, angles, config.VedicBias)
			finalEnergy += config.VedicBias.VedicWeight * vedicEnergy * 1000.0 // Scale to kcal/mol
		}

		// Apply contact restraints if enabled
		if config.UseContactMap && len(contacts) > 0 {
			contactEnergy := prediction.ApplyContactRestraints(structure, contacts, 10.0)
			finalEnergy += contactEnergy
		}

		// Track best
		if finalEnergy < bestEnergy {
			bestEnergy = finalEnergy
			bestStructure = structure
			bestOptResult = optResult
		}

		if config.Verbose && i%5 == 0 {
			fmt.Printf("  Optimized %d/%d structures (best energy: %.2f kcal/mol)\r",
				i+1, len(ensemble), bestEnergy)
		}
	}

	if config.Verbose {
		fmt.Printf("\n")
		fmt.Printf("  Optimization complete: %d/%d successful (%.1f%%)\n",
			successful, len(ensemble), 100.0*float64(successful)/float64(len(ensemble)))
		fmt.Printf("  Best energy: %.2f kcal/mol\n", bestEnergy)
		fmt.Printf("\n")
	}

	result.SuccessRate = float64(successful) / float64(len(ensemble))

	if bestStructure == nil {
		return nil, fmt.Errorf("all optimizations failed")
	}

	// PHASE D: SELECTION & VALIDATION
	if config.Verbose {
		fmt.Printf("Phase D: Final Structure Selection\n")
	}

	result.FinalStructure = bestStructure
	result.FinalAngles = geometry.CalculateRamachandran(bestStructure)
	result.FinalEnergy = bestEnergy
	result.OptimizationResult = bestOptResult

	// Calculate Vedic score
	result.FinalVedicScore = prediction.ScoreProteinVedicHarmonics(
		bestStructure,
		result.FinalAngles,
		config.VedicBias,
	)

	// Generate Vedic report
	result.VedicReport = prediction.GenerateVedicHarmonicReport(
		bestStructure,
		result.FinalAngles,
		ssPred,
		config.VedicBias,
	)

	// Combined score
	result.CombinedScore = (1.0 - config.VedicBias.VedicWeight) * result.FinalEnergy +
		config.VedicBias.VedicWeight * (1.0 - result.FinalVedicScore) * 1000.0

	// Validate against experimental if provided
	if experimental != nil {
		comp := validation.CompareStructures(bestStructure, experimental)
		result.Validation = &comp

		if config.Verbose {
			fmt.Printf("  RMSD: %.2f Å\n", comp.RMSD)
			fmt.Printf("  TM-score: %.3f\n", comp.TMScore)
			fmt.Printf("  GDT_TS: %.3f\n", comp.GDT_TS)
		}

		// Quality score: Harmonic mean of metrics
		rmsdScore := 1.0 / (1.0 + comp.RMSD/10.0)
		tmScore := comp.TMScore
		vedicScore := result.FinalVedicScore

		sumInverses := 1.0/rmsdScore + 1.0/tmScore + 1.0/vedicScore
		result.QualityScore = 3.0 / sumInverses
	} else {
		// No experimental: quality based on energy and Vedic
		energyScore := 1.0 / (1.0 + bestEnergy/10000.0)
		vedicScore := result.FinalVedicScore

		result.QualityScore = 2.0 / (1.0/energyScore + 1.0/vedicScore)
	}

	result.TotalTimeSeconds = time.Since(startTime).Seconds()

	if config.Verbose {
		fmt.Printf("  Vedic Score: %.3f\n", result.FinalVedicScore)
		fmt.Printf("  Quality Score: %.3f\n", result.QualityScore)
		fmt.Printf("\n")
		fmt.Printf("=== Pipeline Complete (%.2f seconds) ===\n", result.TotalTimeSeconds)
	}

	return result, nil
}

// initializeFromSSPrediction creates initial structure from SS prediction
//
// BIOCHEMIST:
// Use predicted helix/sheet regions to set initial (φ, ψ) angles
// - Helix: φ=-60°, ψ=-45°
// - Sheet: φ=-120°, ψ=+120°
// - Coil: Extended φ=-120°, ψ=+120°
func initializeFromSSPrediction(sequence string, ssPred []prediction.SecondaryStructurePrediction) *parser.Protein {
	// Create template protein
	protein := &parser.Protein{
		Name:     "initialized",
		Residues: make([]*parser.Residue, len(sequence)),
	}

	for i := range sequence {
		protein.Residues[i] = &parser.Residue{
			Name:    string(sequence[i]),
			SeqNum:  i + 1,
			ChainID: "A",
		}
	}

	// Build structure from SS-based angles
	angles := make([]geometry.RamachandranAngles, len(sequence))

	for i := range sequence {
		var phi, psi float64

		if i < len(ssPred) {
			switch ssPred[i].PredictedType {
			case prediction.AlphaHelix:
				phi = -60.0 * 3.14159 / 180.0
				psi = -45.0 * 3.14159 / 180.0
			case prediction.BetaSheet:
				phi = -120.0 * 3.14159 / 180.0
				psi = +120.0 * 3.14159 / 180.0
			default: // Coil
				phi = -120.0 * 3.14159 / 180.0
				psi = +120.0 * 3.14159 / 180.0
			}
		} else {
			// Fallback: extended
			phi = -120.0 * 3.14159 / 180.0
			psi = +120.0 * 3.14159 / 180.0
		}

		angles[i] = geometry.RamachandranAngles{Phi: phi, Psi: psi}
	}

	// Build 3D structure from angles
	// WAVE 11.1: Use quaternion-based coordinate builder (NOVEL!)
	structure, err := geometry.BuildProteinFromAngles(sequence, angles)
	if err != nil {
		// Fallback to simple builder if quaternion method fails
		return initializeFallback(sequence)
	}

	return structure
}

// initializeFallback creates extended chain if coordinate builder fails
func initializeFallback(sequence string) *parser.Protein {
	// Simple extended chain as last resort
	protein := &parser.Protein{
		Name:     "fallback",
		Residues: make([]*parser.Residue, len(sequence)),
		Atoms:    make([]*parser.Atom, 0, len(sequence)*4),
	}

	atomSerial := 1
	x := 0.0

	for i := range sequence {
		res := &parser.Residue{
			Name:    string(sequence[i]),
			SeqNum:  i + 1,
			ChainID: "A",
		}

		// N, CA, C, O along X-axis
		res.N = &parser.Atom{Serial: atomSerial, Name: "N", ResName: string(sequence[i]),
			ChainID: "A", ResSeq: i + 1, X: x, Y: 0, Z: 0, Element: "N"}
		protein.Atoms = append(protein.Atoms, res.N)
		atomSerial++
		x += 1.46

		res.CA = &parser.Atom{Serial: atomSerial, Name: "CA", ResName: string(sequence[i]),
			ChainID: "A", ResSeq: i + 1, X: x, Y: 0, Z: 0, Element: "C"}
		protein.Atoms = append(protein.Atoms, res.CA)
		atomSerial++
		x += 1.52

		res.C = &parser.Atom{Serial: atomSerial, Name: "C", ResName: string(sequence[i]),
			ChainID: "A", ResSeq: i + 1, X: x, Y: 0, Z: 0, Element: "C"}
		protein.Atoms = append(protein.Atoms, res.C)
		atomSerial++

		res.O = &parser.Atom{Serial: atomSerial, Name: "O", ResName: string(sequence[i]),
			ChainID: "A", ResSeq: i + 1, X: x, Y: 1.23, Z: 0, Element: "O"}
		protein.Atoms = append(protein.Atoms, res.O)
		atomSerial++
		x += 1.33

		protein.Residues[i] = res
	}

	return protein
}

// QuickFold provides simple interface for protein folding
//
// CONVENIENCE FUNCTION:
// Folds protein using all default Phase 2 enhancements
func QuickFold(sequence string, verbose bool) (*UnifiedPipelineV2Result, error) {
	config := DefaultUnifiedPipelineV2Config(sequence)
	config.Verbose = verbose
	return RunUnifiedPipelineV2(config, nil)
}
