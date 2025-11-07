// Package physics implements energy calculation and minimization for protein structures.
//
// PHYSICIST: Total energy = bonded + non-bonded terms
// MATHEMATICIAN: Gradient-based minimization (L-BFGS)
// BIOCHEMIST: Validates against known protein energies
package physics

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// EnergyComponents holds breakdown of total energy
type EnergyComponents struct {
	Bond          float64 // Bond stretching energy
	Angle         float64 // Angle bending energy
	Dihedral      float64 // Ramachandran dihedral energy (backbone constraints)
	VanDerWaals   float64 // Lennard-Jones energy
	Electrostatic float64 // Coulomb energy
	Total         float64 // Sum of all components
}

// CalculateTotalEnergy computes all energy terms for a protein
//
// PHYSICIST:
// E_total = E_bond + E_angle + E_dihedral + E_vdw + E_elec
//
// Parameters:
// - protein: Protein structure with atomic coordinates
// - vdwCutoff: Van der Waals cutoff distance (typically 8-12 Å)
// - elecCutoff: Electrostatic cutoff distance (typically 12-15 Å)
//
// Returns: Energy components in kcal/mol
func CalculateTotalEnergy(protein *parser.Protein, vdwCutoff, elecCutoff float64) EnergyComponents {
	energy := EnergyComponents{}

	// Bond energy: Sum over all covalent bonds
	energy.Bond = calculateBondEnergyTotal(protein)

	// Angle energy: Sum over all bond angles
	energy.Angle = calculateAngleEnergyTotal(protein)

	// Dihedral energy: Ramachandran potential (backbone φ,ψ constraints)
	energy.Dihedral = RamachandranPotential(protein)

	// Van der Waals: Sum over all non-bonded pairs
	energy.VanDerWaals = calculateVanDerWaalsTotal(protein, vdwCutoff)

	// Electrostatic: Sum over all non-bonded pairs
	energy.Electrostatic = calculateElectrostaticTotal(protein, elecCutoff)

	// Total
	energy.Total = energy.Bond + energy.Angle + energy.Dihedral + energy.VanDerWaals + energy.Electrostatic

	// Cap energy to prevent overflow
	// Realistic protein energies: -500 to +2000 kcal/mol
	// >10,000 indicates severe steric clashes or coordinate corruption
	// <-10,000 indicates unphysical attraction
	if energy.Total > 10000.0 {
		energy.Total = 10000.0
	}
	if energy.Total < -10000.0 {
		energy.Total = -10000.0
	}

	return energy
}

// calculateBondEnergyTotal sums bond energies for all bonds in protein
func calculateBondEnergyTotal(protein *parser.Protein) float64 {
	totalEnergy := 0.0

	// Iterate over residues
	for _, res := range protein.Residues {
		if !res.HasCompleteBackbone() {
			continue
		}

		// N-CA bond
		if res.N != nil && res.CA != nil {
			params := GetBondParams("N", "CA")
			totalEnergy += CalculateBondEnergy(res.N, res.CA, params)
		}

		// CA-C bond
		if res.CA != nil && res.C != nil {
			params := GetBondParams("CA", "C")
			totalEnergy += CalculateBondEnergy(res.CA, res.C, params)
		}

		// C-O bond (carbonyl)
		if res.C != nil && res.O != nil {
			params := GetBondParams("C", "O")
			totalEnergy += CalculateBondEnergy(res.C, res.O, params)
		}
	}

	// Peptide bonds between residues
	for i := 0; i < len(protein.Residues)-1; i++ {
		res1 := protein.Residues[i]
		res2 := protein.Residues[i+1]

		if res1.C != nil && res2.N != nil {
			params := GetBondParams("C", "N")
			totalEnergy += CalculateBondEnergy(res1.C, res2.N, params)
		}
	}

	return totalEnergy
}

// calculateAngleEnergyTotal sums angle energies for all angles in protein
func calculateAngleEnergyTotal(protein *parser.Protein) float64 {
	totalEnergy := 0.0

	// Iterate over residues
	for _, res := range protein.Residues {
		if !res.HasCompleteBackbone() {
			continue
		}

		// N-CA-C angle
		if res.N != nil && res.CA != nil && res.C != nil {
			params := GetAngleParams("N", "CA", "C")
			totalEnergy += CalculateAngleEnergy(res.N, res.CA, res.C, params)
		}

		// CA-C-O angle
		if res.CA != nil && res.C != nil && res.O != nil {
			params := GetAngleParams("CA", "C", "O")
			totalEnergy += CalculateAngleEnergy(res.CA, res.C, res.O, params)
		}
	}

	// Inter-residue angles
	for i := 0; i < len(protein.Residues)-1; i++ {
		res1 := protein.Residues[i]
		res2 := protein.Residues[i+1]

		// CA-C-N angle (across peptide bond)
		if res1.CA != nil && res1.C != nil && res2.N != nil {
			params := GetAngleParams("CA", "C", "N")
			totalEnergy += CalculateAngleEnergy(res1.CA, res1.C, res2.N, params)
		}

		// C-N-CA angle (across peptide bond)
		if res1.C != nil && res2.N != nil && res2.CA != nil {
			params := GetAngleParams("C", "N", "CA")
			totalEnergy += CalculateAngleEnergy(res1.C, res2.N, res2.CA, params)
		}
	}

	return totalEnergy
}

