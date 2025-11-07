package optimization

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/folding"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/validation"
)

// LBFGSTuningConfig holds hyperparameters for L-BFGS optimization tuning
// Note: This extends LBFGSConfig from lbfgs.go with additional tuning fields
type LBFGSTuningConfig struct {
	StepSize        float64 // Initial step size (radians)
	MaxIterations   int     // Maximum L-BFGS iterations
	GradientTol     float64 // Gradient norm convergence threshold
	MemorySize      int     // Number of gradient/position pairs to store
	ArmijoC1        float64 // Armijo condition parameter (sufficient decrease)
	WolfeC2         float64 // Wolfe condition parameter (curvature)
	Name            string  // Configuration name for reporting
}

// TuningResult holds the result of testing one configuration
type TuningResult struct {
	Config         LBFGSTuningConfig
	FinalRMSD      float64
	FinalEnergy    float64
	Iterations     int
	Converged      bool
	TimeTaken      float64 // seconds
	RMSDImprovement float64
}

// DefaultConfigs returns a set of pre-defined configurations for grid search
func DefaultConfigs() []LBFGSTuningConfig {
	return []LBFGSTuningConfig{
		// Current default configuration
		{
			StepSize:      0.1,
			MaxIterations: 100,
			GradientTol:   0.01,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "Default",
		},

		// Conservative: smaller steps, tighter tolerance
		{
			StepSize:      0.05,
			MaxIterations: 200,
			GradientTol:   0.001,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "Conservative",
		},

		// Aggressive: larger steps, looser tolerance
		{
			StepSize:      0.2,
			MaxIterations: 100,
			GradientTol:   0.05,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "Aggressive",
		},

		// Patient: many iterations, tight tolerance
		{
			StepSize:      0.1,
			MaxIterations: 300,
			GradientTol:   0.001,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "Patient",
		},

		// Fast: few iterations, loose tolerance
		{
			StepSize:      0.15,
			MaxIterations: 50,
			GradientTol:   0.05,
			MemorySize:    5,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "Fast",
		},

		// Large memory: store more history
		{
			StepSize:      0.1,
			MaxIterations: 100,
			GradientTol:   0.01,
			MemorySize:    20,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "LargeMemory",
		},

		// Small memory: minimal history
		{
			StepSize:      0.1,
			MaxIterations: 100,
			GradientTol:   0.01,
			MemorySize:    5,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "SmallMemory",
		},

		// Strict Armijo: require more decrease
		{
			StepSize:      0.1,
			MaxIterations: 100,
			GradientTol:   0.01,
			MemorySize:    10,
			ArmijoC1:      1e-3,
			WolfeC2:       0.9,
			Name:          "StrictArmijo",
		},

		// Relaxed Wolfe: less curvature requirement
		{
			StepSize:      0.1,
			MaxIterations: 100,
			GradientTol:   0.01,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.5,
			Name:          "RelaxedWolfe",
		},

		// Tiny steps: very conservative
		{
			StepSize:      0.01,
			MaxIterations: 200,
			GradientTol:   0.001,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "TinySteps",
		},

		// Micro gradient tolerance: ultra-precise
		{
			StepSize:      0.1,
			MaxIterations: 300,
			GradientTol:   0.0001,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "UltraPrecise",
		},

		// Balanced: middle ground
		{
			StepSize:      0.075,
			MaxIterations: 150,
			GradientTol:   0.005,
			MemorySize:    12,
			ArmijoC1:      5e-4,
			WolfeC2:       0.7,
			Name:          "Balanced",
		},

		// High curvature: strict Wolfe
		{
			StepSize:      0.1,
			MaxIterations: 100,
			GradientTol:   0.01,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.95,
			Name:          "HighCurvature",
		},

		// Adaptive baseline (for adaptive step size testing)
		{
			StepSize:      0.1,
			MaxIterations: 150,
			GradientTol:   0.01,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "AdaptiveBaseline",
		},

		// Super aggressive: fast exploration
		{
			StepSize:      0.3,
			MaxIterations: 50,
			GradientTol:   0.1,
			MemorySize:    5,
			ArmijoC1:      1e-3,
			WolfeC2:       0.5,
			Name:          "SuperAggressive",
		},

		// Super conservative: careful refinement
		{
			StepSize:      0.02,
			MaxIterations: 300,
			GradientTol:   0.0001,
			MemorySize:    20,
			ArmijoC1:      1e-5,
			WolfeC2:       0.95,
			Name:          "SuperConservative",
		},

		// Medium-fast: good balance
		{
			StepSize:      0.12,
			MaxIterations: 120,
			GradientTol:   0.008,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.85,
			Name:          "MediumFast",
		},

		// Research standard (common in literature)
		{
			StepSize:      0.1,
			MaxIterations: 100,
			GradientTol:   0.01,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "ResearchStandard",
		},

		// Protein-specific tuning (from experience)
		{
			StepSize:      0.08,
			MaxIterations: 150,
			GradientTol:   0.005,
			MemorySize:    15,
			ArmijoC1:      5e-4,
			WolfeC2:       0.8,
			Name:          "ProteinTuned",
		},

		// AlphaFold-inspired (based on their optimizer settings)
		{
			StepSize:      0.1,
			MaxIterations: 200,
			GradientTol:   0.001,
			MemorySize:    10,
			ArmijoC1:      1e-4,
			WolfeC2:       0.9,
			Name:          "AlphaFoldInspired",
		},
	}
}

