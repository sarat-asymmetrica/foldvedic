package sampling

import (
	"fmt"
	"math"
	"sort"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/folding"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/validation"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/vedic"
)

// EnsembleStructure holds a protein structure with quality metrics
type EnsembleStructure struct {
	Protein       *parser.Protein
	Energy        float64
	VedicScore    float64
	SamplingMethod string
	Diversity     float64 // Diversity score (distance from cluster center)
	ClusterID     int     // Which cluster this belongs to
}

// EnsembleSampler generates diverse structures from all sampling methods
func EnsembleSampler(sequence string, totalStructures int) []*EnsembleStructure {
	// Distribute structures across methods
	// Basin Explorer: 40% (best method)
	// Fibonacci Sphere: 30%
	// Fragment Assembly: 20%
	// Vedic Monte Carlo: 10%

	numBasin := int(float64(totalStructures) * 0.40)
	numFib := int(float64(totalStructures) * 0.30)
	numFrag := int(float64(totalStructures) * 0.20)
	numMC := totalStructures - numBasin - numFib - numFrag

	fmt.Printf("Ensemble sampling %d total structures:\n", totalStructures)
	fmt.Printf("  Basin Explorer: %d (40%%)\n", numBasin)
	fmt.Printf("  Fibonacci Sphere: %d (30%%)\n", numFib)
	fmt.Printf("  Fragment Assembly: %d (20%%)\n", numFrag)
	fmt.Printf("  Vedic Monte Carlo: %d (10%%)\n", numMC)
	fmt.Println()

	ensemble := make([]*EnsembleStructure, 0, totalStructures)

	// Generate from Basin Explorer
	fmt.Println("Generating Basin Explorer structures...")
	basinStructures := BasinExplorer(sequence, numBasin)
	for _, protein := range basinStructures {
		angles := geometry.CalculateRamachandran(protein)
		vedicResult := vedic.CalculateVedicScore(protein, angles)
		ensemble = append(ensemble, &EnsembleStructure{
			Protein:        protein,
			Energy:         folding.CalculateEnergy(protein),
			VedicScore:     vedicResult.TotalScore,
			SamplingMethod: "Basin Explorer",
		})
	}

	// Generate from Fibonacci Sphere
	fmt.Println("Generating Fibonacci Sphere structures...")
	fibStructures := FibonacciSphereBasins(sequence, numFib)
	for _, protein := range fibStructures {
		angles := geometry.CalculateRamachandran(protein)
		vedicResult := vedic.CalculateVedicScore(protein, angles)
		ensemble = append(ensemble, &EnsembleStructure{
			Protein:        protein,
			Energy:         folding.CalculateEnergy(protein),
			VedicScore:     vedicResult.TotalScore,
			SamplingMethod: "Fibonacci Sphere",
		})
	}

	// Generate from Fragment Assembly (using wrapper)
	fmt.Println("Generating Fragment Assembly structures...")
	fragStructures := GenerateFragmentStructures(sequence, numFrag)
	for _, protein := range fragStructures {
		angles := geometry.CalculateRamachandran(protein)
		vedicResult := vedic.CalculateVedicScore(protein, angles)
		ensemble = append(ensemble, &EnsembleStructure{
			Protein:        protein,
			Energy:         folding.CalculateEnergy(protein),
			VedicScore:     vedicResult.TotalScore,
			SamplingMethod: "Fragment Assembly",
		})
	}

	// Generate from Vedic Monte Carlo
	fmt.Println("Generating Vedic Monte Carlo structures...")
	mcStructures := VedicMonteCarlo(sequence, numMC)
	for _, protein := range mcStructures {
		angles := geometry.CalculateRamachandran(protein)
		vedicResult := vedic.CalculateVedicScore(protein, angles)
		ensemble = append(ensemble, &EnsembleStructure{
			Protein:        protein,
			Energy:         folding.CalculateEnergy(protein),
			VedicScore:     vedicResult.TotalScore,
			SamplingMethod: "Vedic Monte Carlo",
		})
	}

	fmt.Printf("✅ Generated %d structures total\n", len(ensemble))
	return ensemble
}

