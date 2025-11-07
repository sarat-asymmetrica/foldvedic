package physics

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// HBond represents a hydrogen bond between donor and acceptor
type HBond struct {
	DonorResidue    *parser.Residue
	AcceptorResidue *parser.Residue
	DonorAtom       *parser.Atom // N-H hydrogen
	AcceptorAtom    *parser.Atom // O=C oxygen
	Distance        float64       // Å
	Angle           float64       // degrees (N-H···O angle)
	Energy          float64       // kcal/mol
}

// DetectHydrogenBonds finds all hydrogen bonds in a protein
// Uses standard geometric criteria from structural biology
func DetectHydrogenBonds(protein *parser.Protein) []HBond {
	hbonds := []HBond{}

	// Hydrogen bond geometric criteria
	const (
		minDistance = 2.5  // Å (minimum N···O distance)
		maxDistance = 3.5  // Å (maximum N···O distance)
		minAngle    = 120.0 // degrees (minimum N-H···O angle)
		maxAngle    = 180.0 // degrees (maximum, ideal is 180°)
	)

	// WAVE 11.4 FIX: Look for explicit H atoms first
	// Build map of residue -> backbone H atom
	hAtomMap := make(map[int]*parser.Atom) // residue SeqNum -> H atom
	for _, atom := range protein.Atoms {
		// Look for backbone amide H (named "H" or "HN")
		if atom.Element == "H" && (atom.Name == "H" || atom.Name == "HN") {
			hAtomMap[atom.ResSeq] = atom
		}
	}

	// Find all potential donors (backbone N atoms)
	donors := []*parser.Atom{}
	for _, residue := range protein.Residues {
		if residue.N != nil {
			donors = append(donors, residue.N)
		}
	}

	// Find all potential acceptors (backbone O atoms)
	acceptors := []*parser.Atom{}
	for _, residue := range protein.Residues {
		if residue.O != nil {
			acceptors = append(acceptors, residue.O)
		}
	}

	// Check all donor-acceptor pairs
	for _, donor := range donors {
		donorResidue := findResidueForAtom(protein, donor)
		if donorResidue == nil {
			continue
		}

		// Try to find explicit H atom for this donor
		donorH := hAtomMap[donorResidue.SeqNum]

		// Fallback: Get donor CA for angle calculation if no H atom
		donorCA := donorResidue.CA
		if donorH == nil && donorCA == nil {
			continue
		}

		for _, acceptor := range acceptors {
			acceptorResidue := findResidueForAtom(protein, acceptor)
			if acceptorResidue == nil {
				continue
			}

			// Don't allow H-bonds between adjacent residues (too close in sequence)
			if abs(donorResidue.SeqNum-acceptorResidue.SeqNum) <= 1 {
				continue
			}

			var distance float64
			var angle float64

			if donorH != nil {
				// Use explicit H atom for accurate geometry
				// H-bond distance is H···O (not N···O)
				distance = calculateDistance(donorH, acceptor)

				// H-bond angle is N-H···O
				angle = calculateHBondAngleWithH(donor, donorH, acceptor)
			} else {
				// Fallback: Use N···O distance and N-CA proxy
				distance = calculateDistance(donor, acceptor)
				angle = calculateHBondAngle(donor, donorCA, acceptor)
			}

			// Check distance criterion
			// For H···O: 1.5 - 2.5 Å
			// For N···O: 2.5 - 3.5 Å
			var minDist, maxDist float64
			if donorH != nil {
				minDist = 1.5
				maxDist = 2.5
			} else {
				minDist = minDistance
				maxDist = maxDistance
			}

			if distance < minDist || distance > maxDist {
				continue
			}

			// Check angle criterion
			if angle < minAngle {
				continue
			}

			// Calculate energy using empirical potential
			energy := calculateHBondEnergy(distance, angle)

			hbond := HBond{
				DonorResidue:    donorResidue,
				AcceptorResidue: acceptorResidue,
				DonorAtom:       donor,
				AcceptorAtom:    acceptor,
				Distance:        distance,
				Angle:           angle,
				Energy:          energy,
			}

			hbonds = append(hbonds, hbond)
		}
	}

	return hbonds
}

// CalculateHydrogenBondEnergy calculates total H-bond energy for a protein
func CalculateHydrogenBondEnergy(protein *parser.Protein) float64 {
	hbonds := DetectHydrogenBonds(protein)

	totalEnergy := 0.0
	for _, hbond := range hbonds {
		totalEnergy += hbond.Energy
	}

	return totalEnergy
}

// calculateHBondEnergy calculates energy of a single H-bond
// Uses simplified empirical potential from Mayo et al. (1990)
func calculateHBondEnergy(distance, angle float64) float64 {
	// Energy components:
	// 1. Distance-dependent term (Lennard-Jones-like)
	// 2. Angle-dependent term (cos(theta) preference for linearity)

	const (
		optimalDistance = 2.9  // Å (optimal N···O distance)
		energyScale     = -5.0 // kcal/mol (maximum H-bond energy)
	)

	// Distance term: Gaussian-like potential
	distanceTerm := math.Exp(-((distance - optimalDistance) * (distance - optimalDistance)) / 0.2)

	// Angle term: Favor linear H-bonds (180°)
	// cos(180°) = -1, cos(120°) = -0.5
	angleRad := angle * math.Pi / 180.0
	angleTerm := (1.0 + math.Cos(angleRad)) / 2.0 // Maps [120°, 180°] to [0.25, 1.0]

	// Combined energy
	energy := energyScale * distanceTerm * angleTerm

	return energy
}

