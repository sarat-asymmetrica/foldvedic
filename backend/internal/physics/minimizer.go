// Package physics implements energy minimization for protein structures.
//
// MATHEMATICIAN: Steepest descent minimization with line search
// PHYSICIST: Minimize total potential energy to find stable conformation
// BIOCHEMIST: Validates converged structures against known proteins
package physics

import (
	"fmt"
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// MinimizerConfig holds parameters for energy minimization
type MinimizerConfig struct {
	// Maximum number of minimization steps
	MaxSteps int

	// Energy convergence threshold (kcal/mol)
	// Stop if |E_new - E_old| < EnergyTolerance
	EnergyTolerance float64

	// Force convergence threshold (kcal/(mol·Å))
	// Stop if max(|F_i|) < ForceTolerance
	ForceTolerance float64

	// Step size for steepest descent (Å)
	StepSize float64

	// Van der Waals cutoff (Å)
	VdWCutoff float64

	// Electrostatic cutoff (Å)
	ElecCutoff float64
}

// DefaultMinimizerConfig returns default minimization parameters
func DefaultMinimizerConfig() MinimizerConfig {
	return MinimizerConfig{
		MaxSteps:        1000,
		EnergyTolerance: 0.01,     // 0.01 kcal/mol
		ForceTolerance:  0.1,      // 0.1 kcal/(mol·Å)
		StepSize:        0.0001,   // 0.0001 Å (very small for stability with stiff bonds)
		VdWCutoff:       10.0,     // 10 Å
		ElecCutoff:      12.0,     // 12 Å
	}
}

// MinimizationResult holds results of energy minimization
type MinimizationResult struct {
	// Number of steps taken
	Steps int

	// Final energy components
	FinalEnergy EnergyComponents

	// Initial energy
	InitialEnergy float64

	// Energy change
	DeltaEnergy float64

	// Converged flag
	Converged bool

	// Convergence reason
	Reason string
}

// MinimizeEnergy performs steepest descent energy minimization
//
// MATHEMATICIAN:
// Steepest descent algorithm:
// 1. Calculate forces F = -∇E for all atoms
// 2. Update positions: x_new = x_old + α × F (α = step size)
// 3. Repeat until convergence
//
// PHYSICIST:
// This is a simple but robust method for initial minimization
// More sophisticated methods (conjugate gradient, L-BFGS) in Wave 3
//
// Citation: Press, W. H., et al. (2007). "Numerical Recipes."
// Cambridge University Press. Chapter 10.
func MinimizeEnergy(protein *parser.Protein, config MinimizerConfig) (*MinimizationResult, error) {
	result := &MinimizationResult{}

	// Calculate initial energy
	initialEnergy := CalculateTotalEnergy(protein, config.VdWCutoff, config.ElecCutoff)
	result.InitialEnergy = initialEnergy.Total
	prevEnergy := initialEnergy.Total

	// Minimization loop
	for step := 0; step < config.MaxSteps; step++ {
		result.Steps = step + 1

		// Calculate forces on all atoms
		forces := CalculateForces(protein, config.VdWCutoff, config.ElecCutoff)

		// Update atom positions using steepest descent
		// x_new = x_old + step_size × F
		maxForce := 0.0
		for _, atom := range protein.Atoms {
			force, exists := forces[atom.Serial]
			if !exists {
				continue
			}

			// Track maximum force magnitude for convergence check
			forceMag := force.Magnitude()
			if forceMag > maxForce {
				maxForce = forceMag
			}

			// Update position
			atom.X += config.StepSize * force.X
			atom.Y += config.StepSize * force.Y
			atom.Z += config.StepSize * force.Z
		}

		// Calculate new energy
		currentEnergy := CalculateTotalEnergy(protein, config.VdWCutoff, config.ElecCutoff)

		// SAFETY: Check for numerical explosion
		// If energy increased drastically, step size is too large - stop minimization
		if currentEnergy.Total > prevEnergy*10.0 || math.IsNaN(currentEnergy.Total) || math.IsInf(currentEnergy.Total, 0) {
			result.FinalEnergy = currentEnergy
			result.DeltaEnergy = result.InitialEnergy - currentEnergy.Total
			result.Converged = false
			result.Reason = "Numerical instability detected (step size too large)"
			return result, fmt.Errorf("energy minimization unstable: energy increased from %.2f to %.2e", prevEnergy, currentEnergy.Total)
		}

		// Check for convergence
		deltaE := math.Abs(currentEnergy.Total - prevEnergy)

		// Energy convergence
		if deltaE < config.EnergyTolerance {
			result.FinalEnergy = currentEnergy
			result.DeltaEnergy = result.InitialEnergy - currentEnergy.Total
			result.Converged = true
			result.Reason = fmt.Sprintf("Energy converged (ΔE = %.6f < %.6f kcal/mol)",
				deltaE, config.EnergyTolerance)
			return result, nil
		}

		// Force convergence
		if maxForce < config.ForceTolerance {
			result.FinalEnergy = currentEnergy
			result.DeltaEnergy = result.InitialEnergy - currentEnergy.Total
			result.Converged = true
			result.Reason = fmt.Sprintf("Forces converged (max F = %.6f < %.6f kcal/(mol·Å))",
				maxForce, config.ForceTolerance)
			return result, nil
		}

		// Update previous energy
		prevEnergy = currentEnergy.Total

		// Log progress every 100 steps
		if step%100 == 0 {
			// Can add logging here if needed
		}
	}

	// Max steps reached without convergence
	finalEnergy := CalculateTotalEnergy(protein, config.VdWCutoff, config.ElecCutoff)
	result.FinalEnergy = finalEnergy
	result.DeltaEnergy = result.InitialEnergy - finalEnergy.Total
	result.Converged = false
	result.Reason = fmt.Sprintf("Max steps reached (%d)", config.MaxSteps)

	return result, nil
}

// ValidateEnergy checks if energy values are physically reasonable
//
// BIOCHEMIST:
// Typical protein energies (per residue):
// - Bond: ~100-200 kcal/mol/residue
// - Angle: ~50-100 kcal/mol/residue
// - VdW: ~-10 to -50 kcal/mol/residue
// - Elec: ~-100 to -300 kcal/mol/residue (varies greatly)
//
// Total: Usually negative (attractive forces dominate)
func ValidateEnergy(energy EnergyComponents, numResidues int) []string {
	warnings := make([]string, 0)

	if numResidues == 0 {
		return warnings
	}

	// Energy per residue
	bondPerRes := energy.Bond / float64(numResidues)
	anglePerRes := energy.Angle / float64(numResidues)
	vdwPerRes := energy.VanDerWaals / float64(numResidues)
	_ = energy.Electrostatic / float64(numResidues) // elecPerRes - varies too much to validate

	// Check for unphysical values
	if bondPerRes > 500 {
		warnings = append(warnings, fmt.Sprintf("Bond energy very high: %.1f kcal/mol/res (expect 100-200)", bondPerRes))
	}

	if anglePerRes > 300 {
		warnings = append(warnings, fmt.Sprintf("Angle energy very high: %.1f kcal/mol/res (expect 50-100)", anglePerRes))
	}

	if vdwPerRes > 100 {
		warnings = append(warnings, "VdW energy positive and large (steric clashes?)")
	}

	if math.IsNaN(energy.Total) || math.IsInf(energy.Total, 0) {
		warnings = append(warnings, "Total energy is NaN or Inf (numerical instability)")
	}

	return warnings
}