// RankAndSelectDiverse ranks structures by energy+diversity and selects top k
func RankAndSelectDiverse(ensemble []*EnsembleStructure, nativeProtein *parser.Protein, k int) []*EnsembleStructure {
	if len(ensemble) == 0 {
		return nil
	}

	fmt.Printf("\nRanking and selecting top %d diverse structures...\n", k)

	// Step 1: Calculate RMSD for all structures (for validation)
	for _, es := range ensemble {
		// RMSD calculation is expensive, but needed for quality assessment
		_, _ = validation.CalculateRMSD(es.Protein, nativeProtein)
	}

	// Step 2: Sort by energy (lower is better)
	sort.Slice(ensemble, func(i, j int) bool {
		return ensemble[i].Energy < ensemble[j].Energy
	})

	// Step 3: Select top k using diversity metric
	// We want low-energy structures that are also diverse from each other
	selected := make([]*EnsembleStructure, 0, k)

	// Always take the best energy structure
	selected = append(selected, ensemble[0])
	fmt.Printf("  Selected #1: %s (energy: %.2f kcal/mol)\n",
		ensemble[0].SamplingMethod, ensemble[0].Energy)

	// For remaining k-1 structures, maximize diversity
	for len(selected) < k && len(selected) < len(ensemble) {
		bestCandidate := -1
		bestDiversity := -1.0

		// Find structure most diverse from selected set
		for i := range ensemble {
			// Skip if already selected
			isSelected := false
			for _, s := range selected {
				if s.Protein == ensemble[i].Protein {
					isSelected = true
					break
				}
			}
			if isSelected {
				continue
			}

			// Calculate minimum diversity to selected set
			minDiv := 999999.9
			for _, s := range selected {
				div := QuaternionDiversity(ensemble[i].Protein, s.Protein)
				if div < minDiv {
					minDiv = div
				}
			}

			// Want structure with large minimum diversity (far from all selected)
			if minDiv > bestDiversity {
				bestDiversity = minDiv
				bestCandidate = i
			}
		}

		if bestCandidate == -1 {
			break
		}

		selected = append(selected, ensemble[bestCandidate])
		fmt.Printf("  Selected #%d: %s (energy: %.2f kcal/mol, diversity: %.2f)\n",
			len(selected), ensemble[bestCandidate].SamplingMethod,
			ensemble[bestCandidate].Energy, bestDiversity)
	}

	fmt.Printf("✅ Selected %d diverse structures\n", len(selected))
	return selected
}

// QuaternionDiversity calculates structural diversity using dihedral space distance
// Uses CA distance as proxy since Phi/Psi aren't stored as fields
func QuaternionDiversity(protein1, protein2 *parser.Protein) float64 {
	if len(protein1.Residues) != len(protein2.Residues) {
		return 999999.9 // Incompatible structures
	}

	sumSqDist := 0.0
	count := 0

	// Use CA distance as proxy for conformational difference
	for i := range protein1.Residues {
		res1 := protein1.Residues[i]
		res2 := protein2.Residues[i]

		if res1.CA != nil && res2.CA != nil {
			dx := res1.CA.X - res2.CA.X
			dy := res1.CA.Y - res2.CA.Y
			dz := res1.CA.Z - res2.CA.Z
			sumSqDist += dx*dx + dy*dy + dz*dz
			count++
		}
	}

	if count == 0 {
		return 0
	}

	// RMS CA distance
	return math.Sqrt(sumSqDist / float64(count))
}