// calculateHBondAngleWithH calculates N-H···O angle using explicit H atom
// WAVE 11.4: NEW FUNCTION for accurate H-bond geometry
//
// The N-H···O angle is measured at H in the N-H-O chain.
// For a linear H-bond, this angle should be 180° (H between N and O in a line).
// We measure angle between vectors H→N (backwards) and H→O (forward).
func calculateHBondAngleWithH(donorN, donorH, acceptorO *parser.Atom) float64 {
	// Vector from H to N (backwards along N-H bond)
	hToN := Vector3D{
		X: donorN.X - donorH.X,
		Y: donorN.Y - donorH.Y,
		Z: donorN.Z - donorH.Z,
	}

	// Vector from H to O (forward to acceptor)
	hToO := Vector3D{
		X: acceptorO.X - donorH.X,
		Y: acceptorO.Y - donorH.Y,
		Z: acceptorO.Z - donorH.Z,
	}

	// Normalize vectors
	hToN = normalizeVector(hToN)
	hToO = normalizeVector(hToO)

	// Calculate angle using dot product
	// If vectors are parallel (same direction), angle = 0°
	// If vectors are antiparallel (opposite direction), angle = 180°
	// For linear H-bond: H→N and H→O point opposite, so angle = 180°
	dotProduct := hToN.X*hToO.X + hToN.Y*hToO.Y + hToN.Z*hToO.Z

	// Clamp to [-1, 1] to avoid numerical errors
	if dotProduct > 1.0 {
		dotProduct = 1.0
	}
	if dotProduct < -1.0 {
		dotProduct = -1.0
	}

	angleRad := math.Acos(dotProduct)
	angleDeg := angleRad * 180.0 / math.Pi

	return angleDeg
}

// calculateHBondAngle calculates N-H···O angle
// Uses N-CA direction as proxy for N-H since we don't have explicit H atoms
// FALLBACK method when H atoms not available
func calculateHBondAngle(donorN, donorCA, acceptorO *parser.Atom) float64 {
	// Vector from N to CA (opposite of N-H direction)
	nToCA := Vector3D{
		X: donorCA.X - donorN.X,
		Y: donorCA.Y - donorN.Y,
		Z: donorCA.Z - donorN.Z,
	}

	// Vector from N to O
	nToO := Vector3D{
		X: acceptorO.X - donorN.X,
		Y: acceptorO.Y - donorN.Y,
		Z: acceptorO.Z - donorN.Z,
	}

	// Normalize vectors
	nToCA = normalizeVector(nToCA)
	nToO = normalizeVector(nToO)

	// Calculate angle using dot product
	dotProduct := nToCA.X*nToO.X + nToCA.Y*nToO.Y + nToCA.Z*nToO.Z

	// Clamp to [-1, 1] to avoid numerical errors
	if dotProduct > 1.0 {
		dotProduct = 1.0
	}
	if dotProduct < -1.0 {
		dotProduct = -1.0
	}

	angleRad := math.Acos(dotProduct)
	angleDeg := angleRad * 180.0 / math.Pi

	// Since we're using N-CA (opposite of N-H), invert the angle
	angleDeg = 180.0 - angleDeg

	return angleDeg
}

// Vector3D is a simple 3D vector
type Vector3D struct {
	X, Y, Z float64
}

func normalizeVector(v Vector3D) Vector3D {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if length < 1e-10 {
		return Vector3D{0, 0, 0}
	}
	return Vector3D{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func calculateDistance(atom1, atom2 *parser.Atom) float64 {
	dx := atom1.X - atom2.X
	dy := atom1.Y - atom2.Y
	dz := atom1.Z - atom2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func findResidueForAtom(protein *parser.Protein, atom *parser.Atom) *parser.Residue {
	for _, residue := range protein.Residues {
		if residue.N == atom || residue.CA == atom || residue.C == atom || residue.O == atom {
			return residue
		}
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetHBondStatistics returns statistics about H-bonds in a protein
type HBondStatistics struct {
	NumHBonds       int
	AverageDistance float64
	AverageAngle    float64
	AverageEnergy   float64
	TotalEnergy     float64
	HelixHBonds     int // i → i+4 pattern
	SheetHBonds     int // Between different strands
	LoopHBonds      int // Everything else
}

func GetHBondStatistics(protein *parser.Protein) HBondStatistics {
	hbonds := DetectHydrogenBonds(protein)

	stats := HBondStatistics{
		NumHBonds: len(hbonds),
	}

	if len(hbonds) == 0 {
		return stats
	}

	sumDistance := 0.0
	sumAngle := 0.0
	sumEnergy := 0.0

	for _, hbond := range hbonds {
		sumDistance += hbond.Distance
		sumAngle += hbond.Angle
		sumEnergy += hbond.Energy

		// Classify H-bond type by sequence separation
		separation := abs(hbond.DonorResidue.SeqNum - hbond.AcceptorResidue.SeqNum)
		if separation == 4 {
			stats.HelixHBonds++ // α-helix pattern (i → i+4)
		} else if separation >= 5 {
			stats.SheetHBonds++ // β-sheet or long-range
		} else {
			stats.LoopHBonds++ // Short-range loops
		}
	}

	stats.AverageDistance = sumDistance / float64(len(hbonds))
	stats.AverageAngle = sumAngle / float64(len(hbonds))
	stats.AverageEnergy = sumEnergy / float64(len(hbonds))
	stats.TotalEnergy = sumEnergy

	return stats
}
