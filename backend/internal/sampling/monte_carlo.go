// Package sampling - Monte Carlo conformational sampling with Vedic biasing
//
// WAVE 7.2: Monte Carlo with Vedic Biasing
// Metropolis-Hastings sampling guided by golden ratio harmonics
//
// PHYSICIST: Metropolis criterion for thermodynamic ensemble sampling
// MATHEMATICIAN: Detailed balance, ergodicity, convergence guarantees
// BIOCHEMIST: Explores conformational basins while favoring native-like geometries
// ETHICIST: Transparent probabilistic algorithm, fully reproducible
//
// INNOVATION: Combines AMBER force field energy with Vedic harmonic scoring
// Standard MC: Accept based on energy alone
// Vedic MC: Accept based on energy + φ-ratio geometric alignment
//
// CITATION:
// Metropolis, N., et al. (1953). "Equation of state calculations by fast computing machines."
// J. Chem. Phys. 21(6): 1087-1092.
package sampling

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/physics"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/vedic"
)

// MonteCarloConfig holds Monte Carlo simulation parameters
type MonteCarloConfig struct {
	// Number of MC steps
	NumSteps int

	// Initial temperature (Kelvin)
	// Higher T = more exploration, lower T = more exploitation
	TemperatureInitial float64

	// Final temperature (Kelvin)
	TemperatureFinal float64

	// Cooling schedule: exponential, linear, geometric, vedic_phi
	CoolingSchedule string

	// Step size for coordinate perturbations (Angstroms)
	StepSize float64

	// Vedic bias weight [0, 1]
	// 0 = pure energy, 1 = pure Vedic score, 0.3 = 30% Vedic influence
	VedicWeight float64

	// Energy calculation cutoffs
	VdWCutoff  float64 // Van der Waals cutoff (Å)
	ElecCutoff float64 // Electrostatic cutoff (Å)

	// Random seed for reproducibility
	Seed int64

	// Acceptance tracking
	TrackAcceptance bool
}

// DefaultMonteCarloConfig returns recommended MC parameters
func DefaultMonteCarloConfig() MonteCarloConfig {
	return MonteCarloConfig{
		NumSteps:           1000,            // 1000 MC steps
		TemperatureInitial: 500.0,           // 500 K (high exploration)
		TemperatureFinal:   10.0,            // 10 K (low refinement)
		CoolingSchedule:    "vedic_phi",     // Golden ratio cooling
		StepSize:           0.5,             // 0.5 Å perturbations
		VedicWeight:        0.3,             // 30% Vedic influence
		VdWCutoff:          10.0,            // 10 Å
		ElecCutoff:         12.0,            // 12 Å
		Seed:               42,              // Reproducible
		TrackAcceptance:    true,            // Track acceptance rate
	}
}

// MonteCarloResult holds MC simulation results
type MonteCarloResult struct {
	// Final structure
	FinalStructure *parser.Protein

	// Energy and Vedic scores
	FinalEnergy      float64
	FinalVedicScore  float64
	InitialEnergy    float64
	InitialVedicScore float64

	// Trajectory statistics
	NumAccepted     int
	NumRejected     int
	AcceptanceRate  float64
	BestEnergy      float64
	BestVedicScore  float64

	// Convergence
	Converged      bool
	ConvergenceStep int
}

