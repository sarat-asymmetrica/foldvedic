// Package optimization - Simulated Annealing
//
// WAVE 8.2: Simulated Annealing Optimizer
// Global optimization via thermodynamic cooling with L-BFGS refinement
//
// PHYSICIST: Mimics thermal annealing in metallurgy - slow cooling finds global minimum
// MATHEMATICIAN: Escapes local minima via probabilistic acceptance at high temperature
// BIOCHEMIST: Explores conformational space like protein folding in vivo
// ETHICIST: Well-established physical analogy, interpretable, reproducible
//
// INNOVATION: Vedic golden ratio cooling + L-BFGS hybrid
// Standard SA: Exponential cooling + steepest descent
// Vedic SA: φ-ratio cooling + L-BFGS refinement at low temperature
//
// CITATION:
// Kirkpatrick, S., et al. (1983). "Optimization by simulated annealing."
// Science, 220(4598), 671-680.
//
// Metropolis, N., et al. (1953). "Equation of state calculations by fast computing machines."
// J. Chem. Phys. 21(6): 1087-1092.
package optimization

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// SimulatedAnnealingConfig holds SA optimization parameters
type SimulatedAnnealingConfig struct {
	// Temperature parameters (Kelvin)
	TemperatureInitial float64
	TemperatureFinal   float64

	// Number of SA steps
	NumSteps int

	// Cooling schedule: "exponential", "linear", "geometric", "vedic_phi"
	CoolingSchedule string

	// Perturbation size (Angstroms)
	// At high T: large perturbations, low T: small perturbations
	PerturbationInitial float64
	PerturbationFinal   float64

	// Refinement: Use L-BFGS at low temperature
	UseLBFGSRefinement bool
	RefinementThreshold float64 // Temperature below which to use L-BFGS
	LBFGSSteps         int      // Number of L-BFGS iterations per refinement

	// Energy calculation cutoffs
	VdWCutoff  float64
	ElecCutoff float64

	// Random seed
	Seed int64

	// Verbose logging
	Verbose bool
}

// DefaultSimulatedAnnealingConfig returns recommended SA parameters
func DefaultSimulatedAnnealingConfig() SimulatedAnnealingConfig {
	return SimulatedAnnealingConfig{
		TemperatureInitial:  1000.0,          // 1000 K (high exploration)
		TemperatureFinal:    1.0,             // 1 K (low exploitation)
		NumSteps:            5000,            // 5000 SA steps
		CoolingSchedule:     "vedic_phi",     // Golden ratio cooling
		PerturbationInitial: 2.0,             // 2.0 Å at high T
		PerturbationFinal:   0.1,             // 0.1 Å at low T
		UseLBFGSRefinement:  true,            // Hybrid SA+LBFGS
		RefinementThreshold: 50.0,            // Refine below 50 K
		LBFGSSteps:          50,              // 50 L-BFGS iterations
		VdWCutoff:           10.0,
		ElecCutoff:          12.0,
		Seed:                42,
		Verbose:             false,
	}
}

// SimulatedAnnealingResult holds SA optimization results
type SimulatedAnnealingResult struct {
	// Optimization statistics
	Steps             int
	AcceptedSteps     int
	RejectedSteps     int
	AcceptanceRate    float64

	// Energy statistics
	InitialEnergy     float64
	FinalEnergy       float64
	BestEnergy        float64
	EnergyChange      float64

	// Convergence
	Converged         bool
	Reason            string

	// Performance
	FunctionEvaluations int
	LBFGSRefinements    int
}