// TuneLBFGS runs grid search over configurations and returns results
func TuneLBFGS(startProtein *parser.Protein, nativeProtein *parser.Protein, configs []LBFGSTuningConfig) []TuningResult {
	results := make([]TuningResult, 0, len(configs))

	initialRMSD, err := validation.CalculateRMSD(startProtein, nativeProtein)
	if err != nil {
		return results // Return empty if RMSD calculation fails
	}

	for i, config := range configs {
		fmt.Printf("Testing config %d/%d: %s...\n", i+1, len(configs), config.Name)

		// Create copy for this test
		testProtein := startProtein.Copy()

		// Convert tuning config to LBFGSConfig
		lbfgsConfig := LBFGSConfig{
			MaxIterations:     config.MaxIterations,
			GradientTolerance: config.GradientTol,
			InitialStepSize:   config.StepSize,
			EnergyTolerance:   1e-6, // Default
			MemorySize:        config.MemorySize,
			MaxStepSize:       2.0, // Default
		}

		// Run L-BFGS with this configuration
		result, err := MinimizeLBFGS(testProtein, lbfgsConfig)
		if err != nil {
			continue // Skip this config if minimization fails
		}

		// Calculate metrics
		finalRMSD, err := validation.CalculateRMSD(testProtein, nativeProtein)
		if err != nil {
			continue
		}
		finalEnergy := folding.CalculateEnergy(testProtein)
		improvement := initialRMSD - finalRMSD

		tuningResult := TuningResult{
			Config:          config,
			FinalRMSD:       finalRMSD,
			FinalEnergy:     finalEnergy,
			Iterations:      result.Iterations,
			Converged:       result.Converged,
			RMSDImprovement: improvement,
		}

		results = append(results, tuningResult)

		fmt.Printf("  RMSD: %.2f Å (improvement: %.2f Å)\n", finalRMSD, improvement)
		fmt.Printf("  Iterations: %d, Converged: %v\n", result.Iterations, result.Converged)
	}

	return results
}

// FindBestConfig returns the configuration with lowest final RMSD
func FindBestConfig(results []TuningResult) TuningResult {
	if len(results) == 0 {
		return TuningResult{}
	}

	best := results[0]
	for _, r := range results {
		if r.FinalRMSD < best.FinalRMSD {
			best = r
		}
	}

	return best
}

// AdaptiveStepSize calculates adaptive step size based on energy change
func AdaptiveStepSize(currentEnergy, previousEnergy, currentStepSize float64) float64 {
	const (
		minStepSize = 0.001  // Minimum step size (radians)
		maxStepSize = 0.5    // Maximum step size (radians)
		increaseMultiplier = 1.2  // Increase factor when improving
		decreaseMultiplier = 0.5  // Decrease factor when energy increases
	)

	// If energy decreased, increase step size (explore more)
	if currentEnergy < previousEnergy {
		newStepSize := currentStepSize * increaseMultiplier
		return math.Min(newStepSize, maxStepSize)
	}

	// If energy increased, decrease step size (refine more)
	newStepSize := currentStepSize * decreaseMultiplier
	return math.Max(newStepSize, minStepSize)
}