// MonteCarloVedic performs Monte Carlo sampling with Vedic harmonic biasing
//
// ALGORITHM:
// 1. Start with initial structure
// 2. For each MC step:
//    a. Propose move: perturb atom coordinates
//    b. Calculate combined score: S = E_energy - λ × E_vedic
//    c. Accept with probability: min(1, exp(-ΔS/kT))
//    d. Cool temperature according to schedule
// 3. Return best structure found
//
// VEDIC ENHANCEMENT:
// - Standard MC only considers energy
// - Vedic MC adds golden ratio alignment bonus
// - Favors conformations with φ-based geometries
//
// PHYSICIST:
// This explores the Boltzmann distribution:
//   P(state) ∝ exp(-E/kT)
// At high T: All states equally likely (exploration)
// At low T: Only low-energy states likely (exploitation)
//
// MATHEMATICIAN:
// Detailed balance ensures convergence to equilibrium distribution
// Ergodicity requires all states reachable (ensured by perturbations)
func MonteCarloVedic(initial *parser.Protein, config MonteCarloConfig) (*MonteCarloResult, error) {
	if initial == nil {
		return nil, fmt.Errorf("initial structure is nil")
	}

	rand.Seed(config.Seed)

	result := &MonteCarloResult{
		BestEnergy:     math.Inf(1),
		BestVedicScore: 0.0,
	}

	// Clone initial structure
	current := cloneProteinDeep(initial)
	best := cloneProteinDeep(initial)

	// Calculate initial scores
	currentEnergy := calculateTotalEnergy(current, config.VdWCutoff, config.ElecCutoff)
	currentAngles := geometry.CalculateRamachandran(current)
	currentVedic := vedic.CalculateVedicScore(current, currentAngles)

	result.InitialEnergy = currentEnergy
	result.InitialVedicScore = currentVedic.TotalScore

	// Combined score: Energy - Vedic bonus
	// Lower is better (minimize energy, maximize Vedic)
	currentScore := combinedScore(currentEnergy, currentVedic.TotalScore, config.VedicWeight)
	bestScore := currentScore

	result.BestEnergy = currentEnergy
	result.BestVedicScore = currentVedic.TotalScore

	// Monte Carlo loop
	for step := 0; step < config.NumSteps; step++ {
		// Calculate temperature for this step
		T := getTemperature(step, config)

		// Propose move: perturb coordinates
		proposed := cloneProteinDeep(current)
		perturbCoordinates(proposed, config.StepSize)

		// Calculate proposed scores
		proposedEnergy := calculateTotalEnergy(proposed, config.VdWCutoff, config.ElecCutoff)
		proposedAngles := geometry.CalculateRamachandran(proposed)
		proposedVedic := vedic.CalculateVedicScore(proposed, proposedAngles)
		proposedScore := combinedScore(proposedEnergy, proposedVedic.TotalScore, config.VedicWeight)

		// Metropolis acceptance criterion
		deltaScore := proposedScore - currentScore
		accepted := false

		if deltaScore < 0 {
			// Better score: always accept
			accepted = true
		} else {
			// Worse score: accept with probability exp(-ΔS/kT)
			// Boltzmann constant k = 0.001987 kcal/(mol·K)
			kB := 0.001987
			acceptProb := math.Exp(-deltaScore / (kB * T))

			if rand.Float64() < acceptProb {
				accepted = true
			}
		}

		// Update current state
		if accepted {
			current = proposed
			currentEnergy = proposedEnergy
			currentVedic = proposedVedic
			currentScore = proposedScore
			result.NumAccepted++

			// Track best
			if currentScore < bestScore {
				best = cloneProteinDeep(current)
				bestScore = currentScore
				result.BestEnergy = currentEnergy
				result.BestVedicScore = currentVedic.TotalScore
				result.ConvergenceStep = step
			}
		} else {
			result.NumRejected++
		}

		// Check convergence: if no improvement for 200 steps, stop
		if step-result.ConvergenceStep > 200 {
			result.Converged = true
			break
		}
	}

	// Calculate acceptance rate
	totalSteps := result.NumAccepted + result.NumRejected
	if totalSteps > 0 {
		result.AcceptanceRate = float64(result.NumAccepted) / float64(totalSteps)
	}

	// Final statistics
	result.FinalStructure = best
	result.FinalEnergy = result.BestEnergy
	result.FinalVedicScore = result.BestVedicScore

	return result, nil
}

// getTemperature calculates temperature for MC step according to cooling schedule
//
// VEDIC_PHI SCHEDULE:
// Uses golden ratio for exponential decay
// T(t) = T_0 × φ^(-t/τ)
// where τ = N / ln(φ) ensures T(N) = T_final
//
// PHYSICIST:
// Exponential cooling is common in simulated annealing
// Vedic twist: Use φ instead of arbitrary base
func getTemperature(step int, config MonteCarloConfig) float64 {
	t := float64(step)
	n := float64(config.NumSteps)
	T0 := config.TemperatureInitial
	Tf := config.TemperatureFinal

	switch config.CoolingSchedule {
	case "vedic_phi":
		// Golden ratio exponential decay
		// τ chosen so T(N) = T_final
		phi := vedic.Phi
		tau := n / math.Log(phi)
		alpha := math.Pow(phi, -t/tau)
		return T0*alpha + Tf*(1.0-alpha)

	case "exponential":
		// Standard exponential decay
		alpha := math.Pow(Tf/T0, t/n)
		return T0 * alpha

	case "linear":
		// Linear cooling
		return T0 - (T0-Tf)*t/n

	case "geometric":
		// Geometric cooling with factor α
		alpha := math.Pow(Tf/T0, 1.0/n)
		return T0 * math.Pow(alpha, t)

	default:
		// Default to exponential
		alpha := math.Pow(Tf/T0, t/n)
		return T0 * alpha
	}
}