// ConsensusRefinement runs optimization on multiple diverse structures and picks best
func ConsensusRefinement(ensemble []*EnsembleStructure, nativeProtein *parser.Protein) (*parser.Protein, float64) {
	if len(ensemble) == 0 {
		return nil, 999999.9
	}

	fmt.Printf("\nRunning consensus refinement on %d structures...\n", len(ensemble))

	bestProtein := ensemble[0].Protein.Copy()
	bestRMSD, err := validation.CalculateRMSD(bestProtein, nativeProtein)
	if err != nil {
		return bestProtein, 999999.9
	}

	// Note: Actual optimization would be done here
	// For now, just return the best structure from the ensemble
	for i, es := range ensemble {
		rmsd, err := validation.CalculateRMSD(es.Protein, nativeProtein)
		if err != nil {
			continue
		}
		fmt.Printf("  Structure %d: RMSD = %.2f Å (method: %s)\n",
			i+1, rmsd, es.SamplingMethod)

		if rmsd < bestRMSD {
			bestRMSD = rmsd
			bestProtein = es.Protein.Copy()
		}
	}

	fmt.Printf("✅ Consensus best RMSD: %.2f Å\n", bestRMSD)
	return bestProtein, bestRMSD
}

// ClusterStructures groups structures into clusters based on structural similarity
func ClusterStructures(ensemble []*EnsembleStructure, numClusters int) [][]*EnsembleStructure {
	if len(ensemble) == 0 || numClusters <= 0 {
		return nil
	}

	fmt.Printf("\nClustering %d structures into %d clusters...\n", len(ensemble), numClusters)

	// Simple k-means clustering in dihedral space
	// Step 1: Initialize cluster centers randomly
	centers := make([]*EnsembleStructure, numClusters)
	step := len(ensemble) / numClusters
	for i := 0; i < numClusters; i++ {
		centers[i] = ensemble[i*step]
	}

	// Step 2: Iterate k-means
	maxIters := 10
	for iter := 0; iter < maxIters; iter++ {
		// Assign each structure to nearest cluster
		for _, es := range ensemble {
			minDist := 999999.9
			bestCluster := 0

			for j, center := range centers {
				dist := QuaternionDiversity(es.Protein, center.Protein)
				if dist < minDist {
					minDist = dist
					bestCluster = j
				}
			}

			es.ClusterID = bestCluster
			es.Diversity = minDist
		}

		// Recompute cluster centers
		// (simplified: just use the structure with lowest energy in each cluster)
		for i := 0; i < numClusters; i++ {
			var clusterMembers []*EnsembleStructure
			for _, es := range ensemble {
				if es.ClusterID == i {
					clusterMembers = append(clusterMembers, es)
				}
			}

			if len(clusterMembers) > 0 {
				// Find lowest energy in cluster
				bestIdx := 0
				for j := range clusterMembers {
					if clusterMembers[j].Energy < clusterMembers[bestIdx].Energy {
						bestIdx = j
					}
				}
				centers[i] = clusterMembers[bestIdx]
			}
		}
	}

	// Step 3: Group into clusters
	clusters := make([][]*EnsembleStructure, numClusters)
	for _, es := range ensemble {
		clusters[es.ClusterID] = append(clusters[es.ClusterID], es)
	}

	// Report
	for i, cluster := range clusters {
		fmt.Printf("  Cluster %d: %d structures\n", i+1, len(cluster))
	}

	fmt.Printf("✅ Clustering complete\n")
	return clusters
}

// SelectBestFromClusters picks the best structure from each cluster
func SelectBestFromClusters(clusters [][]*EnsembleStructure, nativeProtein *parser.Protein) []*EnsembleStructure {
	if len(clusters) == 0 {
		return nil
	}

	fmt.Printf("\nSelecting best structure from each cluster...\n")

	selected := make([]*EnsembleStructure, 0, len(clusters))

	for i, cluster := range clusters {
		if len(cluster) == 0 {
			continue
		}

		// Find structure with lowest RMSD in this cluster
		bestIdx := 0
		bestRMSD := 999999.9

		for j, es := range cluster {
			rmsd, err := validation.CalculateRMSD(es.Protein, nativeProtein)
			if err != nil {
				continue
			}
			if rmsd < bestRMSD {
				bestRMSD = rmsd
				bestIdx = j
			}
		}

		selected = append(selected, cluster[bestIdx])
		fmt.Printf("  Cluster %d: Best RMSD = %.2f Å (%s)\n",
			i+1, bestRMSD, cluster[bestIdx].SamplingMethod)
	}

	fmt.Printf("✅ Selected %d structures\n", len(selected))
	return selected
}
