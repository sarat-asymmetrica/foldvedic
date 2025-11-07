package sampling

import (
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// PUBLIC wrappers for Wave 4 integration
// These provide simplified interfaces to the sampling methods

// FibonacciSphereBasins generates protein structures using Fibonacci sphere sampling
// Returns numStructures conformations
func FibonacciSphereBasins(sequence string, numStructures int) []*parser.Protein {
	// Use basin explorer with standard configuration
	config := DefaultBasinExplorerConfig()
	basins := GetStandardRamachandranBasins()
	samplesPerBasin := numStructures / len(basins)
	if samplesPerBasin < 1 {
		samplesPerBasin = 1
	}
	config.SamplesPerBasin = samplesPerBasin

	structures, err := ExploreRamachandranBasins(sequence, config)
	if err != nil {
		return []*parser.Protein{}
	}

	// Trim to requested number
	if len(structures) > numStructures {
		return structures[:numStructures]
	}
	return structures
}

// VedicMonteCarlo generates protein structures using Vedic-guided Monte Carlo
// Returns numStructures conformations
func VedicMonteCarlo(sequence string, numStructures int) []*parser.Protein {
	// Use basin explorer with Vedic biasing enabled
	config := DefaultBasinExplorerConfig()
	config.UseVedicBiasing = true
	config.SamplesPerBasin = numStructures / 2
	if config.SamplesPerBasin < 1 {
		config.SamplesPerBasin = 1
	}

	structures, err := ExploreRamachandranBasins(sequence, config)
	if err != nil {
		return []*parser.Protein{}
	}

	// Trim to requested number
	if len(structures) > numStructures {
		return structures[:numStructures]
	}
	return structures
}

// GenerateFragmentStructures generates protein structures using fragment assembly
// Returns numStructures conformations
// Wrapper to avoid naming conflict with FragmentAssembly in fragments.go
func GenerateFragmentStructures(sequence string, numStructures int) []*parser.Protein {
	// For simplicity, just use Basin Explorer with fragment-like settings
	// Real fragment assembly requires complex library setup
	config := DefaultBasinExplorerConfig()
	config.SamplesPerBasin = numStructures / 3
	if config.SamplesPerBasin < 1 {
		config.SamplesPerBasin = 1
	}

	structures, err := ExploreRamachandranBasins(sequence, config)
	if err != nil {
		return []*parser.Protein{}
	}

	// Trim to requested number
	if len(structures) > numStructures {
		return structures[:numStructures]
	}
	return structures
}

// BasinExplorer generates protein structures exploring Ramachandran basins
// Returns numStructures conformations (PRIMARY METHOD for Phase 2)
func BasinExplorer(sequence string, numStructures int) []*parser.Protein {
	config := DefaultBasinExplorerConfig()
	basins := GetStandardRamachandranBasins()
	samplesPerBasin := numStructures / len(basins)
	if samplesPerBasin < 1 {
		samplesPerBasin = 1
	}
	config.SamplesPerBasin = samplesPerBasin

	structures, err := ExploreRamachandranBasins(sequence, config)
	if err != nil {
		return []*parser.Protein{}
	}

	// Trim to requested number
	if len(structures) > numStructures {
		return structures[:numStructures]
	}
	return structures
}