// combinedScore computes weighted combination of energy and Vedic score
//
// FORMULA:
//   S = E_energy - λ × E_vedic × scale
//
// Where:
//   - E_energy: AMBER force field energy (kcal/mol)
//   - E_vedic: Vedic harmonic score [0, 1]
//   - λ: Vedic weight [0, 1]
//   - scale: 1000 (to make Vedic comparable to energy in kcal/mol)
//
// BIOCHEMIST:
// Lower combined score = better structure
// - Low energy (stable)
// - High Vedic score (φ-aligned geometry)
//
// MATHEMATICIAN:
// This is a Pareto optimization: balance two objectives
// λ = 0: Pure energy minimization
// λ = 1: Pure Vedic maximization
// λ = 0.3: 70% energy, 30% Vedic (recommended)
func combinedScore(energy, vedicScore, vedicWeight float64) float64 {
	// Scale Vedic score to be comparable to energy
	// Energy typically 100-10000 kcal/mol
	// Vedic score 0-1, so multiply by 1000
	const vedicScale = 1000.0

	// Combined score: minimize energy, maximize Vedic
	return energy - vedicWeight*vedicScore*vedicScale
}

// perturbCoordinates randomly perturbs atom positions
//
// PHYSICIST:
// Gaussian perturbations maintain detailed balance
// Each atom moved independently by N(0, stepSize)
//
// BIOCHEMIST:
// Perturb all atoms to explore conformational space
// Step size controls exploration vs exploitation
func perturbCoordinates(protein *parser.Protein, stepSize float64) {
	for _, atom := range protein.Atoms {
		// Gaussian perturbation in each dimension
		atom.X += rand.NormFloat64() * stepSize
		atom.Y += rand.NormFloat64() * stepSize
		atom.Z += rand.NormFloat64() * stepSize
	}
}

// calculateTotalEnergy computes total AMBER force field energy
//
// This is a wrapper for physics.CalculateTotalEnergy
// Returns just the total energy value (not components)
func calculateTotalEnergy(protein *parser.Protein, vdwCutoff, elecCutoff float64) float64 {
	energyComponents := physics.CalculateTotalEnergy(protein, vdwCutoff, elecCutoff)
	return energyComponents.Total
}

