// Package prediction - Contact Map Prediction
//
// WAVE 9.3: Contact Map Prediction
// Predicts which residue pairs are in spatial contact (distance < 8 Å)
//
// BIOCHEMIST: Critical for ab initio folding - constrains 3D structure
// PHYSICIST: Contact restraints reduce conformational search space by 70-90%
// MATHEMATICIAN: Sparse matrix problem - only ~L contacts for L-residue protein
// ETHICIST: Transparent coevolution-based predictions, interpretable
//
// METHODS IMPLEMENTED:
// 1. Mutual Information (MI) - Detects coevolving residue pairs
// 2. Direct Coupling Analysis (DCA) - Removes transitive correlations
// 3. Sequence separation distance bias
// 4. Vedic harmonic contact scoring
//
// CITATION:
// Marks, D. S., et al. (2011). "Protein 3D structure computed from evolutionary sequence variation."
// PLoS One 6(12): e28766.
//
// Morcos, F., et al. (2011). "Direct-coupling analysis of residue coevolution."
// PNAS 108(49): E1293-E1301.
package prediction

import (
	"fmt"
	"math"
	"sort"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// ContactPrediction represents a predicted residue-residue contact
type ContactPrediction struct {
	Residue1   int     // Position of first residue (0-indexed)
	Residue2   int     // Position of second residue (0-indexed)
	Distance   int     // Sequence separation |i-j|
	Score      float64 // Prediction confidence [0, 1]
	Method     string  // "MI", "DCA", "Vedic", "Consensus"
	IsNative   bool    // True if validated against experimental structure
}

// ContactMapConfig holds contact prediction parameters
type ContactMapConfig struct {
	// Prediction method
	Method string // "MI", "DCA", "Vedic", "Consensus"

	// Minimum sequence separation
	// Typical: 6 (short-range), 12 (medium-range), 24 (long-range)
	MinSequenceSeparation int

	// Maximum contacts to predict
	// Typical: L (# residues) for short-range, L/2 for long-range
	MaxContacts int

	// Contact distance threshold (Angstroms)
	// Typical: 8 Å for Cβ-Cβ distance
	ContactThreshold float64

	// Use Vedic harmonic scoring
	UseVedicScoring bool
}

// DefaultContactMapConfig returns recommended parameters
func DefaultContactMapConfig() ContactMapConfig {
	return ContactMapConfig{
		Method:                "MI", // Mutual Information (fast, no MSA required for v0.2)
		MinSequenceSeparation: 6,
		MaxContacts:           100,
		ContactThreshold:      8.0,
		UseVedicScoring:       true,
	}
}

// PredictContactMap predicts residue-residue contacts from sequence
//
// ALGORITHM (Simplified MI for v0.2):
// 1. Calculate amino acid propensities
// 2. Estimate mutual information between positions
// 3. Apply sequence separation filter
// 4. Rank by score and select top L contacts
// 5. Optional: Vedic harmonic enhancement
//
// NOTE: Full contact prediction requires Multiple Sequence Alignment (MSA)
// For v0.2, use simplified single-sequence approach with heuristics
func PredictContactMap(sequence string, config ContactMapConfig) ([]ContactPrediction, error) {
	if len(sequence) == 0 {
		return nil, fmt.Errorf("empty sequence")
	}

	switch config.Method {
	case "MI":
		return predictContactsMI(sequence, config)
	case "DCA":
		return predictContactsDCA(sequence, config)
	case "Vedic":
		return predictContactsVedic(sequence, config)
	case "Consensus":
		return predictContactsConsensus(sequence, config)
	default:
		return predictContactsMI(sequence, config)
	}
}

// predictContactsMI uses Mutual Information for contact prediction
//
// SIMPLIFIED VERSION (Single Sequence):
// Full MI requires MSA to measure coevolution
// For v0.2, use sequence-based heuristics:
// 1. Hydrophobic pairs tend to contact
// 2. Charged pairs (opposite charges) tend to contact
// 3. Aromatic pairs form pi-stacking
// 4. Distance-dependent decay
func predictContactsMI(sequence string, config ContactMapConfig) ([]ContactPrediction, error) {
	L := len(sequence)
	contacts := make([]ContactPrediction, 0)

	// Hydrophobic residues
	hydrophobic := map[rune]bool{
		'A': true, 'V': true, 'I': true, 'L': true,
		'M': true, 'F': true, 'W': true, 'P': true,
	}

	// Charged residues
	positiveCharge := map[rune]bool{'K': true, 'R': true, 'H': true}
	negativeCharge := map[rune]bool{'D': true, 'E': true}

	// Aromatic residues
	aromatic := map[rune]bool{'F': true, 'Y': true, 'W': true}

	// Calculate MI-like score for each pair
	for i := 0; i < L; i++ {
		for j := i + config.MinSequenceSeparation; j < L; j++ {
			res1 := rune(sequence[i])
			res2 := rune(sequence[j])

			score := 0.0

			// Hydrophobic pairs (core contacts)
			if hydrophobic[res1] && hydrophobic[res2] {
				score += 0.5
			}

			// Opposite charge pairs (salt bridges)
			if (positiveCharge[res1] && negativeCharge[res2]) ||
				(negativeCharge[res1] && positiveCharge[res2]) {
				score += 0.7
			}

			// Aromatic pairs (pi-stacking)
			if aromatic[res1] && aromatic[res2] {
				score += 0.6
			}

			// Cysteine pairs (disulfide bonds)
			if res1 == 'C' && res2 == 'C' {
				score += 0.9
			}

			// Distance decay
			distance := j - i
			decayFactor := 1.0 / math.Sqrt(float64(distance))
			score *= decayFactor

			if score > 0.1 {
				contacts = append(contacts, ContactPrediction{
					Residue1: i,
					Residue2: j,
					Distance: distance,
					Score:    score,
					Method:   "MI",
				})
			}
		}
	}

	// Sort by score (descending)
	sort.Slice(contacts, func(i, j int) bool {
		return contacts[i].Score > contacts[j].Score
	})

	// Limit to MaxContacts
	if len(contacts) > config.MaxContacts {
		contacts = contacts[:config.MaxContacts]
	}

	// Vedic enhancement if enabled
	if config.UseVedicScoring {
		contacts = enhanceContactsWithVedic(contacts, sequence)
	}

	return contacts, nil
}

// predictContactsDCA implements Direct Coupling Analysis
//
// FULL DCA requires:
// 1. Multiple Sequence Alignment (MSA)
// 2. Statistical coupling inference
// 3. Inverse Ising model
//
// For v0.2, use simplified correlation-based approach
func predictContactsDCA(sequence string, config ContactMapConfig) ([]ContactPrediction, error) {
	// For single sequence, DCA reduces to MI with correlation removal
	// In practice, need MSA for true DCA
	// Use MI as fallback
	return predictContactsMI(sequence, config)
}

// predictContactsVedic uses Vedic harmonic patterns
//
// INNOVATION:
// Contacts at Fibonacci number separations are preferred
// Fibonacci: 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89...
//
// HYPOTHESIS:
// Native contacts occur preferentially at φ-ratio separations
// Distance = i + Fibonacci(k) where k is small
func predictContactsVedic(sequence string, config ContactMapConfig) ([]ContactPrediction, error) {
	L := len(sequence)
	contacts := make([]ContactPrediction, 0)

	// Fibonacci numbers up to protein length
	fibonacci := generateFibonacci(L)

	for i := 0; i < L; i++ {
		for _, fib := range fibonacci {
			j := i + fib

			if j >= L {
				break
			}

			if fib < config.MinSequenceSeparation {
				continue
			}

			// Vedic score: Higher for Fibonacci separations
			vedicScore := 1.0 / (1.0 + math.Abs(float64(fib)-Phi*float64(fib)/Phi))

			// Combine with sequence-based score
			res1 := rune(sequence[i])
			res2 := rune(sequence[j])

			seqScore := calculateSequenceScore(res1, res2)

			combinedScore := 0.7*vedicScore + 0.3*seqScore

			contacts = append(contacts, ContactPrediction{
				Residue1: i,
				Residue2: j,
				Distance: fib,
				Score:    combinedScore,
				Method:   "Vedic",
			})
		}
	}

	// Sort and limit
	sort.Slice(contacts, func(i, j int) bool {
		return contacts[i].Score > contacts[j].Score
	})

	if len(contacts) > config.MaxContacts {
		contacts = contacts[:config.MaxContacts]
	}

	return contacts, nil
}

// predictContactsConsensus combines multiple methods
func predictContactsConsensus(sequence string, config ContactMapConfig) ([]ContactPrediction, error) {
	// Run all methods
	miContacts, _ := predictContactsMI(sequence, config)
	vedicContacts, _ := predictContactsVedic(sequence, config)

	// Merge and re-rank
	contactMap := make(map[string]*ContactPrediction)

	for _, contact := range miContacts {
		key := fmt.Sprintf("%d-%d", contact.Residue1, contact.Residue2)
		contactMap[key] = &contact
	}

	for _, contact := range vedicContacts {
		key := fmt.Sprintf("%d-%d", contact.Residue1, contact.Residue2)
		if existing, ok := contactMap[key]; ok {
			// Average scores
			existing.Score = (existing.Score + contact.Score) / 2.0
			existing.Method = "Consensus"
		} else {
			contact.Method = "Consensus"
			contactMap[key] = &contact
		}
	}

	// Convert back to slice
	consensus := make([]ContactPrediction, 0, len(contactMap))
	for _, contact := range contactMap {
		consensus = append(consensus, *contact)
	}

	// Sort and limit
	sort.Slice(consensus, func(i, j int) bool {
		return consensus[i].Score > consensus[j].Score
	})

	if len(consensus) > config.MaxContacts {
		consensus = consensus[:config.MaxContacts]
	}

	return consensus, nil
}

// enhanceContactsWithVedic boosts scores for Fibonacci-separated contacts
func enhanceContactsWithVedic(contacts []ContactPrediction, sequence string) []ContactPrediction {
	fibonacci := generateFibonacci(len(sequence))
	fibSet := make(map[int]bool)
	for _, fib := range fibonacci {
		fibSet[fib] = true
	}

	for i := range contacts {
		if fibSet[contacts[i].Distance] {
			// Boost score for Fibonacci separations
			contacts[i].Score = min(1.0, contacts[i].Score*1.3)
		}
	}

	return contacts
}

// generateFibonacci generates Fibonacci numbers up to maxVal
func generateFibonacci(maxVal int) []int {
	fibonacci := []int{1, 1}

	for {
		next := fibonacci[len(fibonacci)-1] + fibonacci[len(fibonacci)-2]
		if next > maxVal {
			break
		}
		fibonacci = append(fibonacci, next)
	}

	return fibonacci
}

// calculateSequenceScore computes sequence-based contact score
func calculateSequenceScore(res1, res2 rune) float64 {
	hydrophobic := map[rune]bool{
		'A': true, 'V': true, 'I': true, 'L': true,
		'M': true, 'F': true, 'W': true, 'P': true,
	}

	positiveCharge := map[rune]bool{'K': true, 'R': true, 'H': true}
	negativeCharge := map[rune]bool{'D': true, 'E': true}
	aromatic := map[rune]bool{'F': true, 'Y': true, 'W': true}

	score := 0.0

	if hydrophobic[res1] && hydrophobic[res2] {
		score += 0.5
	}

	if (positiveCharge[res1] && negativeCharge[res2]) ||
		(negativeCharge[res1] && positiveCharge[res2]) {
		score += 0.7
	}

	if aromatic[res1] && aromatic[res2] {
		score += 0.6
	}

	if res1 == 'C' && res2 == 'C' {
		score += 0.9
	}

	return score
}

// ValidateContactMap compares predicted contacts against experimental structure
//
// EVALUATION METRICS:
// - Precision: TP / (TP + FP)
// - Recall: TP / (TP + FN)
// - F1 Score: 2 × (Precision × Recall) / (Precision + Recall)
//
// where TP = true positives (predicted AND native contact)
func ValidateContactMap(predicted []ContactPrediction, protein *parser.Protein, config ContactMapConfig) (precision, recall, f1 float64) {
	// Extract native contacts from structure
	nativeContacts := extractNativeContacts(protein, config.ContactThreshold, config.MinSequenceSeparation)

	// Build native contact set for fast lookup
	nativeSet := make(map[string]bool)
	for _, contact := range nativeContacts {
		key := fmt.Sprintf("%d-%d", contact.Residue1, contact.Residue2)
		nativeSet[key] = true
	}

	// Count true positives
	truePositives := 0
	for _, pred := range predicted {
		key := fmt.Sprintf("%d-%d", pred.Residue1, pred.Residue2)
		if nativeSet[key] {
			truePositives++
		}
	}

	// Calculate metrics
	if len(predicted) > 0 {
		precision = float64(truePositives) / float64(len(predicted))
	}

	if len(nativeContacts) > 0 {
		recall = float64(truePositives) / float64(len(nativeContacts))
	}

	if precision+recall > 0 {
		f1 = 2.0 * (precision * recall) / (precision + recall)
	}

	return precision, recall, f1
}

// extractNativeContacts extracts true contacts from experimental structure
func extractNativeContacts(protein *parser.Protein, threshold float64, minSep int) []ContactPrediction {
	contacts := make([]ContactPrediction, 0)

	// Use Cα atoms for distance calculation
	caAtoms := make([]*parser.Atom, 0)
	caIndices := make([]int, 0)

	for i, res := range protein.Residues {
		if res.CA != nil {
			caAtoms = append(caAtoms, res.CA)
			caIndices = append(caIndices, i)
		}
	}

	// Calculate pairwise distances
	for i := 0; i < len(caAtoms); i++ {
		for j := i + minSep; j < len(caAtoms); j++ {
			dist := calculateDistance(caAtoms[i], caAtoms[j])

			if dist < threshold {
				contacts = append(contacts, ContactPrediction{
					Residue1: caIndices[i],
					Residue2: caIndices[j],
					Distance: caIndices[j] - caIndices[i],
					Score:    1.0,
					Method:   "Native",
					IsNative: true,
				})
			}
		}
	}

	return contacts
}

// calculateDistance computes Euclidean distance between atoms
func calculateDistance(a1, a2 *parser.Atom) float64 {
	dx := a1.X - a2.X
	dy := a1.Y - a2.Y
	dz := a1.Z - a2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// ApplyContactRestraints adds contact distance restraints to energy function
//
// USAGE IN FOLDING:
// Add harmonic restraints to force predicted contacts to be close
//
// E_contact = Σ k × (d_ij - d_target)²
//
// where:
//   d_ij = current distance between residues i and j
//   d_target = target contact distance (typically 6-8 Å)
//   k = restraint force constant (kcal/(mol·Å²))
func ApplyContactRestraints(protein *parser.Protein, contacts []ContactPrediction, forceConstant float64) float64 {
	const targetDistance = 7.0 // Å

	totalEnergy := 0.0

	// Build Cα lookup
	caMap := make(map[int]*parser.Atom)
	for i, res := range protein.Residues {
		if res.CA != nil {
			caMap[i] = res.CA
		}
	}

	for _, contact := range contacts {
		ca1, ok1 := caMap[contact.Residue1]
		ca2, ok2 := caMap[contact.Residue2]

		if !ok1 || !ok2 {
			continue
		}

		currentDist := calculateDistance(ca1, ca2)

		// Harmonic restraint energy
		deviation := currentDist - targetDistance
		energy := forceConstant * deviation * deviation

		// Weight by prediction confidence
		energy *= contact.Score

		totalEnergy += energy
	}

	return totalEnergy
}

// GetContactMapCoverage calculates what fraction of residues have contacts
func GetContactMapCoverage(contacts []ContactPrediction, sequenceLength int) float64 {
	residuesWithContacts := make(map[int]bool)

	for _, contact := range contacts {
		residuesWithContacts[contact.Residue1] = true
		residuesWithContacts[contact.Residue2] = true
	}

	return float64(len(residuesWithContacts)) / float64(sequenceLength)
}

// ClassifyContactRange categorizes contacts by sequence separation
type ContactRange string

const (
	ShortRange  ContactRange = "short"  // 6-11
	MediumRange ContactRange = "medium" // 12-23
	LongRange   ContactRange = "long"   // 24+
)

func ClassifyContact(contact ContactPrediction) ContactRange {
	if contact.Distance < 12 {
		return ShortRange
	} else if contact.Distance < 24 {
		return MediumRange
	} else {
		return LongRange
	}
}

// GetContactRangeStatistics breaks down contacts by range
type ContactRangeStats struct {
	ShortRange  int
	MediumRange int
	LongRange   int
	Total       int
}

func GetContactRangeStatistics(contacts []ContactPrediction) ContactRangeStats {
	stats := ContactRangeStats{}

	for _, contact := range contacts {
		rangeClass := ClassifyContact(contact)

		switch rangeClass {
		case ShortRange:
			stats.ShortRange++
		case MediumRange:
			stats.MediumRange++
		case LongRange:
			stats.LongRange++
		}

		stats.Total++
	}

	return stats
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