// calculateVanDerWaalsTotal sums Lennard-Jones energies for all non-bonded pairs
//
// PHYSICIST:
// Only calculate for atoms separated by >3 bonds (1-4 and beyond)
// Use cutoff distance to reduce O(n²) cost
func calculateVanDerWaalsTotal(protein *parser.Protein, cutoff float64) float64 {
	totalEnergy := 0.0

	// Simple O(n²) loop for now
	// TODO: Use spatial hashing for O(n) performance (Wave 3 - Williams Optimizer)
	atoms := protein.Atoms

	for i := 0; i < len(atoms); i++ {
		for j := i + 1; j < len(atoms); j++ {
			// Skip if atoms are in same residue or adjacent residues (bonded/1-4 interactions)
			if math.Abs(float64(atoms[i].ResSeq-atoms[j].ResSeq)) <= 1 {
				continue
			}

			energy := CalculateLennardJonesEnergy(atoms[i], atoms[j], cutoff)
			totalEnergy += energy
		}
	}

	return totalEnergy
}

// calculateElectrostaticTotal sums Coulomb energies for all non-bonded pairs
func calculateElectrostaticTotal(protein *parser.Protein, cutoff float64) float64 {
	totalEnergy := 0.0

	// Simplified partial charges (backbone only, from AMBER ff14SB)
	charges := map[string]float64{
		"N":  -0.4157, // Backbone nitrogen
		"CA": 0.0337,  // Alpha carbon
		"C":  0.5973,  // Carbonyl carbon
		"O":  -0.5679, // Carbonyl oxygen
	}

	atoms := protein.Atoms

	for i := 0; i < len(atoms); i++ {
		for j := i + 1; j < len(atoms); j++ {
			// Skip if atoms are in same residue or adjacent residues
			if math.Abs(float64(atoms[i].ResSeq-atoms[j].ResSeq)) <= 1 {
				continue
			}

			// Get charges
			charge1, ok1 := charges[atoms[i].Name]
			charge2, ok2 := charges[atoms[j].Name]

			if !ok1 || !ok2 {
				continue // Skip atoms with unknown charges
			}

			energy := CalculateElectrostaticEnergy(atoms[i], atoms[j], charge1, charge2, cutoff)
			totalEnergy += energy
		}
	}

	return totalEnergy
}

// CalculateForces computes forces on all atoms from all energy terms
//
// MATHEMATICIAN:
// F = -∇E (force is negative gradient of energy)
// Returns force vector for each atom
func CalculateForces(protein *parser.Protein, vdwCutoff, elecCutoff float64) map[int]Vector3 {
	forces := make(map[int]Vector3)

	// Initialize all forces to zero
	for _, atom := range protein.Atoms {
		forces[atom.Serial] = Vector3{X: 0, Y: 0, Z: 0}
	}

	// Bond forces
	addBondForces(protein, forces)

	// TODO: Angle forces, VdW forces, electrostatic forces
	// For Wave 1, we're focusing on basic bond forces
	// Full implementation in Wave 2

	return forces
}

// addBondForces adds bond forces to force map
func addBondForces(protein *parser.Protein, forces map[int]Vector3) {
	// Iterate over residues
	for _, res := range protein.Residues {
		if !res.HasCompleteBackbone() {
			continue
		}

		// N-CA bond
		if res.N != nil && res.CA != nil {
			params := GetBondParams("N", "CA")
			force := CalculateBondForce(res.N, res.CA, params)

			// Newton's third law: equal and opposite forces
			forces[res.N.Serial] = forces[res.N.Serial].Add(force.Mul(-1))
			forces[res.CA.Serial] = forces[res.CA.Serial].Add(force)
		}

		// CA-C bond
		if res.CA != nil && res.C != nil {
			params := GetBondParams("CA", "C")
			force := CalculateBondForce(res.CA, res.C, params)

			forces[res.CA.Serial] = forces[res.CA.Serial].Add(force.Mul(-1))
			forces[res.C.Serial] = forces[res.C.Serial].Add(force)
		}

		// C-O bond
		if res.C != nil && res.O != nil {
			params := GetBondParams("C", "O")
			force := CalculateBondForce(res.C, res.O, params)

			forces[res.C.Serial] = forces[res.C.Serial].Add(force.Mul(-1))
			forces[res.O.Serial] = forces[res.O.Serial].Add(force)
		}
	}

	// Peptide bonds
	for i := 0; i < len(protein.Residues)-1; i++ {
		res1 := protein.Residues[i]
		res2 := protein.Residues[i+1]

		if res1.C != nil && res2.N != nil {
			params := GetBondParams("C", "N")
			force := CalculateBondForce(res1.C, res2.N, params)

			forces[res1.C.Serial] = forces[res1.C.Serial].Add(force.Mul(-1))
			forces[res2.N.Serial] = forces[res2.N.Serial].Add(force)
		}
	}
}