// SimulatedAnnealing performs simulated annealing optimization
//
// ALGORITHM:
// 1. Initialize: T = T_initial, x = x_0
// 2. For each step:
//    a. Perturb structure: x' = x + δ (random perturbation)
//    b. Calculate energy change: ΔE = E(x') - E(x)
//    c. Accept with probability:
//       - If ΔE < 0: always accept (better)
//       - If ΔE > 0: accept with P = exp(-ΔE/kT) (Metropolis criterion)
//    d. Cool temperature: T → T × α
//    e. If T < T_refine: run L-BFGS refinement
// 3. Return best structure found
//
// PHYSICIST:
// This mimics physical annealing:
// - High T: atoms jump freely (exploration)
// - Low T: atoms settle into minima (exploitation)
// - Slow cooling: finds global minimum (not just local)
//
// MATHEMATICIAN:
// Convergence to global minimum guaranteed if:
// - Cooling is infinitely slow: T(t) > C / log(1+t)
// - In practice: Use finite cooling for efficiency
func SimulatedAnnealing(protein *parser.Protein, config SimulatedAnnealingConfig) (*SimulatedAnnealingResult, error) {
	if protein == nil {
		return nil, fmt.Errorf("protein is nil")
	}

	rand.Seed(config.Seed)

	result := &SimulatedAnnealingResult{}

	// Calculate initial energy
	currentEnergy := evaluateEnergy(protein, LBFGSConfig{VdWCutoff: config.VdWCutoff, ElecCutoff: config.ElecCutoff})
	result.InitialEnergy = currentEnergy
	result.BestEnergy = currentEnergy
	result.FunctionEvaluations = 1

	// Clone best structure
	bestProtein := cloneProtein(protein)

	if config.Verbose {
		fmt.Printf("Simulated Annealing: Initial energy = %.2f kcal/mol\n", currentEnergy)
	}

	// Boltzmann constant: k_B = 0.001987 kcal/(mol·K)
	const kB = 0.001987

	lastRefinement := 0 // Track when we last did L-BFGS refinement

	// Simulated annealing loop
	for step := 0; step < config.NumSteps; step++ {
		result.Steps = step + 1

		// Calculate temperature for this step
		T := getTemperatureSchedule(step, config)

		// Calculate perturbation size (decreases with temperature)
		perturbSize := getPerturbationSize(step, config)

		// Perturb structure
		proposedProtein := cloneProtein(protein)
		perturbStructure(proposedProtein, perturbSize)

		// Calculate proposed energy
		proposedEnergy := evaluateEnergy(proposedProtein, LBFGSConfig{VdWCutoff: config.VdWCutoff, ElecCutoff: config.ElecCutoff})
		result.FunctionEvaluations++

		// Metropolis acceptance criterion
		deltaE := proposedEnergy - currentEnergy
		accepted := false

		if deltaE < 0 {
			// Better energy: always accept
			accepted = true
		} else {
			// Worse energy: accept with probability exp(-ΔE/kT)
			acceptProb := math.Exp(-deltaE / (kB * T))
			if rand.Float64() < acceptProb {
				accepted = true
			}
		}

		// Update statistics
		if accepted {
			result.AcceptedSteps++
			protein = proposedProtein
			currentEnergy = proposedEnergy

			// Track best
			if currentEnergy < result.BestEnergy {
				result.BestEnergy = currentEnergy
				bestProtein = cloneProtein(protein)
			}
		} else {
			result.RejectedSteps++
		}

		// L-BFGS refinement at low temperature
		if config.UseLBFGSRefinement && T < config.RefinementThreshold {
			// Refine every 100 steps at low temperature
			if step-lastRefinement >= 100 {
				if config.Verbose {
					fmt.Printf("  Step %d (T=%.1f K): Refining with L-BFGS...\n", step, T)
				}

				lbfgsConfig := DefaultLBFGSConfig()
				lbfgsConfig.MaxIterations = config.LBFGSSteps
				lbfgsConfig.Verbose = false

				lbfgsResult, err := MinimizeLBFGS(protein, lbfgsConfig)
				if err == nil && lbfgsResult.Converged {
					currentEnergy = lbfgsResult.FinalEnergy
					result.FunctionEvaluations += lbfgsResult.FunctionEvaluations
					result.LBFGSRefinements++

					// Update best if improved
					if currentEnergy < result.BestEnergy {
						result.BestEnergy = currentEnergy
						bestProtein = cloneProtein(protein)
					}

					if config.Verbose {
						fmt.Printf("    L-BFGS: E = %.2f kcal/mol (%d iters)\n", currentEnergy, lbfgsResult.Iterations)
					}
				}

				lastRefinement = step
			}
		}

		// Progress logging
		if config.Verbose && (step%500 == 0 || step < 10) {
			acceptRate := float64(result.AcceptedSteps) / float64(result.AcceptedSteps+result.RejectedSteps)
			fmt.Printf("  Step %4d: T=%6.1f K, E=%8.2f, Best=%8.2f, Accept=%.2f, δ=%.3f Å\n",
				step, T, currentEnergy, result.BestEnergy, acceptRate, perturbSize)
		}

		// Early stopping: if temperature is very low and no improvement for 500 steps
		if T < config.TemperatureFinal*2.0 && step-lastRefinement > 500 {
			// Check if best energy hasn't improved
			stagnant := math.Abs(currentEnergy-result.BestEnergy) < 0.1

			if stagnant {
				result.Converged = true
				result.Reason = fmt.Sprintf("Converged at step %d (T=%.2f K, no improvement)", step, T)
				break
			}
		}
	}

	// Final statistics
	result.FinalEnergy = result.BestEnergy
	result.EnergyChange = result.InitialEnergy - result.BestEnergy

	totalSteps := result.AcceptedSteps + result.RejectedSteps
	if totalSteps > 0 {
		result.AcceptanceRate = float64(result.AcceptedSteps) / float64(totalSteps)
	}

	if !result.Converged {
		result.Reason = fmt.Sprintf("Completed %d SA steps", config.NumSteps)
	}

	// Apply best structure
	copyProteinCoordinates(bestProtein, protein)

	if config.Verbose {
		fmt.Printf("\nSimulated Annealing Complete:\n")
		fmt.Printf("  Steps: %d, Accepted: %d (%.1f%%)\n", result.Steps, result.AcceptedSteps, result.AcceptanceRate*100)
		fmt.Printf("  Energy: %.2f → %.2f (Δ = %.2f kcal/mol)\n",
			result.InitialEnergy, result.FinalEnergy, result.EnergyChange)
		fmt.Printf("  L-BFGS refinements: %d\n", result.LBFGSRefinements)
	}

	return result, nil
}