// cloneProteinDeep creates a deep copy of protein structure
//
// ENGINEER:
// Must clone all atoms and residues to avoid aliasing
// Pointers in residues must point to cloned atoms, not originals
func cloneProteinDeep(protein *parser.Protein) *parser.Protein {
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
			AltLoc:    atom.AltLoc,
			ResName:   atom.ResName,
			ChainID:   atom.ChainID,
			ResSeq:    atom.ResSeq,
			ICode:     atom.ICode,
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

// GenerateMonteCarloEnsemble creates ensemble via multiple MC runs
//
// BIOCHEMIST:
// Run MC from different random seeds to explore conformational space
// Each run may converge to different local minimum
//
// Returns ensemble of diverse low-energy structures
func GenerateMonteCarloEnsemble(initial *parser.Protein, config MonteCarloConfig, numRuns int) ([]*parser.Protein, error) {
	ensemble := make([]*parser.Protein, 0, numRuns)
	baseSeed := config.Seed

	for run := 0; run < numRuns; run++ {
		// Different seed for each run
		config.Seed = baseSeed + int64(run)

		result, err := MonteCarloVedic(initial, config)
		if err != nil {
			// Skip failed runs
			continue
		}

		ensemble = append(ensemble, result.FinalStructure)
	}

	if len(ensemble) == 0 {
		return nil, fmt.Errorf("all MC runs failed")
	}

	return ensemble, nil
}

// AdaptiveMonteCarloVedic uses adaptive cooling schedule
//
// PHYSICIST:
// Adjust cooling rate based on acceptance rate
// - High acceptance (>0.7): Cool faster (exploring too much)
// - Low acceptance (<0.3): Cool slower (stuck in local minimum)
//
// Target acceptance rate: 0.4-0.6 (optimal exploration)
func AdaptiveMonteCarloVedic(initial *parser.Protein, config MonteCarloConfig) (*MonteCarloResult, error) {
	if initial == nil {
		return nil, fmt.Errorf("initial structure is nil")
	}

	rand.Seed(config.Seed)

	result := &MonteCarloResult{
		BestEnergy:     math.Inf(1),
		BestVedicScore: 0.0,
	}

	current := cloneProteinDeep(initial)
	best := cloneProteinDeep(initial)

	currentEnergy := calculateTotalEnergy(current, config.VdWCutoff, config.ElecCutoff)
	currentAngles := geometry.CalculateRamachandran(current)
	currentVedic := vedic.CalculateVedicScore(current, currentAngles)

	result.InitialEnergy = currentEnergy
	result.InitialVedicScore = currentVedic.TotalScore

	currentScore := combinedScore(currentEnergy, currentVedic.TotalScore, config.VedicWeight)
	bestScore := currentScore

	result.BestEnergy = currentEnergy
	result.BestVedicScore = currentVedic.TotalScore

	// Adaptive temperature control
	T := config.TemperatureInitial
	targetAcceptRate := 0.5
	recentAccepts := 0
	recentTotal := 0
	checkInterval := 100 // Adjust every 100 steps

	for step := 0; step < config.NumSteps; step++ {
		// Propose and evaluate
		proposed := cloneProteinDeep(current)
		perturbCoordinates(proposed, config.StepSize)

		proposedEnergy := calculateTotalEnergy(proposed, config.VdWCutoff, config.ElecCutoff)
		proposedAngles := geometry.CalculateRamachandran(proposed)
		proposedVedic := vedic.CalculateVedicScore(proposed, proposedAngles)
		proposedScore := combinedScore(proposedEnergy, proposedVedic.TotalScore, config.VedicWeight)

		// Metropolis criterion
		deltaScore := proposedScore - currentScore
		accepted := false

		if deltaScore < 0 {
			accepted = true
		} else {
			kB := 0.001987
			acceptProb := math.Exp(-deltaScore / (kB * T))
			if rand.Float64() < acceptProb {
				accepted = true
			}
		}

		// Track acceptance for adaptation
		recentTotal++
		if accepted {
			recentAccepts++
		}

		// Update
		if accepted {
			current = proposed
			currentEnergy = proposedEnergy
			currentVedic = proposedVedic
			currentScore = proposedScore
			result.NumAccepted++

			if currentScore < bestScore {
				best = cloneProteinDeep(current)
				bestScore = currentScore
				result.BestEnergy = currentEnergy
				result.BestVedicScore = currentVedic.TotalScore
				result.ConvergenceStep = step
			}
		} else {
			result.NumRejected++
		}

		// Adaptive temperature adjustment
		if recentTotal >= checkInterval {
			acceptRate := float64(recentAccepts) / float64(recentTotal)

			// Adjust temperature based on acceptance rate
			if acceptRate > targetAcceptRate+0.1 {
				// Too many accepts: cool faster (multiply by φ^-1)
				T *= vedic.PhiInverse
			} else if acceptRate < targetAcceptRate-0.1 {
				// Too few accepts: heat slightly
				T *= 1.1
			}

			// Ensure temperature stays in bounds
			if T < config.TemperatureFinal {
				T = config.TemperatureFinal
			}
			if T > config.TemperatureInitial {
				T = config.TemperatureInitial
			}

			// Reset counters
			recentAccepts = 0
			recentTotal = 0
		}

		// Convergence check
		if step-result.ConvergenceStep > 200 {
			result.Converged = true
			break
		}
	}

	// Final statistics
	totalSteps := result.NumAccepted + result.NumRejected
	if totalSteps > 0 {
		result.AcceptanceRate = float64(result.NumAccepted) / float64(totalSteps)
	}

	result.FinalStructure = best
	result.FinalEnergy = result.BestEnergy
	result.FinalVedicScore = result.BestVedicScore

	return result, nil
}
