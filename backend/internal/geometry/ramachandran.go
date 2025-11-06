// Package geometry provides geometric calculations for protein structures.
//
// BIOCHEMIST: Ramachandran plot analysis - phi/psi angles define allowed conformational space
// PHYSICIST: Dihedral angle calculations using vector geometry and cross products
// MATHEMATICIAN: Bijective mapping to quaternion space for smooth interpolation
// ETHICIST: Well-documented formulas with literature citations for reproducibility
package geometry

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// Vector3 represents a 3D vector
type Vector3 struct {
	X, Y, Z float64
}

// Sub subtracts two vectors
func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// Cross computes the cross product of two vectors
func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Dot computes the dot product of two vectors
func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Magnitude returns the length of the vector
func (v Vector3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns a unit vector in the same direction
func (v Vector3) Normalize() Vector3 {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector3{X: 0, Y: 0, Z: 0}
	}
	return Vector3{X: v.X / mag, Y: v.Y / mag, Z: v.Z / mag}
}

// Add two vectors
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// Scale vector by scalar
func (v Vector3) Scale(s float64) Vector3 {
	return Vector3{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

// Length returns the length of the vector (alias for Magnitude)
func (v Vector3) Length() float64 {
	return v.Magnitude()
}

// RotateByQuaternion rotates vector by quaternion
func (v Vector3) RotateByQuaternion(q Quaternion) Vector3 {
	// Convert quaternion to rotation matrix
	qx, qy, qz, qw := q.X, q.Y, q.Z, q.W

	// Rotation matrix elements
	r00 := 1 - 2*(qy*qy+qz*qz)
	r01 := 2 * (qx*qy - qz*qw)
	r02 := 2 * (qx*qz + qy*qw)

	r10 := 2 * (qx*qy + qz*qw)
	r11 := 1 - 2*(qx*qx+qz*qz)
	r12 := 2 * (qy*qz - qx*qw)

	r20 := 2 * (qx*qz - qy*qw)
	r21 := 2 * (qy*qz + qx*qw)
	r22 := 1 - 2*(qx*qx+qy*qy)

	// Apply rotation
	return Vector3{
		X: r00*v.X + r01*v.Y + r02*v.Z,
		Y: r10*v.X + r11*v.Y + r12*v.Z,
		Z: r20*v.X + r21*v.Y + r22*v.Z,
	}
}

// RamachandranAngles holds phi and psi dihedral angles for a residue
type RamachandranAngles struct {
	Phi float64 // Phi dihedral angle (radians)
	Psi float64 // Psi dihedral angle (radians)
}

// CalculateRamachandran computes phi and psi angles for all residues in a protein
//
// BIOCHEMIST:
//   - Phi (φ): C(i-1) - N(i) - Cα(i) - C(i) dihedral angle
//   - Psi (ψ): N(i) - Cα(i) - C(i) - N(i+1) dihedral angle
//   - Terminal residues have undefined angles (first has no phi, last has no psi)
//
// PHYSICIST:
//   - Dihedral angle calculated using atan2 for proper quadrant handling
//   - Range: [-π, +π] radians (or [-180°, +180°])
//
// Citation: Ramachandran, G. N., et al. (1963). "Stereochemistry of polypeptide chain configurations."
// J. Mol. Biol. 7.1: 95-99.
func CalculateRamachandran(protein *parser.Protein) []RamachandranAngles {
	residues := protein.Residues
	angles := make([]RamachandranAngles, len(residues))

	for i := range residues {
		// Phi requires previous residue's C
		if i > 0 && residues[i-1].HasCompleteBackbone() && residues[i].HasCompleteBackbone() {
			angles[i].Phi = calculateDihedral(
				atomToVector(residues[i-1].C),
				atomToVector(residues[i].N),
				atomToVector(residues[i].CA),
				atomToVector(residues[i].C),
			)
		} else {
			angles[i].Phi = math.NaN() // Undefined for N-terminal
		}

		// Psi requires next residue's N
		if i < len(residues)-1 && residues[i].HasCompleteBackbone() && residues[i+1].HasCompleteBackbone() {
			angles[i].Psi = calculateDihedral(
				atomToVector(residues[i].N),
				atomToVector(residues[i].CA),
				atomToVector(residues[i].C),
				atomToVector(residues[i+1].N),
			)
		} else {
			angles[i].Psi = math.NaN() // Undefined for C-terminal
		}
	}

	return angles
}

// calculateDihedral computes the dihedral angle defined by four points
//
// PHYSICIST:
// Dihedral angle between planes (p1,p2,p3) and (p2,p3,p4)
// Formula uses cross products to get plane normals, then atan2 for signed angle
//
// Algorithm from:
// Prisant, M. G., et al. (2020). "New tools in MolProbity validation."
// Protein Sci. 29.1: 315-329.
//
// MATHEMATICIAN:
// Returns angle in radians [-π, +π]
// atan2 ensures proper quadrant (no ambiguity from acos)
func calculateDihedral(p1, p2, p3, p4 Vector3) float64 {
	// Vectors along bonds
	b1 := p2.Sub(p1) // Vector from p1 to p2
	b2 := p3.Sub(p2) // Vector from p2 to p3
	b3 := p4.Sub(p3) // Vector from p3 to p4

	// Normal vectors to the two planes
	n1 := b1.Cross(b2) // Normal to plane (p1, p2, p3)
	n2 := b2.Cross(b3) // Normal to plane (p2, p3, p4)

	// Normalize the cross product of b2 (central bond)
	b2norm := b2.Normalize()

	// Calculate angle using atan2 for proper sign
	// m1 is the cross product of the two normals projected onto b2
	m1 := n1.Cross(b2norm)

	x := n1.Dot(n2)
	y := m1.Dot(n2)

	return math.Atan2(y, x)
}

// atomToVector converts an Atom to a Vector3
func atomToVector(atom *parser.Atom) Vector3 {
	return Vector3{X: atom.X, Y: atom.Y, Z: atom.Z}
}

// ToDegressPhi converts phi angle from radians to degrees
func (ra RamachandranAngles) ToDegressPhi() float64 {
	return ra.Phi * 180.0 / math.Pi
}

// ToDegressPsi converts psi angle from radians to degrees
func (ra RamachandranAngles) ToDegressPsi() float64 {
	return ra.Psi * 180.0 / math.Pi
}

// IsInAllowedRegion checks if (phi, psi) is in an allowed Ramachandran region
//
// BIOCHEMIST:
// Allowed regions from Ramachandran plot:
// - Alpha helix: φ ≈ -60°, ψ ≈ -45° (±30° tolerance)
// - Beta sheet: φ ≈ -120°, ψ ≈ +120° (±40° tolerance)
// - Left-handed helix: φ ≈ +60°, ψ ≈ +45° (rare, mainly for Glycine)
//
// Simplified check - full Ramachandran regions are more complex
// For production, use DSSP or MolProbity definitions
func (ra RamachandranAngles) IsInAllowedRegion() bool {
	if math.IsNaN(ra.Phi) || math.IsNaN(ra.Psi) {
		return false
	}

	phiDeg := ra.ToDegressPhi()
	psiDeg := ra.ToDegressPsi()

	// Alpha helix region: φ ≈ -60°, ψ ≈ -45°
	if phiDeg >= -90 && phiDeg <= -30 && psiDeg >= -75 && psiDeg <= -15 {
		return true
	}

	// Beta sheet region: φ ≈ -120°, ψ ≈ +120°
	if phiDeg >= -160 && phiDeg <= -80 && psiDeg >= 80 && psiDeg <= 160 {
		return true
	}

	// Left-handed helix region (mainly Gly): φ ≈ +60°, ψ ≈ +45°
	if phiDeg >= 30 && phiDeg <= 90 && psiDeg >= 15 && psiDeg <= 75 {
		return true
	}

	// Additional allowed regions (extended, turn conformations)
	// Simplified - full Ramachandran has more nuanced regions
	if phiDeg >= -180 && phiDeg <= -30 {
		return true // Most of left half of Ramachandran plot
	}

	return false
}
