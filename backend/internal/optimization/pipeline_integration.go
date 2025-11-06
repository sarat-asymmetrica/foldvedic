// Package optimization - Pipeline Integration
//
// WAVE 8.3: Pipeline Integration with Advanced Optimization
// Upgrades folding pipeline from 100-step steepest descent to 1000+ step advanced optimization
//
// BIOCHEMIST: Provides sufficient minimization for native-like structures
// PHYSICIST: Combines global search (SA) with local refinement (L-BFGS)
// MATHEMATICIAN: Adaptive step budgets based on protein size
// ETHICIST: Transparent optimization strategy, reproducible results
//
// UPGRADE PATH:
// Phase 1 (v0.1): 100 steps steepest descent
// Phase 2 (v0.2): 1000+ steps SA + L-BFGS hybrid
// Expected improvement: 20-30% better convergence
package optimization

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// OptimizationStrategy defines which optimization method to use
type OptimizationStrategy string

const (
	// StrategyLBFGS: L-BFGS only (fast, local optimization)
	StrategyLBFGS OptimizationStrategy = "lbfgs"

	// StrategySimulatedAnnealing: SA only (global search, slower)
	StrategySimulatedAnnealing OptimizationStrategy = "simulated_annealing"

	// StrategyHybrid: SA + L-BFGS (best of both worlds, recommended)
	StrategyHybrid OptimizationStrategy = "hybrid"

	// StrategySteepestDescent: Legacy Phase 1 method (for comparison)
	StrategySteepestDescent OptimizationStrategy = "steepest_descent"
)

// AdaptiveOptimizationConfig holds parameters for adaptive optimization
type AdaptiveOptimizationConfig struct {
	// Optimization strategy
	Strategy OptimizationStrategy

	// Adaptive step budget based on protein size
	UseAdaptiveBudget bool
	BaseSteps         int // Base number of steps for reference protein (76 residues = ubiquitin)

	// Energy and gradient tolerances
	EnergyTolerance   float64
	GradientTolerance float64

	// Force field cutoffs
	VdWCutoff  float64
	ElecCutoff float64

	// Verbose logging
	Verbose bool
}

// DefaultAdaptiveOptimizationConfig returns recommended parameters for Phase 2
func DefaultAdaptiveOptimizationConfig() AdaptiveOptimizationConfig {
	return AdaptiveOptimizationConfig{
		Strategy:          StrategyHybrid,  // SA + L-BFGS
		UseAdaptiveBudget: true,            // Scale steps with protein size
		BaseSteps:         1000,            // 1000 steps for 76-residue protein
		EnergyTolerance:   0.01,            // 0.01 kcal/mol
		GradientTolerance: 0.1,             // 0.1 kcal/(mol·Å)
		VdWCutoff:         10.0,
		ElecCutoff:        12.0,
		Verbose:           false,
	}
}

// OptimizationResult holds unified optimization results
type OptimizationResult struct {
	// Which strategy was used
	Strategy OptimizationStrategy

	// Energy statistics
	InitialEnergy float64
	FinalEnergy   float64
	EnergyChange  float64

	// Optimization statistics
	Iterations          int
	FunctionEvaluations int
	Converged           bool
	Reason              string

	// Strategy-specific results
	LBFGSResult *LBFGSResult
	SAResult    *SimulatedAnnealingResult
}

