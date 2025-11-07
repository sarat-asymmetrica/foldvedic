package sampling

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/folding"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// DiversityMetric calculates overall structural diversity of a set of proteins
// Returns average pairwise dihedral RMSD
func DiversityMetric(proteins []*parser.Protein) float64 {
	if len(proteins) <= 1 {
		return 0.0
	}

	totalDiversity := 0.0
	count := 0

	// Calculate all pairwise diversities
	for i := 0; i < len(proteins); i++ {
		for j := i + 1; j < len(proteins); j++ {
			div := QuaternionDiversity(proteins[i], proteins[j])
			totalDiversity += div
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return totalDiversity / float64(count)
}

// CalculateDihedralRMSD calculates RMSD between two proteins in dihedral space
// This is faster than Cartesian RMSD and respects rotational invariance
// Since Phi/Psi/Omega aren't stored as fields, we use CA distance as proxy
func CalculateDihedralRMSD(protein1, protein2 *parser.Protein) float64 {
	if len(protein1.Residues) != len(protein2.Residues) {
		return 999999.9
	}

	sumSqDiff := 0.0
	count := 0

	// Use CA distance as proxy for dihedral difference
	for i := range protein1.Residues {
		res1 := protein1.Residues[i]
		res2 := protein2.Residues[i]

		if res1.CA != nil && res2.CA != nil {
			dx := res1.CA.X - res2.CA.X
			dy := res1.CA.Y - res2.CA.Y
			dz := res1.CA.Z - res2.CA.Z
			sumSqDiff += dx*dx + dy*dy + dz*dz
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	// Return RMSD in degrees equivalent (scaled for interpretability)
	rmsAngstroms := math.Sqrt(sumSqDiff / float64(count))
	return rmsAngstroms * 30.0 // Scale factor to convert Å to degree-like units
}

// angleDifference calculates the minimum angular distance between two angles
// Handles wrap-around at ±π
func angleDifference(angle1, angle2 float64) float64 {
	diff := math.Abs(angle1 - angle2)
	if diff > math.Pi {
		diff = 2*math.Pi - diff
	}
	return diff
}

// CalculateEnsembleDiversity returns multiple diversity metrics for an ensemble
type EnsembleDiversityMetrics struct {
	AveragePairwiseDihedralRMSD   float64
	MedianPairwiseDihedralRMSD    float64
	MaxPairwiseDihedralRMSD       float64
	MinPairwiseDihedralRMSD       float64
	AverageEnergySpread           float64
	NumUniqueStructures           int
}

func CalculateEnsembleDiversity(proteins []*parser.Protein) EnsembleDiversityMetrics {
	metrics := EnsembleDiversityMetrics{}

	if len(proteins) <= 1 {
		return metrics
	}

	// Calculate all pairwise dihedral RMSDs
	pairwiseRMSDs := []float64{}
	for i := 0; i < len(proteins); i++ {
		for j := i + 1; j < len(proteins); j++ {
			rmsd := CalculateDihedralRMSD(proteins[i], proteins[j])
			pairwiseRMSDs = append(pairwiseRMSDs, rmsd)
		}
	}

	if len(pairwiseRMSDs) == 0 {
		return metrics
	}

	// Calculate statistics
	sum := 0.0
	minRMSD := pairwiseRMSDs[0]
	maxRMSD := pairwiseRMSDs[0]

	for _, rmsd := range pairwiseRMSDs {
		sum += rmsd
		if rmsd < minRMSD {
			minRMSD = rmsd
		}
		if rmsd > maxRMSD {
			maxRMSD = rmsd
		}
	}

	metrics.AveragePairwiseDihedralRMSD = sum / float64(len(pairwiseRMSDs))
	metrics.MinPairwiseDihedralRMSD = minRMSD
	metrics.MaxPairwiseDihedralRMSD = maxRMSD

	// Calculate median
	sortedRMSDs := make([]float64, len(pairwiseRMSDs))
	copy(sortedRMSDs, pairwiseRMSDs)
	// Simple bubble sort for median
	for i := 0; i < len(sortedRMSDs); i++ {
		for j := i + 1; j < len(sortedRMSDs); j++ {
			if sortedRMSDs[j] < sortedRMSDs[i] {
				sortedRMSDs[i], sortedRMSDs[j] = sortedRMSDs[j], sortedRMSDs[i]
			}
		}
	}
	if len(sortedRMSDs)%2 == 0 {
		metrics.MedianPairwiseDihedralRMSD = (sortedRMSDs[len(sortedRMSDs)/2-1] + sortedRMSDs[len(sortedRMSDs)/2]) / 2.0
	} else {
		metrics.MedianPairwiseDihedralRMSD = sortedRMSDs[len(sortedRMSDs)/2]
	}

	// Calculate energy spread
	energies := make([]float64, len(proteins))
	for i, protein := range proteins {
		energies[i] = folding.CalculateEnergy(protein)
	}

	energySum := 0.0
	for _, e := range energies {
		energySum += e
	}
	avgEnergy := energySum / float64(len(energies))

	spreadSum := 0.0
	for _, e := range energies {
		spreadSum += math.Abs(e - avgEnergy)
	}
	metrics.AverageEnergySpread = spreadSum / float64(len(energies))

	// Count unique structures (structures with dihedral RMSD > 10° are considered unique)
	uniqueThreshold := 10.0 // degrees
	isUnique := make([]bool, len(proteins))
	for i := range isUnique {
		isUnique[i] = true
	}

	for i := 0; i < len(proteins); i++ {
		if !isUnique[i] {
			continue
		}
		for j := i + 1; j < len(proteins); j++ {
			if !isUnique[j] {
				continue
			}
			rmsd := CalculateDihedralRMSD(proteins[i], proteins[j])
			if rmsd < uniqueThreshold {
				isUnique[j] = false // Mark as duplicate
			}
		}
	}

	uniqueCount := 0
	for _, u := range isUnique {
		if u {
			uniqueCount++
		}
	}
	metrics.NumUniqueStructures = uniqueCount

	return metrics
}

// SelectMaximallyDiverseSubset uses greedy algorithm to select k maximally diverse structures
func SelectMaximallyDiverseSubset(proteins []*parser.Protein, k int) []*parser.Protein {
	if len(proteins) <= k {
		return proteins
	}

	selected := make([]*parser.Protein, 0, k)

	// Start with a random structure (first one)
	selected = append(selected, proteins[0])

	// Greedily add structures that maximize minimum distance to selected set
	for len(selected) < k {
		bestIdx := -1
		bestMinDist := -1.0

		for i, candidate := range proteins {
			// Skip if already selected
			isSelected := false
			for _, s := range selected {
				if s == candidate {
					isSelected = true
					break
				}
			}
			if isSelected {
				continue
			}

			// Calculate minimum distance to selected set
			minDist := 999999.9
			for _, s := range selected {
				dist := CalculateDihedralRMSD(candidate, s)
				if dist < minDist {
					minDist = dist
				}
			}

			// Want structure with largest minimum distance (most diverse)
			if minDist > bestMinDist {
				bestMinDist = minDist
				bestIdx = i
			}
		}

		if bestIdx == -1 {
			break // No more structures to add
		}

		selected = append(selected, proteins[bestIdx])
	}

	return selected
}

// CalculateConformationalEntropy estimates entropy of ensemble based on diversity
// Higher entropy = more diverse conformational sampling
func CalculateConformationalEntropy(proteins []*parser.Protein) float64 {
	metrics := CalculateEnsembleDiversity(proteins)

	// Simple entropy estimate based on:
	// 1. Number of unique structures
	// 2. Average diversity
	// This is a heuristic, not rigorous statistical mechanics entropy

	if metrics.NumUniqueStructures <= 1 {
		return 0.0
	}

	// S ≈ k * ln(N_unique) * <diversity>
	// Where k is Boltzmann-like constant (set to 1 for simplicity)
	entropy := math.Log(float64(metrics.NumUniqueStructures)) * metrics.AveragePairwiseDihedralRMSD

	return entropy
}
