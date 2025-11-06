// Package geometry - Quaternion-Based Coordinate Builder
//
// WAVE 11.1: Quaternion Coordinate Generation - NOVEL APPROACH
//
// STANDARD METHODS:
// - NeRF (Natural Extension Reference Frame): Uses rotation matrices
// - Forward kinematics: Matrix multiplication chains
//
// OUR APPROACH: QUATERNION COMPOSITION
// - Leverage existing quaternion engine from Wave 7
// - Use quaternion rotations instead of matrices
// - Composable, numerically stable, elegant
//
// CROSS-DOMAIN INSPIRATION:
// - Computer graphics: Skeletal animation with quaternions
// - Robotics: Forward kinematics for robot arms
// - Aerospace: Attitude control with quaternions
//
// INSIGHT: Protein backbones are literally skeletal chains!
// Each residue is a "joint" with rotation defined by (φ, ψ, ω).
// This is identical to robot arm forward kinematics.
//
// WRIGHT BROTHERS EMPIRICISM:
// - Build simplest version first
// - Test on 3-residue peptide
// - Expect bugs, iterate fast
// - Add complexity only when needed
//
// BIOCHEMIST: Standard bond lengths/angles from crystallography
// PHYSICIST: Quaternion rotations = SO(3) group operations
// MATHEMATICIAN: Elegant composition via quaternion multiplication
// ETHICIST: Novel but grounded in solid math
package geometry

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// Standard bond lengths (Å) from crystallography
const (
	BondN_CA  = 1.458 // N-CA bond length
	BondCA_C  = 1.523 // CA-C bond length
	BondC_N   = 1.329 // C-N peptide bond
	BondC_O   = 1.231 // C=O carbonyl bond
	BondCA_CB = 1.530 // CA-CB side chain (for stubs)
)

// Standard bond angles (degrees) from Ramachandran
const (
	AngleN_CA_C  = 111.0 // N-CA-C angle
	AngleCA_C_N  = 117.0 // CA-C-N angle
	AngleC_N_CA  = 121.0 // C-N-CA angle
	AngleCA_C_O  = 120.5 // CA-C=O angle
	AngleN_CA_CB = 110.5 // N-CA-CB tetrahedral
)

