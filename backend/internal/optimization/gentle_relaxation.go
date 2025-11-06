// Package optimization - Gentle Energy Relaxation
//
// WAVE 11.2: Quick & Dirty Energy Minimization
//
// WRIGHT BROTHERS EMPIRICISM:
// - L-BFGS is too aggressive → numerical explosion
// - SA is too slow → takes minutes
// - Solution: GENTLE steepest descent with tiny steps
//
// GOAL: Remove severe clashes, don't find global minimum
// - Just enough to get realistic energies
// - Fast (seconds not minutes)
// - Stable (no explosions)
//
// CROSS-DOMAIN: Chemistry software often uses "gentle equilibration"
// before full MD simulations to remove bad contacts
package optimization

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/physics"
)

// GentleRelaxationConfig holds gentle minimization parameters
type GentleRelaxationConfig struct {
	MaxSteps        int     // Number of steepest descent steps
	StepSize        float64 // Å per step (SMALL!)
	EnergyTolerance float64 // Stop if energy change < this
	VdWCutoff       float64 // Van der Waals cutoff
	ElecCutoff      float64 // Electrostatic cutoff
}

// DefaultGentleRelaxationConfig returns safe parameters
func DefaultGentleRelaxationConfig() GentleRelaxationConfig {
	return GentleRelaxationConfig{
		MaxSteps:        50,    // Quick relaxation
		StepSize:        0.01,  // Tiny steps (0.01 Å)
		EnergyTolerance: 0.1,   // Stop if energy changes < 0.1 kcal/mol
		VdWCutoff:       10.0,
		ElecCutoff:      12.0,
	}
}

// GentleRelaxationResult holds relaxation results
type GentleRelaxationResult struct {
	InitialEnergy float64
	FinalEnergy   float64
	EnergyChange  float64
	Steps         int
	Converged     bool
}

// GentleRelax performs gentle energy minimization
//
// ALGORITHM: Steepest descent with tiny steps
// 1. Calculate forces (negative gradient)
// 2. Move atoms slightly in direction of forces
// 3. Repeat until energy stops decreasing
//
// WHY THIS WORKS:
// - Small steps → numerically stable
// - Few iterations → fast
// - No second-order methods → simple
//
// TRADEOFFS:
// - Won't find global minimum (don't care)
// - Won't fully optimize (don't care)
// - WILL remove severe clashes (what we need!)
func GentleRelax(protein *parser.Protein, config GentleRelaxationConfig) (*GentleRelaxationResult, error) {
	result := &GentleRelaxationResult{}

	// Calculate initial energy
	energyComps := physics.CalculateTotalEnergy(protein, config.VdWCutoff, config.ElecCutoff)
	result.InitialEnergy = energyComps.Total
	prevEnergy := energyComps.Total

	for step := 0; step < config.MaxSteps; step++ {
		// Calculate forces on all atoms
		forces := physics.CalculateForces(protein, config.VdWCutoff, config.ElecCutoff)

		// Move atoms in direction of forces (with TINY steps)
		moved := false
		for _, atom := range protein.Atoms {
			if force, exists := forces[atom.Serial]; exists {
				// Force points toward lower energy
				// Move atom by small fraction of force
				displacement := math.Sqrt(force.X*force.X + force.Y*force.Y + force.Z*force.Z)

				if displacement > 1e-6 {
					// Normalize and scale by step size
					scale := config.StepSize / displacement

					// Cap maximum displacement per step
					if scale > 0.1 {
						scale = 0.1
					}

					atom.X += force.X * scale
					atom.Y += force.Y * scale
					atom.Z += force.Z * scale
					moved = true
				}
			}
		}

		// Recalculate energy
		energyComps = physics.CalculateTotalEnergy(protein, config.VdWCutoff, config.ElecCutoff)
		currentEnergy := energyComps.Total

		// Check convergence
		energyDelta := prevEnergy - currentEnergy

		if math.Abs(energyDelta) < config.EnergyTolerance {
			result.FinalEnergy = currentEnergy
			result.Steps = step + 1
			result.Converged = true
			result.EnergyChange = result.InitialEnergy - result.FinalEnergy
			return result, nil
		}

		// Safety: If energy increases dramatically, stop
		if currentEnergy > prevEnergy*2.0 && step > 5 {
			// Energy exploding, revert or stop
			result.FinalEnergy = prevEnergy
			result.Steps = step
			result.Converged = false
			result.EnergyChange = result.InitialEnergy - result.FinalEnergy
			return result, nil
		}

		prevEnergy = currentEnergy

		if !moved {
			// No atoms moved, converged
			result.FinalEnergy = currentEnergy
			result.Steps = step + 1
			result.Converged = true
			result.EnergyChange = result.InitialEnergy - result.FinalEnergy
			return result, nil
		}
	}

	// Max steps reached
	result.FinalEnergy = prevEnergy
	result.Steps = config.MaxSteps
	result.Converged = false
	result.EnergyChange = result.InitialEnergy - result.FinalEnergy

	return result, nil
}

// QuickClashRemoval removes severe atomic clashes
//
// WILD IDEA: Even simpler than GentleRelax
// Just push apart atoms that are too close
//
// ALGORITHM:
// 1. Find all atom pairs closer than 2 Å (severe clash)
// 2. Push them apart to 2.5 Å
// 3. Done!
//
// This is like emergency surgery before the real optimization
func QuickClashRemoval(protein *parser.Protein) int {
	clashesFixed := 0
	minDist := 2.0   // Å - anything closer is a clash
	targetDist := 2.5 // Å - push apart to this distance

	atoms := protein.Atoms

	for i := 0; i < len(atoms); i++ {
		for j := i + 1; j < len(atoms); j++ {
			a1 := atoms[i]
			a2 := atoms[j]

			// Calculate distance
			dx := a2.X - a1.X
			dy := a2.Y - a1.Y
			dz := a2.Z - a1.Z
			dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

			// Check for clash
			if dist < minDist && dist > 0.1 {
				// Severe clash! Push apart
				// Move each atom half the distance needed
				moveEach := (targetDist - dist) / 2.0

				// Direction vector (normalized)
				nx := dx / dist
				ny := dy / dist
				nz := dz / dist

				// Move atoms apart
				a1.X -= nx * moveEach
				a1.Y -= ny * moveEach
				a1.Z -= nz * moveEach

				a2.X += nx * moveEach
				a2.Y += ny * moveEach
				a2.Z += nz * moveEach

				clashesFixed++
			}
		}
	}

	return clashesFixed
}