// OptimizeProtein performs adaptive optimization based on protein size and strategy
//
// ALGORITHM:
// 1. Determine optimization budget (adaptive to protein size)
// 2. Select optimization strategy (LBFGS, SA, or Hybrid)
// 3. Execute optimization
// 4. Return unified results
//
// ADAPTIVE BUDGET:
// Steps = BaseSteps × sqrt(N / N_ref)
// where N = number of residues, N_ref = 76 (ubiquitin reference)
//
// BIOCHEMIST:
// Larger proteins need more steps but not linearly
// sqrt scaling balances thoroughness vs computational cost
func OptimizeProtein(protein *parser.Protein, config AdaptiveOptimizationConfig) (*OptimizationResult, error) {
	if protein == nil {
		return nil, fmt.Errorf("protein is nil")
	}

	result := &OptimizationResult{
		Strategy: config.Strategy,
	}

	// Calculate adaptive step budget
	numSteps := calculateAdaptiveSteps(protein, config)

	if config.Verbose {
		fmt.Printf("Optimizing %s (%d residues) with strategy: %s\n",
			protein.Name, len(protein.Residues), config.Strategy)
		fmt.Printf("  Optimization budget: %d steps\n", numSteps)
	}

	// Calculate initial energy
	initialEnergy := evaluateEnergy(protein, LBFGSConfig{
		VdWCutoff:  config.VdWCutoff,
		ElecCutoff: config.ElecCutoff,
	})
	result.InitialEnergy = initialEnergy

	// Execute optimization based on strategy
	switch config.Strategy {
	case StrategyLBFGS:
		lbfgsConfig := DefaultLBFGSConfig()
		lbfgsConfig.MaxIterations = numSteps
		lbfgsConfig.EnergyTolerance = config.EnergyTolerance
		lbfgsConfig.GradientTolerance = config.GradientTolerance
		lbfgsConfig.VdWCutoff = config.VdWCutoff
		lbfgsConfig.ElecCutoff = config.ElecCutoff
		lbfgsConfig.Verbose = config.Verbose

		lbfgsResult, err := MinimizeLBFGS(protein, lbfgsConfig)
		if err != nil {
			return nil, fmt.Errorf("L-BFGS optimization failed: %w", err)
		}

		result.LBFGSResult = lbfgsResult
		result.FinalEnergy = lbfgsResult.FinalEnergy
		result.EnergyChange = lbfgsResult.EnergyChange
		result.Iterations = lbfgsResult.Iterations
		result.FunctionEvaluations = lbfgsResult.FunctionEvaluations
		result.Converged = lbfgsResult.Converged
		result.Reason = lbfgsResult.Reason

	case StrategySimulatedAnnealing:
		saConfig := DefaultSimulatedAnnealingConfig()
		saConfig.NumSteps = numSteps
		saConfig.VdWCutoff = config.VdWCutoff
		saConfig.ElecCutoff = config.ElecCutoff
		saConfig.Verbose = config.Verbose
		saConfig.UseLBFGSRefinement = false // Pure SA

		saResult, err := SimulatedAnnealing(protein, saConfig)
		if err != nil {
			return nil, fmt.Errorf("simulated annealing failed: %w", err)
		}

		result.SAResult = saResult
		result.FinalEnergy = saResult.FinalEnergy
		result.EnergyChange = saResult.EnergyChange
		result.Iterations = saResult.Steps
		result.FunctionEvaluations = saResult.FunctionEvaluations
		result.Converged = saResult.Converged
		result.Reason = saResult.Reason

	case StrategyHybrid:
		// Hybrid: SA (70% of budget) + L-BFGS (30% of budget)
		saSteps := int(float64(numSteps) * 0.7)
		lbfgsSteps := int(float64(numSteps) * 0.3)

		saConfig := DefaultSimulatedAnnealingConfig()
		saConfig.NumSteps = saSteps
		saConfig.VdWCutoff = config.VdWCutoff
		saConfig.ElecCutoff = config.ElecCutoff
		saConfig.Verbose = config.Verbose
		saConfig.UseLBFGSRefinement = true
		saConfig.RefinementThreshold = 50.0
		saConfig.LBFGSSteps = 50

		lbfgsConfig := DefaultLBFGSConfig()
		lbfgsConfig.MaxIterations = lbfgsSteps
		lbfgsConfig.EnergyTolerance = config.EnergyTolerance
		lbfgsConfig.GradientTolerance = config.GradientTolerance
		lbfgsConfig.VdWCutoff = config.VdWCutoff
		lbfgsConfig.ElecCutoff = config.ElecCutoff
		lbfgsConfig.Verbose = config.Verbose

		// Execute hybrid optimization
		hybridResult, err := HybridOptimization(protein, saConfig, lbfgsConfig)
		if err != nil {
			return nil, fmt.Errorf("hybrid optimization failed: %w", err)
		}

		result.SAResult = hybridResult
		result.FinalEnergy = hybridResult.FinalEnergy
		result.EnergyChange = hybridResult.EnergyChange
		result.Iterations = hybridResult.Steps
		result.FunctionEvaluations = hybridResult.FunctionEvaluations
		result.Converged = hybridResult.Converged
		result.Reason = hybridResult.Reason

	case StrategySteepestDescent:
		// Legacy Phase 1 method (for comparison)
		// Not implementing here - use existing physics.MinimizeEnergy
		return nil, fmt.Errorf("steepest descent strategy should use physics.MinimizeEnergy directly")

	default:
		return nil, fmt.Errorf("unknown optimization strategy: %s", config.Strategy)
	}

	if config.Verbose {
		fmt.Printf("Optimization complete: %.2f → %.2f kcal/mol (Δ = %.2f)\n",
			result.InitialEnergy, result.FinalEnergy, result.EnergyChange)
		fmt.Printf("  Converged: %v (%s)\n", result.Converged, result.Reason)
	}

	return result, nil
}