// getTemperatureSchedule calculates temperature for SA step
//
// VEDIC_PHI SCHEDULE:
// Uses golden ratio for cooling
// T(t) = T_0 × φ^(-t/τ) + T_f × (1 - φ^(-t/τ))
// where τ = N / ln(φ) ensures smooth transition
func getTemperatureSchedule(step int, config SimulatedAnnealingConfig) float64 {
	t := float64(step)
	n := float64(config.NumSteps)
	T0 := config.TemperatureInitial
	Tf := config.TemperatureFinal

	switch config.CoolingSchedule {
	case "vedic_phi":
		// Golden ratio exponential cooling
		phi := 1.618033988749895
		tau := n / math.Log(phi)
		alpha := math.Pow(phi, -t/tau)
		return T0*alpha + Tf*(1.0-alpha)

	case "exponential":
		// Standard exponential: T = T_0 × α^t
		alpha := math.Pow(Tf/T0, 1.0/n)
		return T0 * math.Pow(alpha, t)

	case "linear":
		// Linear cooling: T = T_0 - (T_0 - T_f) × t/N
		return T0 - (T0-Tf)*t/n

	case "geometric":
		// Geometric: T = T_0 / (1 + α×t)
		alpha := (T0 - Tf) / (Tf * n)
		return T0 / (1.0 + alpha*t)

	default:
		// Default to exponential
		alpha := math.Pow(Tf/T0, 1.0/n)
		return T0 * math.Pow(alpha, t)
	}
}

// getPerturbationSize calculates perturbation size for SA step
//
// PHYSICIST:
// At high T: large perturbations (exploration)
// At low T: small perturbations (refinement)
// Decreases with temperature
func getPerturbationSize(step int, config SimulatedAnnealingConfig) float64 {
	t := float64(step)
	n := float64(config.NumSteps)
	d0 := config.PerturbationInitial
	df := config.PerturbationFinal

	// Linear decrease
	return d0 - (d0-df)*t/n
}

// perturbStructure randomly perturbs protein coordinates
func perturbStructure(protein *parser.Protein, perturbSize float64) {
	for _, atom := range protein.Atoms {
		atom.X += (rand.Float64()*2.0 - 1.0) * perturbSize
		atom.Y += (rand.Float64()*2.0 - 1.0) * perturbSize
		atom.Z += (rand.Float64()*2.0 - 1.0) * perturbSize
	}
}

// cloneProtein creates deep copy of protein
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

	// Clone residues
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

// copyProteinCoordinates copies coordinates from source to target
func copyProteinCoordinates(source, target *parser.Protein) {
	for i := range source.Atoms {
		if i < len(target.Atoms) {
			target.Atoms[i].X = source.Atoms[i].X
			target.Atoms[i].Y = source.Atoms[i].Y
			target.Atoms[i].Z = source.Atoms[i].Z
		}
	}
}

// HybridOptimization combines SA global search with L-BFGS local refinement
//
// ALGORITHM:
// 1. Simulated Annealing: Find good basin
// 2. L-BFGS: Refine to local minimum
//
// PHYSICIST:
// Best of both worlds:
// - SA escapes local minima
// - L-BFGS converges quickly to minimum
func HybridOptimization(protein *parser.Protein, saConfig SimulatedAnnealingConfig, lbfgsConfig LBFGSConfig) (*SimulatedAnnealingResult, error) {
	// Phase 1: Simulated Annealing
	saResult, err := SimulatedAnnealing(protein, saConfig)
	if err != nil {
		return nil, fmt.Errorf("SA phase failed: %w", err)
	}

	// Phase 2: L-BFGS refinement
	lbfgsResult, err := MinimizeLBFGS(protein, lbfgsConfig)
	if err != nil {
		// SA succeeded but L-BFGS failed: still return SA result
		return saResult, nil
	}

	// Update SA result with L-BFGS improvements
	saResult.FinalEnergy = lbfgsResult.FinalEnergy
	saResult.EnergyChange = saResult.InitialEnergy - lbfgsResult.FinalEnergy
	saResult.FunctionEvaluations += lbfgsResult.FunctionEvaluations

	if lbfgsConfig.Verbose {
		fmt.Printf("Hybrid Optimization: SA+LBFGS complete\n")
		fmt.Printf("  SA: %.2f → %.2f kcal/mol\n", saResult.InitialEnergy, saResult.BestEnergy)
		fmt.Printf("  LBFGS: %.2f → %.2f kcal/mol\n", saResult.BestEnergy, lbfgsResult.FinalEnergy)
		fmt.Printf("  Total improvement: %.2f kcal/mol\n", saResult.EnergyChange)
	}

	return saResult, nil
}
