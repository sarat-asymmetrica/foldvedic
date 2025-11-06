// Package physics implements force field calculations for protein energy minimization.
//
// BIOCHEMIST: AMBER ff14SB force field - validated parameters from quantum chemistry
// PHYSICIST: Standard molecular mechanics energy function with bonded and non-bonded terms
// MATHEMATICIAN: Analytical gradients for L-BFGS optimization
// ETHICIST: All parameters cited from peer-reviewed literature
package physics

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// Vector3 represents a 3D vector for forces and positions
type Vector3 struct {
	X, Y, Z float64
}

// Add adds two vectors
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// Sub subtracts two vectors
func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// Mul multiplies vector by scalar
func (v Vector3) Mul(scalar float64) Vector3 {
	return Vector3{X: v.X * scalar, Y: v.Y * scalar, Z: v.Z * scalar}
}

// Dot computes dot product
func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Magnitude returns vector length
func (v Vector3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns unit vector
func (v Vector3) Normalize() Vector3 {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector3{X: 0, Y: 0, Z: 0}
	}
	return Vector3{X: v.X / mag, Y: v.Y / mag, Z: v.Z / mag}
}

// BondParameters holds force constants and equilibrium values for covalent bonds
//
// Citation: Maier, J. A., et al. (2015). "ff14SB: Improving the accuracy of protein
// side chain and backbone parameters from ff99SB." J. Chem. Theory Comput. 11.8: 3696-3713.
type BondParameters struct {
	// Force constant k_b in kcal/(mol·Å²)
	// E_bond = k_b × (r - r_0)²
	K0 float64

	// Equilibrium bond length r_0 in Angstroms
	R0 float64
}

// Standard AMBER ff14SB bond parameters for peptide backbone
var backboneBondParams = map[string]BondParameters{
	// Peptide bond (C-N)
	"C-N": {K0: 490.0, R0: 1.335},

	// N-CA bond
	"N-CA": {K0: 337.0, R0: 1.449},

	// CA-C bond
	"CA-C": {K0: 317.0, R0: 1.522},

	// C=O bond (carbonyl)
	"C-O": {K0: 570.0, R0: 1.229},

	// Default for unknown bonds
	"default": {K0: 300.0, R0: 1.5},
}

// AngleParameters holds force constants and equilibrium values for bond angles
type AngleParameters struct {
	// Force constant k_θ in kcal/(mol·rad²)
	// E_angle = k_θ × (θ - θ_0)²
	K0 float64

	// Equilibrium angle θ_0 in radians
	Theta0 float64
}

// Standard AMBER ff14SB angle parameters for peptide backbone
var backboneAngleParams = map[string]AngleParameters{
	// N-CA-C angle
	"N-CA-C": {K0: 63.0, Theta0: 110.1 * math.Pi / 180},

	// CA-C-N angle
	"CA-C-N": {K0: 70.0, Theta0: 116.6 * math.Pi / 180},

	// C-N-CA angle
	"C-N-CA": {K0: 50.0, Theta0: 121.9 * math.Pi / 180},

	// CA-C-O angle (carbonyl)
	"CA-C-O": {K0: 80.0, Theta0: 120.4 * math.Pi / 180},

	// Default for unknown angles
	"default": {K0: 50.0, Theta0: math.Pi * 109.5 / 180}, // Tetrahedral
}

// CalculateBondEnergy computes harmonic bond energy
//
// PHYSICIST:
// E = k × (r - r_0)²
// where k = force constant, r = current length, r_0 = equilibrium length
//
// Returns energy in kcal/mol
func CalculateBondEnergy(atom1, atom2 *parser.Atom, params BondParameters) float64 {
	// Calculate current bond length
	dx := atom2.X - atom1.X
	dy := atom2.Y - atom1.Y
	dz := atom2.Z - atom1.Z
	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	// Harmonic potential
	dr := r - params.R0
	energy := params.K0 * dr * dr

	return energy
}

// CalculateBondForce computes force on atoms from bond potential
//
// MATHEMATICIAN:
// F = -∇E = -dE/dr × ∇r
// For harmonic bond: F = -2k(r-r_0) × (r_vec/r)
//
// Returns force vector pointing from atom1 to atom2
func CalculateBondForce(atom1, atom2 *parser.Atom, params BondParameters) Vector3 {
	// Vector from atom1 to atom2
	dx := atom2.X - atom1.X
	dy := atom2.Y - atom1.Y
	dz := atom2.Z - atom1.Z
	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	if r == 0 {
		return Vector3{X: 0, Y: 0, Z: 0}
	}

	// Force magnitude: F = -2k(r - r_0)
	dr := r - params.R0
	forceMag := -2.0 * params.K0 * dr

	// Force direction (unit vector from atom1 to atom2)
	direction := Vector3{X: dx / r, Y: dy / r, Z: dz / r}

	return direction.Mul(forceMag)
}

// CalculateAngleEnergy computes harmonic angle energy
//
// PHYSICIST:
// E = k_θ × (θ - θ_0)²
// where θ is the angle formed by three atoms
//
// Returns energy in kcal/mol
func CalculateAngleEnergy(atom1, atom2, atom3 *parser.Atom, params AngleParameters) float64 {
	// Vectors: atom2 is the central atom
	v1 := Vector3{X: atom1.X - atom2.X, Y: atom1.Y - atom2.Y, Z: atom1.Z - atom2.Z}
	v2 := Vector3{X: atom3.X - atom2.X, Y: atom3.Y - atom2.Y, Z: atom3.Z - atom2.Z}

	// Calculate angle using dot product
	dot := v1.Dot(v2)
	mag1 := v1.Magnitude()
	mag2 := v2.Magnitude()

	if mag1 == 0 || mag2 == 0 {
		return 0
	}

	// cos(θ) = (v1 · v2) / (|v1| × |v2|)
	cosTheta := dot / (mag1 * mag2)

	// Clamp to [-1, 1] to avoid NaN from acos
	if cosTheta > 1.0 {
		cosTheta = 1.0
	} else if cosTheta < -1.0 {
		cosTheta = -1.0
	}

	theta := math.Acos(cosTheta)

	// Harmonic potential
	dTheta := theta - params.Theta0
	energy := params.K0 * dTheta * dTheta

	return energy
}