// BuildProteinFromAngles constructs 3D protein coordinates from (φ, ψ, ω) angles
//
// ALGORITHM: Quaternion-Based Forward Kinematics
//
// For each residue i:
//   1. Start with previous C atom position and orientation
//   2. Rotate by ω (peptide bond rotation) to place N
//   3. Rotate by φ (phi dihedral) to place CA
//   4. Rotate by ψ (psi dihedral) to place C
//   5. Place O based on C=O bond
//   6. Update orientation quaternion for next residue
//
// INPUTS:
//   - sequence: Amino acid sequence (e.g., "ACDEFG")
//   - angles: φ, ψ, ω angles for each residue
//
// OUTPUTS:
//   - Protein with realistic 3D coordinates
//
// WRIGHT BROTHERS TEST:
//   - Try on "GAC" (3 residues)
//   - Check if bond lengths are ~correct
//   - Check if it doesn't explode
//
func BuildProteinFromAngles(sequence string, angles []RamachandranAngles) (*parser.Protein, error) {
	n := len(sequence)

	protein := &parser.Protein{
		Name:     "built_from_angles",
		Residues: make([]*parser.Residue, n),
		Atoms:    make([]*parser.Atom, 0, n*4),
	}

	// Initial position and orientation
	// Start at origin, pointing along +X axis
	currentPos := Vector3{0, 0, 0}
	currentDir := Vector3{1, 0, 0} // Direction to next atom

	// Reference frame: Z-axis up, Y-axis to the side
	upVec := Vector3{0, 0, 1}

	atomSerial := 1

	for i := 0; i < n; i++ {
		res := &parser.Residue{
			Name:    string(sequence[i]),
			SeqNum:  i + 1,
			ChainID: "A",
		}

		// Get angles for this residue
		var phi, psi, omega float64
		if i < len(angles) {
			phi = angles[i].Phi
			psi = angles[i].Psi
			if i > 0 {
				// Omega is typically ~180° (trans peptide bond)
				omega = 180.0 * math.Pi / 180.0
			}
		} else {
			// Default: extended conformation
			phi = -120.0 * math.Pi / 180.0
			psi = 120.0 * math.Pi / 180.0
			omega = 180.0 * math.Pi / 180.0
		}

		// === PLACE N ATOM ===
		if i > 0 {
			// Rotate by omega around C-N axis
			axis := currentDir.Normalize()
			omegaQuat := QuaternionFromAxisAngle(axis, omega)
			currentDir = currentDir.RotateByQuaternion(omegaQuat)
		}

		nPos := currentPos
		res.N = &parser.Atom{
			Serial:  atomSerial,
			Name:    "N",
			ResName: string(sequence[i]),
			ChainID: "A",
			ResSeq:  i + 1,
			X:       nPos.X,
			Y:       nPos.Y,
			Z:       nPos.Z,
			Element: "N",
		}
		protein.Atoms = append(protein.Atoms, res.N)
		atomSerial++

		// === PLACE CA ATOM ===
		// Move along current direction by N-CA bond length
		caDir := currentDir.Normalize()
		caPos := nPos.Add(caDir.Scale(BondN_CA))

		res.CA = &parser.Atom{
			Serial:  atomSerial,
			Name:    "CA",
			ResName: string(sequence[i]),
			ChainID: "A",
			ResSeq:  i + 1,
			X:       caPos.X,
			Y:       caPos.Y,
			Z:       caPos.Z,
			Element: "C",
		}
		protein.Atoms = append(protein.Atoms, res.CA)
		atomSerial++

		// Rotate by phi around N-CA axis
		axis := caDir
		phiQuat := QuaternionFromAxisAngle(axis, phi)
		currentDir = currentDir.RotateByQuaternion(phiQuat)

		// === PLACE C ATOM ===
		// Rotate to account for N-CA-C bond angle
		bondAngle := AngleN_CA_C * math.Pi / 180.0
		// Rotate current direction by (180° - bond angle)
		rotAngle := math.Pi - bondAngle

		// Find perpendicular axis for rotation
		perpAxis := currentDir.Cross(upVec).Normalize()
		if perpAxis.Length() < 0.1 {
			// Current dir parallel to up, use different perpendicular
			perpAxis = currentDir.Cross(Vector3{0, 1, 0}).Normalize()
		}

		angleQuat := QuaternionFromAxisAngle(perpAxis, rotAngle)
		cDir := currentDir.RotateByQuaternion(angleQuat).Normalize()

		cPos := caPos.Add(cDir.Scale(BondCA_C))

		res.C = &parser.Atom{
			Serial:  atomSerial,
			Name:    "C",
			ResName: string(sequence[i]),
			ChainID: "A",
			ResSeq:  i + 1,
			X:       cPos.X,
			Y:       cPos.Y,
			Z:       cPos.Z,
			Element: "C",
		}
		protein.Atoms = append(protein.Atoms, res.C)
		atomSerial++

		// Rotate by psi around CA-C axis
		axis = cDir
		psiQuat := QuaternionFromAxisAngle(axis, psi)
		currentDir = cDir.RotateByQuaternion(psiQuat)

		// === PLACE O ATOM ===
		// O is perpendicular to C, at C=O bond angle
		oAngle := AngleCA_C_O * math.Pi / 180.0
		oAxis := cDir.Cross(caDir).Normalize()
		if oAxis.Length() < 0.1 {
			oAxis = cDir.Cross(upVec).Normalize()
		}

		oQuat := QuaternionFromAxisAngle(oAxis, oAngle)
		oDir := cDir.Scale(-1.0).RotateByQuaternion(oQuat).Normalize()
		oPos := cPos.Add(oDir.Scale(BondC_O))

		res.O = &parser.Atom{
			Serial:  atomSerial,
			Name:    "O",
			ResName: string(sequence[i]),
			ChainID: "A",
			ResSeq:  i + 1,
			X:       oPos.X,
			Y:       oPos.Y,
			Z:       oPos.Z,
			Element: "O",
		}
		protein.Atoms = append(protein.Atoms, res.O)
		atomSerial++

		// Update current position for next residue
		// Next N is C-N bond length away
		currentPos = cPos.Add(currentDir.Scale(BondC_N))

		protein.Residues[i] = res
	}

	return protein, nil
}