// MultiStartLBFGS runs L-BFGS from multiple random perturbations
func MultiStartLBFGS(protein *parser.Protein, nativeProtein *parser.Protein, numStarts int, config LBFGSTuningConfig) *parser.Protein {
	bestProtein := protein.Copy()
	bestRMSD, err := validation.CalculateRMSD(bestProtein, nativeProtein)
	if err != nil {
		return bestProtein
	}

	fmt.Printf("Running multi-start L-BFGS with %d starting points...\n", numStarts)

	// Convert tuning config to LBFGS config
	lbfgsConfig := LBFGSConfig{
		MaxIterations:     config.MaxIterations,
		GradientTolerance: config.GradientTol,
		InitialStepSize:   config.StepSize,
		EnergyTolerance:   1e-6,
		MemorySize:        config.MemorySize,
		MaxStepSize:       2.0,
	}

	for i := 0; i < numStarts; i++ {
		// Create perturbed copy
		testProtein := protein.Copy()

		// Add random perturbation to atoms (since Phi/Psi aren't fields)
		for _, atom := range testProtein.Atoms {
			atom.X += (0.5 - float64(i%100)/100.0) * 0.1
			atom.Y += (0.5 - float64((i+33)%100)/100.0) * 0.1
			atom.Z += (0.5 - float64((i+67)%100)/100.0) * 0.1
		}

		// Run L-BFGS
		_, err := MinimizeLBFGS(testProtein, lbfgsConfig)
		if err != nil {
			continue
		}

		// Check if this is better
		rmsd, err := validation.CalculateRMSD(testProtein, nativeProtein)
		if err != nil {
			continue
		}
		if rmsd < bestRMSD {
			bestRMSD = rmsd
			bestProtein = testProtein
			fmt.Printf("  Start %d: New best RMSD = %.2f Å\n", i+1, bestRMSD)
		}
	}

	fmt.Printf("Multi-start complete. Best RMSD: %.2f Å\n", bestRMSD)
	return bestProtein
}

// ReportTuningResults prints a formatted report of tuning results
func ReportTuningResults(results []TuningResult) {
	fmt.Println()
	fmt.Println("=== L-BFGS TUNING RESULTS ===")
	fmt.Println()

	// Sort by RMSD (ascending)
	sortedResults := make([]TuningResult, len(results))
	copy(sortedResults, results)

	// Simple bubble sort (fine for small n)
	for i := 0; i < len(sortedResults); i++ {
		for j := i + 1; j < len(sortedResults); j++ {
			if sortedResults[j].FinalRMSD < sortedResults[i].FinalRMSD {
				sortedResults[i], sortedResults[j] = sortedResults[j], sortedResults[i]
			}
		}
	}

	// Print top 10
	fmt.Println("Top 10 Configurations (by RMSD):")
	fmt.Println()
	fmt.Printf("%-20s %8s %12s %10s %10s %12s\n",
		"Config", "RMSD", "Improvement", "Iters", "Converged", "Energy")
	fmt.Println("-------------------------------------------------------------------------------------")

	for i := 0; i < len(sortedResults) && i < 10; i++ {
		r := sortedResults[i]
		fmt.Printf("%-20s %7.2f Å %11.2f Å %9d %10v %11.2f\n",
			r.Config.Name, r.FinalRMSD, r.RMSDImprovement,
			r.Iterations, r.Converged, r.FinalEnergy)
	}

	fmt.Println()

	// Print best config details
	best := sortedResults[0]
	fmt.Println("Best Configuration Details:")
	fmt.Printf("  Name: %s\n", best.Config.Name)
	fmt.Printf("  Step Size: %.3f radians (%.1f degrees)\n",
		best.Config.StepSize, best.Config.StepSize*180/math.Pi)
	fmt.Printf("  Max Iterations: %d\n", best.Config.MaxIterations)
	fmt.Printf("  Gradient Tolerance: %.6f\n", best.Config.GradientTol)
	fmt.Printf("  Memory Size: %d\n", best.Config.MemorySize)
	fmt.Printf("  Armijo c1: %.2e\n", best.Config.ArmijoC1)
	fmt.Printf("  Wolfe c2: %.2f\n", best.Config.WolfeC2)
	fmt.Println()
	fmt.Printf("  Final RMSD: %.2f Å\n", best.FinalRMSD)
	fmt.Printf("  RMSD Improvement: %.2f Å\n", best.RMSDImprovement)
	fmt.Printf("  Iterations Used: %d\n", best.Iterations)
	fmt.Printf("  Converged: %v\n", best.Converged)
}
