// Package main - Quick Benchmark for Unified Pipeline v2
//
// WAVE 10.2: Quick Benchmark Test
// Validates that Phase 2 pipeline works end-to-end
package main

import (
	"fmt"
	"time"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/pipeline"
)

func main() {
	fmt.Println("=== FoldVedic.ai Phase 2 Quick Benchmark ===")
	fmt.Println()

	// Test sequences of increasing complexity
	testCases := []struct {
		name     string
		sequence string
	}{
		{"Tiny (6aa)", "GACDEF"},
		{"Small (10aa)", "GACDEFGHIK"},
		{"Medium (15aa)", "GACDEFGHIKLMNPQ"},
	}

	for _, tc := range testCases {
		fmt.Printf("Testing: %s\n", tc.name)
		fmt.Printf("Sequence: %s\n", tc.sequence)
		fmt.Println()

		start := time.Now()

		// Run pipeline with default config
		result, err := pipeline.QuickFold(tc.sequence, true)

		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			fmt.Println()
			continue
		}

		// Display results
		fmt.Printf("\n")
		fmt.Printf("RESULTS:\n")
		fmt.Printf("  Samples generated: %d\n", result.TotalSamplesGenerated)
		fmt.Printf("  Success rate: %.1f%%\n", result.SuccessRate*100)
		fmt.Printf("  Final energy: %.2f kcal/mol\n", result.FinalEnergy)
		fmt.Printf("  Vedic score: %.3f\n", result.FinalVedicScore)
		fmt.Printf("  Quality score: %.3f\n", result.QualityScore)
		fmt.Printf("  Time: %.2f seconds\n", elapsed.Seconds())
		fmt.Println()

		// Secondary structure
		if len(result.SecondaryStructure) > 0 {
			ssString := ""
			for _, ss := range result.SecondaryStructure {
				switch ss.PredictedType {
				case 1: // AlphaHelix
					ssString += "H"
				case 2: // BetaSheet
					ssString += "E"
				default: // Coil
					ssString += "C"
				}
			}
			fmt.Printf("  Predicted SS: %s\n", ssString)
		}

		// Contact map
		if len(result.ContactMap) > 0 {
			fmt.Printf("  Predicted contacts: %d\n", len(result.ContactMap))
		}

		// Optimization
		if result.OptimizationResult != nil {
			fmt.Printf("  Optimization: %d iterations, %.2f energy reduction\n",
				result.OptimizationResult.Iterations,
				result.OptimizationResult.InitialEnergy-result.OptimizationResult.FinalEnergy)
		}

		fmt.Println()
		fmt.Println("----------------------------------------")
		fmt.Println()
	}

	fmt.Println("=== Benchmark Complete ===")
}