// LennardJonesParams holds van der Waals parameters
//
// Citation: AMBER ff14SB parameters
type LennardJonesParams struct {
	// Well depth ε in kcal/mol
	Epsilon float64

	// Atom radius σ in Angstroms
	Sigma float64
}

// Standard Lennard-Jones parameters for common atoms
var ljParams = map[string]LennardJonesParams{
	"C": {Epsilon: 0.086, Sigma: 1.908}, // Carbon (sp3)
	"N": {Epsilon: 0.170, Sigma: 1.824}, // Nitrogen (amide)
	"O": {Epsilon: 0.210, Sigma: 1.661}, // Oxygen (carbonyl)
	"H": {Epsilon: 0.016, Sigma: 1.487}, // Hydrogen
	"S": {Epsilon: 0.250, Sigma: 2.000}, // Sulfur
}

// CalculateLennardJonesEnergy computes van der Waals energy (Lennard-Jones 12-6 potential)
//
// PHYSICIST:
// E_LJ = 4ε × [(σ/r)¹² - (σ/r)⁶]
// - Repulsive term (r⁻¹²): Pauli exclusion at short range
// - Attractive term (r⁻⁶): London dispersion forces
//
// Citation: Jones, J. E. (1924). "On the determination of molecular fields."
// Proc. R. Soc. Lond. A 106.738: 463-477.
//
// Returns energy in kcal/mol
func CalculateLennardJonesEnergy(atom1, atom2 *parser.Atom, cutoff float64) float64 {
	// Calculate distance
	dx := atom2.X - atom1.X
	dy := atom2.Y - atom1.Y
	dz := atom2.Z - atom1.Z
	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	// Skip if beyond cutoff (typically 8-12 Å)
	if r > cutoff {
		return 0
	}

	// Get LJ parameters (combine using Lorentz-Berthelot rules)
	params1, ok1 := ljParams[atom1.Element]
	params2, ok2 := ljParams[atom2.Element]

	if !ok1 || !ok2 {
		// Default params for unknown elements
		params1 = LennardJonesParams{Epsilon: 0.1, Sigma: 1.8}
		params2 = params1
	}

	// Combining rules:
	// ε_ij = sqrt(ε_i × ε_j)
	// σ_ij = (σ_i + σ_j) / 2
	epsilon := math.Sqrt(params1.Epsilon * params2.Epsilon)
	sigma := (params1.Sigma + params2.Sigma) / 2.0

	// LJ potential
	sigmaOverR := sigma / r
	term6 := math.Pow(sigmaOverR, 6)
	term12 := term6 * term6

	energy := 4.0 * epsilon * (term12 - term6)

	return energy
}

// CalculateElectrostaticEnergy computes Coulomb electrostatic energy
//
// PHYSICIST:
// E_elec = (q_i × q_j) / (4πε_0 × r_ij)
// With implicit solvent: Use distance-dependent dielectric ε(r) = 4r
//
// BIOCHEMIST:
// Partial charges from AMBER ff14SB:
// - Backbone N: -0.4157
// - Backbone CA: 0.0337
// - Backbone C: 0.5973
// - Backbone O: -0.5679
//
// Returns energy in kcal/mol
func CalculateElectrostaticEnergy(atom1, atom2 *parser.Atom, charge1, charge2, cutoff float64) float64 {
	// Calculate distance
	dx := atom2.X - atom1.X
	dy := atom2.Y - atom1.Y
	dz := atom2.Z - atom1.Z
	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	// Skip if beyond cutoff
	if r > cutoff || r == 0 {
		return 0
	}

	// Coulomb's constant in kcal·Å/(mol·e²)
	// k_e = 1/(4πε_0) = 332.06 kcal·Å/(mol·e²)
	kCoulomb := 332.06

	// Distance-dependent dielectric for implicit solvent
	// ε(r) = 4r (mimics water screening)
	dielectric := 4.0 * r

	// Coulomb energy
	energy := (kCoulomb * charge1 * charge2) / dielectric

	return energy
}

// GetBondParams returns force field parameters for a bond
func GetBondParams(atomType1, atomType2 string) BondParameters {
	// Try both orderings
	key1 := atomType1 + "-" + atomType2
	key2 := atomType2 + "-" + atomType1

	if params, ok := backboneBondParams[key1]; ok {
		return params
	}
	if params, ok := backboneBondParams[key2]; ok {
		return params
	}

	// Default
	return backboneBondParams["default"]
}

// GetAngleParams returns force field parameters for an angle
func GetAngleParams(atomType1, atomType2, atomType3 string) AngleParameters {
	// Try both orderings
	key1 := atomType1 + "-" + atomType2 + "-" + atomType3
	key2 := atomType3 + "-" + atomType2 + "-" + atomType1

	if params, ok := backboneAngleParams[key1]; ok {
		return params
	}
	if params, ok := backboneAngleParams[key2]; ok {
		return params
	}

	// Default
	return backboneAngleParams["default"]
}
