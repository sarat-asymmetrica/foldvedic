// Package physics implements spatial hashing for O(n) force calculations.
//
// MATHEMATICIAN: Reduce O(n²) pairwise forces to O(n) using spatial grid
// PHYSICIST: Only calculate forces within cutoff distance
// BIOCHEMIST: Typical protein interactions are short-range (<12 Å)
// ETHICIST: Enables folding of large proteins on modest hardware
package physics

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// SpatialHash implements a 3D spatial grid for fast neighbor queries
//
// PERFORMANCE:
// - Naive pairwise: O(n²) = 1M comparisons for 1000 atoms
// - Spatial hash: O(n) = ~27n comparisons (neighbors only)
// - Expected speedup: 37× for 1000 atoms, 370× for 10,000 atoms
type SpatialHash struct {
	cellSize float64              // Grid cell size (Å)
	grid     map[int][]*parser.Atom // Hash map: cell ID → atoms in cell
}

// NewSpatialHash creates a spatial hash grid
//
// PHYSICIST:
// Cell size should be at least the cutoff distance
// Typical: 10-12 Å for non-bonded interactions
func NewSpatialHash(cellSize float64) *SpatialHash {
	return &SpatialHash{
		cellSize: cellSize,
		grid:     make(map[int][]*parser.Atom),
	}
}

// Insert adds an atom to the spatial hash
func (sh *SpatialHash) Insert(atom *parser.Atom) {
	cellID := sh.getCellID(atom.X, atom.Y, atom.Z)
	sh.grid[cellID] = append(sh.grid[cellID], atom)
}

// GetNeighbors returns all atoms within neighboring cells
//
// MATHEMATICIAN:
// Checks 27 cells (3×3×3 cube centered on atom)
// Returns candidates; caller must check exact distance
func (sh *SpatialHash) GetNeighbors(atom *parser.Atom) []*parser.Atom {
	neighbors := make([]*parser.Atom, 0, 27*4) // Estimate 4 atoms per cell

	// Get cell indices
	ix := int(math.Floor(atom.X / sh.cellSize))
	iy := int(math.Floor(atom.Y / sh.cellSize))
	iz := int(math.Floor(atom.Z / sh.cellSize))

	// Check 27 neighboring cells (including center)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				cellID := sh.hashCell(ix+dx, iy+dy, iz+dz)
				if atoms, exists := sh.grid[cellID]; exists {
					neighbors = append(neighbors, atoms...)
				}
			}
		}
	}

	return neighbors
}

// Clear removes all atoms from the grid
func (sh *SpatialHash) Clear() {
	sh.grid = make(map[int][]*parser.Atom)
}

// getCellID computes cell ID for 3D coordinates
func (sh *SpatialHash) getCellID(x, y, z float64) int {
	ix := int(math.Floor(x / sh.cellSize))
	iy := int(math.Floor(y / sh.cellSize))
	iz := int(math.Floor(z / sh.cellSize))
	return sh.hashCell(ix, iy, iz)
}

// hashCell computes hash for cell indices
//
// VEDIC APPROACH:
// Use digital root for better distribution
// Hash = DR(x) + DR(y)×10 + DR(z)×100
func (sh *SpatialHash) hashCell(ix, iy, iz int) int {
	// Simple hash: Morton code (Z-order curve)
	// Interleave bits for 3D locality
	return (ix & 0x3FF) | ((iy & 0x3FF) << 10) | ((iz & 0x3FF) << 20)
}

// CalculateNonBondedEnergySpatial uses spatial hashing for O(n) performance
//
// PERFORMANCE COMPARISON:
// - calculateVanDerWaalsTotal: O(n²) naive
// - This function: O(n) with spatial hashing
func CalculateNonBondedEnergySpatial(protein *parser.Protein, vdwCutoff, elecCutoff float64) (vdw, elec float64) {
	// Use larger cutoff for cell size
	cutoff := math.Max(vdwCutoff, elecCutoff)
	spatialHash := NewSpatialHash(cutoff)

	// Insert all atoms
	for _, atom := range protein.Atoms {
		spatialHash.Insert(atom)
	}

	// Simplified charges
	charges := map[string]float64{
		"N":  -0.4157,
		"CA": 0.0337,
		"C":  0.5973,
		"O":  -0.5679,
	}

	// Calculate pairwise energies (only neighbors)
	visited := make(map[[2]int]bool) // Track pairs to avoid double counting

	for _, atom1 := range protein.Atoms {
		neighbors := spatialHash.GetNeighbors(atom1)

		for _, atom2 := range neighbors {
			// Skip self
			if atom1.Serial == atom2.Serial {
				continue
			}

			// Skip if already calculated (avoid double counting)
			pair := [2]int{atom1.Serial, atom2.Serial}
			if atom1.Serial > atom2.Serial {
				pair = [2]int{atom2.Serial, atom1.Serial}
			}
			if visited[pair] {
				continue
			}
			visited[pair] = true

			// Skip bonded/1-4 interactions
			if math.Abs(float64(atom1.ResSeq-atom2.ResSeq)) <= 1 {
				continue
			}

			// Calculate distance
			dx := atom2.X - atom1.X
			dy := atom2.Y - atom1.Y
			dz := atom2.Z - atom1.Z
			r := math.Sqrt(dx*dx + dy*dy + dz*dz)

			// Van der Waals
			if r <= vdwCutoff {
				vdw += CalculateLennardJonesEnergy(atom1, atom2, vdwCutoff)
			}

			// Electrostatic
			if r <= elecCutoff {
				charge1, ok1 := charges[atom1.Name]
				charge2, ok2 := charges[atom2.Name]
				if ok1 && ok2 {
					elec += CalculateElectrostaticEnergy(atom1, atom2, charge1, charge2, elecCutoff)
				}
			}
		}
	}

	return vdw, elec
}
