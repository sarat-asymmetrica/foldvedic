package sampling

import (
	"strings"
)

// SSType represents secondary structure type
type SSType int

const (
	SSCoil SSType = iota
	SSHelix
	SSSheet
)

func (s SSType) String() string {
	switch s {
	case SSHelix:
		return "H"
	case SSSheet:
		return "E"
	default:
		return "C"
	}
}

// SecondaryStructurePrediction holds predicted secondary structure
type SecondaryStructurePrediction struct {
	Sequence  string
	SSTypes   []SSType
	Confidence []float64 // Confidence of each prediction (0-1)
}

// PredictSecondaryStructure uses simplified Chou-Fasman algorithm
// Returns predicted secondary structure for each residue
func PredictSecondaryStructure(sequence string) SecondaryStructurePrediction {
	n := len(sequence)
	prediction := SecondaryStructurePrediction{
		Sequence:   sequence,
		SSTypes:    make([]SSType, n),
		Confidence: make([]float64, n),
	}

	if n == 0 {
		return prediction
	}

	// Chou-Fasman propensity parameters (simplified)
	// Based on statistical analysis of known protein structures
	// Scale: >1.0 = favors that structure, <1.0 = disfavors
	helixPropensity := map[rune]float64{
		'A': 1.42, 'C': 0.70, 'D': 1.01, 'E': 1.51, 'F': 1.13,
		'G': 0.57, 'H': 1.00, 'I': 1.08, 'K': 1.16, 'L': 1.21,
		'M': 1.45, 'N': 0.67, 'P': 0.57, 'Q': 1.11, 'R': 0.98,
		'S': 0.77, 'T': 0.83, 'V': 1.06, 'W': 1.08, 'Y': 0.69,
	}

	sheetPropensity := map[rune]float64{
		'A': 0.83, 'C': 1.19, 'D': 0.54, 'E': 0.37, 'F': 1.38,
		'G': 0.75, 'H': 0.87, 'I': 1.60, 'K': 0.74, 'L': 1.30,
		'M': 1.05, 'N': 0.89, 'P': 0.55, 'Q': 1.10, 'R': 0.93,
		'S': 0.75, 'T': 1.19, 'V': 1.70, 'W': 1.37, 'Y': 1.47,
	}

	// Window size for averaging propensities
	windowSize := 6

	// Calculate average propensities for each position
	for i := 0; i < n; i++ {
		// Define window around position i
		start := max(0, i-windowSize/2)
		end := min(n, i+windowSize/2+1)

		helixScore := 0.0
		sheetScore := 0.0
		count := 0

		for j := start; j < end; j++ {
			residue := rune(sequence[j])
			if h, ok := helixPropensity[residue]; ok {
				helixScore += h
				count++
			}
			if s, ok := sheetPropensity[residue]; ok {
				sheetScore += s
			}
		}

		if count > 0 {
			helixScore /= float64(count)
			sheetScore /= float64(count)
		}

		// Assign secondary structure type
		// Thresholds from Chou-Fasman
		helixThreshold := 1.03
		sheetThreshold := 1.05

		if helixScore > helixThreshold && helixScore > sheetScore {
			prediction.SSTypes[i] = SSHelix
			prediction.Confidence[i] = (helixScore - helixThreshold) / helixThreshold
		} else if sheetScore > sheetThreshold && sheetScore > helixScore {
			prediction.SSTypes[i] = SSSheet
			prediction.Confidence[i] = (sheetScore - sheetThreshold) / sheetThreshold
		} else {
			prediction.SSTypes[i] = SSCoil
			prediction.Confidence[i] = 0.5 // Low confidence for coil
		}

		// Clamp confidence to [0, 1]
		if prediction.Confidence[i] < 0 {
			prediction.Confidence[i] = 0
		}
		if prediction.Confidence[i] > 1 {
			prediction.Confidence[i] = 1
		}
	}

	// Post-processing: Remove very short helices/sheets (min length 4)
	minHelixLength := 4
	minSheetLength := 3

	smoothSSTypes := make([]SSType, n)
	copy(smoothSSTypes, prediction.SSTypes)

	// Find and remove short helices
	i := 0
	for i < n {
		if prediction.SSTypes[i] == SSHelix {
			start := i
			for i < n && prediction.SSTypes[i] == SSHelix {
				i++
			}
			length := i - start
			if length < minHelixLength {
				// Too short, convert to coil
				for j := start; j < i; j++ {
					smoothSSTypes[j] = SSCoil
				}
			}
		} else {
			i++
		}
	}

	// Find and remove short sheets
	i = 0
	for i < n {
		if smoothSSTypes[i] == SSSheet {
			start := i
			for i < n && smoothSSTypes[i] == SSSheet {
				i++
			}
			length := i - start
			if length < minSheetLength {
				// Too short, convert to coil
				for j := start; j < i; j++ {
					smoothSSTypes[j] = SSCoil
				}
			}
		} else {
			i++
		}
	}

	prediction.SSTypes = smoothSSTypes

	return prediction
}

// GetHelixRegions returns ranges of helical residues
func (p *SecondaryStructurePrediction) GetHelixRegions() [][2]int {
	regions := [][2]int{}
	inHelix := false
	start := 0

	for i, ssType := range p.SSTypes {
		if ssType == SSHelix && !inHelix {
			start = i
			inHelix = true
		} else if ssType != SSHelix && inHelix {
			regions = append(regions, [2]int{start, i})
			inHelix = false
		}
	}

	if inHelix {
		regions = append(regions, [2]int{start, len(p.SSTypes)})
	}

	return regions
}

// GetSheetRegions returns ranges of sheet residues
func (p *SecondaryStructurePrediction) GetSheetRegions() [][2]int {
	regions := [][2]int{}
	inSheet := false
	start := 0

	for i, ssType := range p.SSTypes {
		if ssType == SSSheet && !inSheet {
			start = i
			inSheet = true
		} else if ssType != SSSheet && inSheet {
			regions = append(regions, [2]int{start, i})
			inSheet = false
		}
	}

	if inSheet {
		regions = append(regions, [2]int{start, len(p.SSTypes)})
	}

	return regions
}

// String returns human-readable secondary structure string
func (p *SecondaryStructurePrediction) String() string {
	var sb strings.Builder
	for _, ssType := range p.SSTypes {
		sb.WriteString(ssType.String())
	}
	return sb.String()
}

// GetAccuracy compares prediction to known structure (for validation)
func (p *SecondaryStructurePrediction) GetAccuracy(trueSSTypes []SSType) float64 {
	if len(p.SSTypes) != len(trueSSTypes) {
		return 0.0
	}

	correct := 0
	for i := range p.SSTypes {
		if p.SSTypes[i] == trueSSTypes[i] {
			correct++
		}
	}

	return float64(correct) / float64(len(p.SSTypes))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
