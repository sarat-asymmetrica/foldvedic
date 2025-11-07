package sampling

import (
	"math"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/folding"
)

// DiversityMetric calculates overall structural diversity of a set of proteins
// Returns average pairwise dihedral RMSD
func DiversityMetric(proteins []*folding.Protein) float64 {
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
func CalculateDihedralRMSD(protein1, protein2 *folding.Protein) float64 {
	if len(protein1.Residues) != len(protein2.Residues) {
		return 999999.9
	}

	sumSqDiff := 0.0
	count := 0

	for i := range protein1.Residues {
		res1 := protein1.Residues[i]
		res2 := protein2.Residues[i]

		// Compare phi angles
		if !math.IsNaN(res1.Phi) && !math.IsNaN(res2.Phi) {
			diff := angleDifference(res1.Phi, res2.Phi)
			sumSqDiff += diff * diff
			count++
		}

		// Compare psi angles
		if !math.IsNaN(res1.Psi) && !math.IsNaN(res2.Psi) {
			diff := angleDifference(res1.Psi, res2.Psi)
			sumSqDiff += diff * diff
			count++
		}

		// Optionally compare omega (usually ~180°, but can vary for cis-peptides)
		if !math.IsNaN(res1.Omega) && !math.IsNaN(res2.Omega) {
			diff := angleDifference(res1.Omega, res2.Omega)
			sumSqDiff += diff * diff
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	// Convert to degrees for interpretability
	rmsRadians := math.Sqrt(sumSqDiff / float64(count))
	return rmsRadians * 180.0 / math.Pi
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

func CalculateEnsembleDiversity(proteins []*folding.Protein) EnsembleDiversityMetrics {
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
func SelectMaximallyDiverseSubset(proteins []*folding.Protein, k int) []*folding.Protein {
	if len(proteins) <= k {
		return proteins
	}

	selected := make([]*folding.Protein, 0, k)

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
func CalculateConformationalEntropy(proteins []*folding.Protein) float64 {
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