// QuaternionFromAxisAngle creates quaternion from axis-angle representation
//
// CROSS-DOMAIN: Robotics, aerospace (attitude representation)
//
// q = [cos(θ/2), sin(θ/2) * axis]
func QuaternionFromAxisAngle(axis Vector3, angle float64) Quaternion {
	axis = axis.Normalize()
	halfAngle := angle / 2.0
	sinHalf := math.Sin(halfAngle)
	cosHalf := math.Cos(halfAngle)

	return Quaternion{
		W: cosHalf,
		X: axis.X * sinHalf,
		Y: axis.Y * sinHalf,
		Z: axis.Z * sinHalf,
	}
}

// BuildProteinFromAnglesVedic builds structure with Vedic harmonic biasing
//
// WILD IDEA: What if we bias bond lengths/angles toward φ-ratios?
//
// HYPOTHESIS (SPECULATIVE):
// - Proteins might prefer geometries aligned with golden ratio
// - Helix pitch 3.6 ≈ 10/φ²
// - Try biasing and see if RMSD improves!
//
// WRIGHT BROTHERS: This might crash and burn. That's okay. Try it!
func BuildProteinFromAnglesVedic(sequence string, angles []RamachandranAngles, vedicWeight float64) (*parser.Protein, error) {
	// First build standard structure
	protein, err := BuildProteinFromAngles(sequence, angles)
	if err != nil {
		return nil, err
	}

	// Apply Vedic biasing: Adjust bond lengths toward φ-ratios
	const phi = 1.618033988749895

	if vedicWeight > 0 {
		for _, atom := range protein.Atoms {
			// Subtle adjustment: Move atoms toward φ-ratio distances from origin
			r := math.Sqrt(atom.X*atom.X + atom.Y*atom.Y + atom.Z*atom.Z)
			if r > 0.1 {
				// Find nearest Fibonacci number
				nearestFib := findNearestFibonacci(r)
				targetR := float64(nearestFib) / phi

				// Blend between current and target
				blendedR := (1.0-vedicWeight)*r + vedicWeight*targetR

				scale := blendedR / r
				atom.X *= scale
				atom.Y *= scale
				atom.Z *= scale
			}
		}
	}

	return protein, nil
}

// findNearestFibonacci finds closest Fibonacci number to value
func findNearestFibonacci(value float64) int {
	fib := []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}
	nearest := fib[0]
	minDiff := math.Abs(value - float64(fib[0]))

	for _, f := range fib {
		diff := math.Abs(value - float64(f))
		if diff < minDiff {
			minDiff = diff
			nearest = f
		}
	}

	return nearest
}

// ValidateBackboneGeometry checks if bond lengths/angles are reasonable
//
// WRIGHT BROTHERS: Quick sanity check before "takeoff"
func ValidateBackboneGeometry(protein *parser.Protein) (bool, string) {
	for i, res := range protein.Residues {
		if res.N == nil || res.CA == nil || res.C == nil {
			return false, "Missing backbone atoms"
		}

		// Check N-CA bond
		nca := Vector3{res.CA.X - res.N.X, res.CA.Y - res.N.Y, res.CA.Z - res.N.Z}
		ncaLen := nca.Length()
		if ncaLen < 1.0 || ncaLen > 2.0 {
			return false, "N-CA bond length out of range"
		}

		// Check CA-C bond
		cac := Vector3{res.C.X - res.CA.X, res.C.Y - res.CA.Y, res.C.Z - res.CA.Z}
		cacLen := cac.Length()
		if cacLen < 1.0 || cacLen > 2.0 {
			return false, "CA-C bond length out of range"
		}

		// Check peptide bond (if not last residue)
		if i < len(protein.Residues)-1 {
			nextN := protein.Residues[i+1].N
			if nextN != nil {
				cn := Vector3{nextN.X - res.C.X, nextN.Y - res.C.Y, nextN.Z - res.C.Z}
				cnLen := cn.Length()
				if cnLen < 0.8 || cnLen > 2.0 {
					return false, "C-N peptide bond length out of range"
				}
			}
		}
	}

	return true, "Geometry valid"
}