// calculateAdaptiveSteps determines optimization budget based on protein size
//
// FORMULA:
// Steps = BaseSteps × sqrt(N / N_ref)
//
// where:
//   N = number of residues in protein
//   N_ref = 76 (ubiquitin reference size)
//   BaseSteps = 1000 (for Phase 2)
//
// EXAMPLES:
//   - 20 residues (Trp-cage): 1000 × sqrt(20/76) = 513 steps
//   - 76 residues (Ubiquitin): 1000 × sqrt(76/76) = 1000 steps
//   - 150 residues (Myoglobin): 1000 × sqrt(150/76) = 1404 steps
//
// MATHEMATICIAN:
// Sublinear scaling: O(sqrt(N)) instead of O(N)
// Balances thoroughness with computational cost
func calculateAdaptiveSteps(protein *parser.Protein, config AdaptiveOptimizationConfig) int {
	if !config.UseAdaptiveBudget {
		return config.BaseSteps
	}

	numResidues := len(protein.Residues)
	if numResidues == 0 {
		return config.BaseSteps
	}

	// Reference protein: ubiquitin (76 residues)
	const referenceSize = 76.0

	// Adaptive scaling: Steps ∝ sqrt(N)
	sizeFactor := math.Sqrt(float64(numResidues) / referenceSize)
	steps := int(float64(config.BaseSteps) * sizeFactor)

	// Ensure minimum and maximum bounds
	const minSteps = 500
	const maxSteps = 5000

	if steps < minSteps {
		steps = minSteps
	}
	if steps > maxSteps {
		steps = maxSteps
	}

	return steps
}

// CompareOptimizationStrategies compares all strategies on the same protein
//
// BENCHMARK TOOL:
// Useful for validating that Phase 2 methods outperform Phase 1
// Run on test proteins to demonstrate improvement
func CompareOptimizationStrategies(protein *parser.Protein) (map[OptimizationStrategy]*OptimizationResult, error) {
	strategies := []OptimizationStrategy{
		StrategyLBFGS,
		StrategySimulatedAnnealing,
		StrategyHybrid,
	}

	results := make(map[OptimizationStrategy]*OptimizationResult)

	for _, strategy := range strategies {
		// Clone protein for independent optimization
		proteinCopy := cloneProtein(protein)

		config := DefaultAdaptiveOptimizationConfig()
		config.Strategy = strategy
		config.BaseSteps = 500 // Shorter for comparison
		config.Verbose = false

		result, err := OptimizeProtein(proteinCopy, config)
		if err != nil {
			fmt.Printf("Strategy %s failed: %v\n", strategy, err)
			continue
		}

		results[strategy] = result

		fmt.Printf("Strategy %s:\n", strategy)
		fmt.Printf("  Energy: %.2f → %.2f (Δ = %.2f kcal/mol)\n",
			result.InitialEnergy, result.FinalEnergy, result.EnergyChange)
		fmt.Printf("  Iterations: %d, Converged: %v\n", result.Iterations, result.Converged)
		fmt.Printf("\n")
	}

	return results, nil
}

// GetRecommendedStrategy returns best strategy for given protein size
//
// HEURISTIC:
// - Small proteins (<50 residues): L-BFGS (fast convergence)
// - Medium proteins (50-150 residues): Hybrid (balanced)
// - Large proteins (>150 residues): Simulated Annealing (avoid local minima)
func GetRecommendedStrategy(numResidues int) OptimizationStrategy {
	if numResidues < 50 {
		return StrategyLBFGS
	} else if numResidues < 150 {
		return StrategyHybrid
	} else {
		return StrategySimulatedAnnealing
	}
}

// QuickOptimize provides simple interface for protein optimization
//
// CONVENIENCE FUNCTION:
// Uses recommended strategy and adaptive budget
// For custom control, use OptimizeProtein directly
func QuickOptimize(protein *parser.Protein, verbose bool) (*OptimizationResult, error) {
	config := DefaultAdaptiveOptimizationConfig()
	config.Strategy = GetRecommendedStrategy(len(protein.Residues))
	config.Verbose = verbose

	return OptimizeProtein(protein, config)
}
